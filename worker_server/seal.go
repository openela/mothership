package mothership_worker_server

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"text/template"
	"time"

	"github.com/openela/mothership/base"
	"github.com/openela/mothership/base/bugtracker"
	mothership_db "github.com/openela/mothership/db"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type issueBody struct {
	Batch     *mothership_db.Batch
	Entries   []*mothership_db.Entry
	PublicURI string
}

func (w *Worker) isEntriesSettled(req *mothershippb.SealBatchRequest) error {
	var entries []*mothership_db.Entry
	var err error
	if req.OperationNames != nil && len(req.OperationNames) > 0 {
		entries, err = base.Q[mothership_db.Entry](w.db).F("batch_name", req.Name).All()
		if err != nil {
			return errors.Wrap(err, "failed to get entries")
		}
		if len(entries) != len(req.OperationNames) {
			return errors.New("not all entries are settled")
		}
	} else {
		entries, err = base.Q[mothership_db.Entry](w.db).F("batch_name", req.Name).All()
		if err != nil {
			return errors.Wrap(err, "failed to get entries")
		}

	}

	allEntriesSettled := true
	for _, entry := range entries {
		if entry.State != mothershippb.Entry_ARCHIVED && entry.State != mothershippb.Entry_ON_HOLD && entry.State != mothershippb.Entry_FAILED && entry.State != mothershippb.Entry_CANCELLED {
			allEntriesSettled = false
			break
		}
	}

	if !allEntriesSettled {
		return errors.New("not all entries are settled")
	}

	return nil
}

func (w *Worker) getTicketInfo(batch *mothership_db.Batch) (string, string, *bugtracker.Options, error) {
	labels := []string{"all-successful"}
	entries, err := base.Q[mothership_db.Entry](w.db).F("batch_name", batch.Name).All()
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to get entries")
	}
	if len(entries) == 0 {
		return "", "", nil, nil
	}

	for _, entry := range entries {
		entryPb := entry.ToPB()
		if entryPb.State != mothershippb.Entry_ARCHIVED {
			labels = []string{"failed-entry"}
			break
		}
	}

	labels = append(labels, "import-batch")

	title := fmt.Sprintf("%s: %s", batch.WorkerID, batch.Name)
	body := `Worker {{.Batch.WorkerID}} sealed {{.Batch.Name}}.

{{if .Entries}}The following entries were in the batch:
{{range .Entries}}- [{{if eq .State 2}}x{{else}} {{end}}] [{{.EntryID}}]({{$.PublicURI}}/{{.Name}})
{{end}}{{else}}No entries were in batch. This is a test, please ignore.{{end}}
`

	bodyTemplate, err := template.New("issueBody").Parse(body)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to parse template")
	}

	var buf bytes.Buffer
	err = bodyTemplate.Execute(&buf, issueBody{
		Batch:     batch,
		Entries:   entries,
		PublicURI: w.publicURI,
	})
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to execute template")
	}
	formattedBody := buf.String()

	opts := bugtracker.Options{
		Labels: labels,
	}

	return title, formattedBody, &opts, nil
}

func (w *Worker) SealBatch(name string) (*mothershippb.Batch, error) {
	batch, err := base.Q[mothership_db.Batch](w.db).F("name", name).GetOrNil()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get batch")
	}
	if batch == nil {
		return nil, temporal.NewNonRetryableApplicationError(
			"batch does not exist",
			"batchDoesNotExist",
			errors.New("batch does not exist"),
		)
	}

	batch.SealTime = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	err = base.Q[mothership_db.Batch](w.db).U(batch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update batch")
	}

	return batch.ToPB(), nil
}

func (w *Worker) WaitForEntriesToSettle(ctx context.Context, req *mothershippb.SealBatchRequest) error {
	for {
		if w.isEntriesSettled(req) == nil {
			break
		}
		activity.RecordHeartbeat(ctx)
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (w *Worker) CreateTicket(ctx context.Context, batchName string) error {
	if w.bugtracker == nil {
		return nil
	}

	batch, err := base.Q[mothership_db.Batch](w.db).F("name", batchName).GetOrNil()
	if err != nil {
		return errors.Wrap(err, "failed to get batch")
	}
	if batch == nil {
		return temporal.NewNonRetryableApplicationError(
			"batch does not exist",
			"batchDoesNotExist",
			errors.New("batch does not exist"),
		)
	}

	title, formattedBody, opts, err := w.getTicketInfo(batch)
	if err != nil {
		return errors.Wrap(err, "failed to get ticket title and body")
	}
	if title == "" {
		return nil
	}

	auth, err := w.bugtracker.GetAuthenticator()
	if err != nil {
		return errors.Wrap(err, "failed to get authenticator")
	}

	ticket, err := w.bugtracker.CreateTicket(auth, title, formattedBody, *opts)
	if err != nil {
		return errors.Wrap(err, "failed to create ticket")
	}

	// Everything beyond this point is best effort.
	// If we fail to update the batch with the ticket URI, it's not the end of the world.
	ticketURI, err := w.bugtracker.TicketURI(ticket)
	if err != nil {
		slog.Info("failed to get ticket URI", "err", err)
	}

	batch.BugtrackerURI = sql.NullString{
		Valid:  true,
		String: ticketURI,
	}

	err = base.Q[mothership_db.Batch](w.db).U(batch)
	if err != nil {
		slog.Info("failed to update batch", "err", err)
	}

	// Close ticket if everything went well, but again not the end of the world if we fail.
	if opts.Labels[0] == "all-successful" {
		err = w.bugtracker.CloseTicket(auth, ticket)
		if err != nil {
			slog.Info("failed to close ticket", "err", err)
		}
	}

	return nil
}

func (w *Worker) UpdateTicketStatus(ctx context.Context, entry *mothershippb.Entry) error {
	if w.bugtracker == nil {
		return nil
	}

	// Wait until batch has a ticket URI
	var batch *mothership_db.Batch
	for {
		var err error
		batch, err = base.Q[mothership_db.Batch](w.db).F("name", entry.Batch.Value).Get()
		if err != nil {
			return errors.Wrap(err, "failed to get batch")
		}

		if batch.BugtrackerURI.Valid {
			break
		}
		activity.RecordHeartbeat(ctx)
		time.Sleep(5 * time.Second)
	}

	auth, err := w.bugtracker.GetAuthenticator()
	if err != nil {
		return errors.Wrap(err, "failed to get authenticator")
	}

	title, formattedBody, opts, err := w.getTicketInfo(batch)
	if err != nil {
		return errors.Wrap(err, "failed to get ticket title and body")
	}
	if title == "" {
		return nil
	}

	ticket, err := w.bugtracker.URIToTicket(batch.BugtrackerURI.String)
	if err != nil {
		return errors.Wrap(err, "failed to get ticket ID")
	}

	err = w.bugtracker.EditTicket(auth, ticket, title, formattedBody, *opts)
	if err != nil {
		return errors.Wrap(err, "failed to edit ticket")
	}

	// Close ticket if everything went well, but again not the end of the world if we fail.
	if opts.Labels[0] == "all-successful" {
		err = w.bugtracker.CloseTicket(auth, ticket)
		if err != nil {
			slog.Info("failed to close ticket", "err", err)
		}
	}

	return nil
}
