package mothership_rpc

import (
	"context"
	"github.com/openela/mothership/base"
	mothership_db "github.com/openela/mothership/db"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/openela/mothership/third_party/googleapis/google/longrunning"
	mothership_worker_server "github.com/openela/mothership/worker_server"
	"go.ciq.dev/pika"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
)

var codenameRegexp = regexp.MustCompile(` \(([^)]+)\)`)

func cleanupTrademarks(s string) string {
	return codenameRegexp.ReplaceAllString(strings.ReplaceAll(s, "Red Hat Enterprise Linux release", "OpenELA"), "")
}

func (s *Server) GetEntry(ctx context.Context, req *mothershippb.GetEntryRequest) (*mothershippb.Entry, error) {
	entry, err := base.Q[mothership_db.Entry](s.db).F("name", req.Name).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get entry: %v", err)
		return nil, status.Error(codes.Internal, "failed to get entry")
	}

	if entry == nil {
		return nil, status.Error(codes.NotFound, "entry not found")
	}

	pb := entry.ToPB()
	pb.OsRelease = cleanupTrademarks(pb.OsRelease)

	// If on hold, let's query temporal for more info.
	if entry.State == mothershippb.Entry_ON_HOLD {
		events := s.temporal.GetWorkflowHistory(ctx, "operations/"+entry.Sha256Sum, "", false, enumspb.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
		// We only need to find the latest ImportRPM event.
		// Return the error from that event.
		pb.ErrorMessage = "Unknown error"
		for events.HasNext() {
			event, err := events.Next()
			if err != nil {
				base.LogErrorf("failed to get next event: %v", err)
				continue
			}
			failedAttrs := event.GetActivityTaskFailedEventAttributes()
			if failedAttrs == nil {
				continue
			}

			pb.ErrorMessage = failedAttrs.Failure.Message
			break
		}
	}

	return pb, nil
}

func (s *Server) ListEntries(_ context.Context, req *mothershippb.ListEntriesRequest) (*mothershippb.ListEntriesResponse, error) {
	aipOptions := pika.ProtoReflect(&mothershippb.Entry{})

	page, nt, err := base.Q[mothership_db.Entry](s.db).GetPage(req, aipOptions)
	if err != nil {
		if strings.Contains(err.Error(), "applying filter from page token") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		base.LogErrorf("failed to get entry page: %v", err)
		return nil, status.Error(codes.Internal, "failed to get entry page")
	}

	var entries []*mothershippb.Entry
	for _, e := range page {
		pb := e.ToPB()
		pb.OsRelease = cleanupTrademarks(pb.OsRelease)
		entries = append(entries, pb)
	}

	return &mothershippb.ListEntriesResponse{
		Entries:       entries,
		NextPageToken: nt,
	}, nil
}

// SubmitEntry handles the RPC request for submitting an entry. This is usually
// called by the worker. The worker must be authenticated. The checksum will "lease"
// the entry for the worker, so that other workers will not submit the same entry.
// This "lease" is enforced using Temporal
func (s *Server) SubmitEntry(ctx context.Context, req *mothershippb.SubmitEntryRequest) (*longrunning.Operation, error) {
	worker, err := s.getWorkerIdentity(ctx)
	if err != nil {
		return nil, err
	}

	// Now make sure the entry doesn't already exist in the ARCHIVED state.
	// If it does, return an error. It should be retracted first.
	entry, err := base.Q[mothership_db.Entry](s.db).F(
		"sha256_sum", req.ProcessRpmRequest.Checksum,
		"state", mothershippb.Entry_ARCHIVED,
	).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get entry: %v", err)
		return nil, status.Error(codes.Internal, "failed to get entry")
	}
	if entry != nil {
		return nil, status.Error(codes.AlreadyExists, "entry already exists, you must retract the entry before submitting again")
	}

	startWorkflowOpts := client.StartWorkflowOptions{
		ID:                                       "operations/" + req.ProcessRpmRequest.Checksum,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
		WorkflowIDReusePolicy:                    enumspb.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	}

	// Submit to Temporal
	run, err := s.temporal.ExecuteWorkflow(
		context.Background(),
		startWorkflowOpts,
		mothership_worker_server.ProcessRPMWorkflow,
		&mothershippb.ProcessRPMArgs{
			Request: req.ProcessRpmRequest,
			InternalRequest: &mothershippb.ProcessRPMInternalRequest{
				WorkerId: worker.WorkerID,
			},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "is already running") {
			return nil, status.Error(codes.AlreadyExists, "entry is already running")
		}
		base.LogErrorf("failed to start workflow: %v", err)
		return nil, status.Error(codes.Internal, "failed to start workflow")
	}

	return s.getOperation(ctx, run.GetID())
}
