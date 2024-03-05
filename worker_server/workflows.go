package mothership_worker_server

import (
	"time"

	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	mothershippb "github.com/openela/mothership/proto/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var w Worker

// processRPMPostHold is a part of the ProcessRPM workflow.
// This part executes the import part, and retries if it fails.
// After the first failure, the workflow is put on hold.
// If the workflow is put on hold, the workflow can be rescued by an admin.
func processRPMPostHold(ctx workflow.Context, entry *mothershippb.Entry, args *mothershippb.ProcessRPMArgs, num int) (*mothershippb.ProcessRPMResponse, error) {
	// If resource exists, then we can start the import.
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		// We'll wait up to 25 minutes for the import to finish.
		// Most imports are fast, but some packages are very large.
		StartToCloseTimeout: 25 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	})
	var importRpmRes mothershippb.ImportRPMResponse
	err := workflow.ExecuteActivity(ctx, w.ImportRPM, args.Request.RpmUri, args.Request.Checksum, args.Request.OsRelease).Get(ctx, &importRpmRes)
	if err != nil {
		// If the import fails, we'll put the workflow on hold.
		// If the workflow is put on hold, an admin can rescue the workflow.
		var err error
		signalChan := workflow.GetSignalChannel(ctx, "rescue")
		workflow.GetLogger(ctx).Info("Import failed, putting workflow on hold")
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(ctx.Done(), func(c workflow.ReceiveChannel, more bool) {
			err = ctx.Err()
		})
		selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, nil)
			err = nil
		})

		// Set state to on hold
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 25 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 0,
			},
		})
		err = workflow.ExecuteActivity(ctx, w.SetEntryState, entry.Name, mothershippb.Entry_ON_HOLD, nil).Get(ctx, entry)
		if err != nil {
			return nil, err
		}

		// Wait until a rescue signal is received. Otherwise, an admin can also
		// cancel the workflow.
		selector.Select(ctx)

		// Check if workflow was cancelled.
		if err != nil {
			ctx, cancel := workflow.NewDisconnectedContext(ctx)
			defer cancel()
			ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 25 * time.Second,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 0,
				},
			})
			_ = workflow.ExecuteActivity(ctx, w.SetEntryState, entry.Name, mothershippb.Entry_CANCELLED, nil).Get(ctx, entry)
			return nil, err
		}

		// Set the entry state to archiving
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 25 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 0,
			},
		})
		err = workflow.ExecuteActivity(ctx, w.SetEntryState, entry.Name, mothershippb.Entry_ARCHIVING, nil).Get(ctx, entry)
		if err != nil {
			return nil, err
		}

		// If the workflow was not cancelled, then we can retry the import.
		return processRPMPostHold(ctx, entry, args, num+1)
	}

	// If the import succeeds, then we can update the entry state.
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 0,
		},
	})
	err = workflow.ExecuteActivity(ctx, w.SetEntryState, entry.Name, mothershippb.Entry_ARCHIVED, &importRpmRes).Get(ctx, entry)
	if err != nil {
		return nil, err
	}

	// If num > 0, this means the import failed at least once.
	// Let's check if the entry was part of a batch, if so we'll update the ticket
	// with the new status.
	if num > 0 && entry.Batch != nil && entry.Batch.Value != "" {
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Hour,
			HeartbeatTimeout:    25 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				// If it fails more than twice, then let's just not care.
				// A maintainer can edit the ticket manually.
				MaximumAttempts: 2,
			},
		})
		err = workflow.ExecuteActivity(ctx, w.UpdateTicketStatus, entry).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	return &mothershippb.ProcessRPMResponse{
		Entry: entry,
	}, nil
}

// ProcessRPMWorkflow processes an SRPM.
// Usually a client worker will first initiate an upload to the storage backend,
// then send a request to the Server `SubmitEntry` method (or send a request
// then upload the resource).
func ProcessRPMWorkflow(ctx workflow.Context, args *mothershippb.ProcessRPMArgs) (*mothershippb.ProcessRPMResponse, error) {
	// First verify that the resource exists.
	// The resource can be uploaded after the request is sent.
	// So we should wait up to 2 hours. The initial timeouts should be low
	// since the worker is most likely to upload the resource immediately.
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			// We're waiting 25 seconds each time
			InitialInterval:    25 * time.Second,
			BackoffCoefficient: 1,
			// Maximum attempts should be set, so it's approximately 2 hours
			MaximumAttempts: (60 * 60 * 2) / 25,
		},
	})
	err := workflow.ExecuteActivity(ctx, w.VerifyResourceExists, args.Request.RpmUri).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Set worker last check in time
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
	})
	err = workflow.ExecuteActivity(ctx, w.SetWorkerLastCheckinTime, args.InternalRequest.WorkerId).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Create an entry, if the import fails, we'll still have an entry.
	// If it succeeds, we'll update the entry state.
	// If it fails we can set the workflow on hold and if the patches are updated
	// an admin can signal and "rescue" the workflow.
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
	})
	var entry mothershippb.Entry
	err = workflow.ExecuteActivity(ctx, w.CreateEntry, args).Get(ctx, &entry)
	if err != nil {
		return nil, err
	}

	// On defer, if the workflow is not completed, then we'll set the entry state
	// to failed.
	defer func() {
		if entry.State == mothershippb.Entry_ARCHIVED || entry.State == mothershippb.Entry_CANCELLED {
			return
		}

		ctx, _ := workflow.NewDisconnectedContext(ctx)
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 25 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 0,
			},
		})

		// Check if entry has EntryID set, if not then we can just delete the entry
		if entry.EntryId == "" {
			_ = workflow.ExecuteActivity(ctx, w.DeleteEntry, entry.Name).Get(ctx, nil)
			return
		}
		_ = workflow.ExecuteActivity(ctx, w.SetEntryState, entry.Name, mothershippb.Entry_FAILED, nil).Get(ctx, nil)
	}()

	// Set the entry name to the RPM NVR
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 45 * time.Second,
	})
	err = workflow.ExecuteActivity(ctx, w.SetEntryIDFromRPM, entry.Name, args.Request.RpmUri, args.Request.Checksum).Get(ctx, &entry)
	if err != nil {
		return nil, err
	}

	// Process the RPM.
	return processRPMPostHold(ctx, &entry, args, 0)
}

// RetractEntryWorkflow retracts an entry.
// Should be used when an entry debranding is not considered fully complete. (Contains upstream trademarks for example)
// This will forcefully remove the commit from the git repository and set the entry state to RETRACTED.
// The same source (for the specific entry) can be re-imported by the client, either by calling DuplicateEntry or
// calling SubmitEntry with the same SRPM URI.
func RetractEntryWorkflow(ctx workflow.Context, name string) (*mshipadminpb.RetractEntryResponse, error) {
	// Set entry state to retracting
	var entry mothershippb.Entry
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
	})
	err := workflow.ExecuteActivity(ctx, w.SetEntryState, name, mothershippb.Entry_RETRACTING, nil).Get(ctx, &entry)
	if err != nil {
		return nil, err
	}

	// Deferring this activity will set the entry state to ARCHIVED if the workflow
	// is not completed.
	defer func() {
		if entry.State == mothershippb.Entry_RETRACTED {
			return
		}

		// This is because the entry is still archived, but the commit was not
		// retracted.
		ctx, cancel := workflow.NewDisconnectedContext(ctx)
		defer cancel()
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 25 * time.Second,
		})
		_ = workflow.ExecuteActivity(ctx, w.SetEntryState, name, mothershippb.Entry_ARCHIVED, nil).Get(ctx, nil)
	}()

	// Retract commit
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
	})

	var res mshipadminpb.RetractEntryResponse
	err = workflow.ExecuteActivity(ctx, w.RetractEntry, name).Get(ctx, &res)
	if err != nil {
		return nil, err
	}

	// Set the entry state to retracted
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
	})
	err = workflow.ExecuteActivity(ctx, w.SetEntryState, name, mothershippb.Entry_RETRACTED, nil).Get(ctx, &entry)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// SealBatchWorkflow seals a batch.
// After a worker finishing submitting their entries, they can seal the batch.
// Sealing the batch will wait for all entries to reach a "stop" condition.
// This can mean "ARCHIVED" or "ON_HOLD". If all entries are in a stop condition,
// a new ticket is created in the ticketing system with the status of each entry.
// After creating the ticket, the batch is sealed and won't accept any more entries.
func SealBatchWorkflow(ctx workflow.Context, req *mothershippb.SealBatchRequest) (*mothershippb.SealBatchResponse, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Hour,
		HeartbeatTimeout:    25 * time.Second,
	})
	err := workflow.ExecuteActivity(ctx, w.WaitForEntriesToSettle, req).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Create ticket
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 40 * time.Second,
	})
	err = workflow.ExecuteActivity(ctx, w.CreateTicket, req.Name).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Seal batch
	var batch mothershippb.Batch
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
	})
	err = workflow.ExecuteActivity(ctx, w.SealBatch, req.Name).Get(ctx, &batch)
	if err != nil {
		return nil, err
	}

	return &mothershippb.SealBatchResponse{
		Batch: &batch,
	}, nil
}
