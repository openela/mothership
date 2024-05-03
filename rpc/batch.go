// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package mothership_rpc

import (
	"context"
	"database/sql"
	"strings"
	"time"

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
)

func (s *Server) GetBatch(_ context.Context, req *mothershippb.GetBatchRequest) (*mothershippb.Batch, error) {
	batch, err := base.Q[mothership_db.BatchView](s.db).F("name", req.Name).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get batch: %v", err)
		return nil, status.Error(codes.Internal, "failed to get batch")
	}

	if batch == nil {
		return nil, status.Error(codes.NotFound, "batch not found")
	}

	return batch.ToPB(), nil
}

func (s *Server) ListBatches(_ context.Context, req *mothershippb.ListBatchesRequest) (*mothershippb.ListBatchesResponse, error) {
	aipOptions := pika.ProtoReflect(&mothershippb.Batch{})

	page, nt, err := base.Q[mothership_db.BatchView](s.db).GetPage(req, aipOptions)
	if err != nil {
		base.LogErrorf("failed to get batch page: %v", err)
		return nil, status.Error(codes.Internal, "failed to get batch page")
	}

	return &mothershippb.ListBatchesResponse{
		Batches:       base.SliceToPB[*mothershippb.Batch, *mothership_db.BatchView](page),
		NextPageToken: nt,
	}, nil
}

func (s *Server) CreateBatch(ctx context.Context, req *mothershippb.CreateBatchRequest) (*mothershippb.Batch, error) {
	worker, err := s.getWorkerIdentity(ctx)
	if err != nil {
		return nil, err
	}

	batch := &mothership_db.Batch{
		Name:          base.NameGen("batches"),
		WorkerID:      worker.WorkerID,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		SealTime:      sql.NullTime{},
		BugtrackerURI: sql.NullString{},
	}
	if req.BatchId != "" {
		batch.BatchID = sql.NullString{String: req.BatchId, Valid: true}
	}

	if err := base.Q[mothership_db.Batch](s.db).Create(batch); err != nil {
		base.LogErrorf("failed to create batch: %v", err)
		return nil, status.Error(codes.Internal, "failed to create batch")
	}

	return batch.ToPB(), nil
}

func (s *Server) SealBatch(ctx context.Context, req *mothershippb.SealBatchRequest) (*longrunning.Operation, error) {
	_, err := s.getWorkerIdentity(ctx)
	if err != nil {
		return nil, err
	}

	batch, err := base.Q[mothership_db.Batch](s.db).F("name", req.Name).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get batch: %v", err)
		return nil, status.Error(codes.Internal, "failed to get batch")
	}
	if batch == nil {
		return nil, status.Error(codes.NotFound, "batch not found")
	}

	startWorkflowOpts := client.StartWorkflowOptions{
		ID:                                       "operations/seal/" + batch.Name,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
		WorkflowIDReusePolicy:                    enumspb.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}

	// Submit to Temporal
	run, err := s.temporal.ExecuteWorkflow(
		context.Background(),
		startWorkflowOpts,
		mothership_worker_server.SealBatchWorkflow,
		req,
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
