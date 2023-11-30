package mothership_rpc

import (
	"context"

	"github.com/openela/mothership/base"
	mothership_db "github.com/openela/mothership/db"
	mothershippb "github.com/openela/mothership/proto/v1"
	"go.ciq.dev/pika"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetBatch(_ context.Context, req *mothershippb.GetBatchRequest) (*mothershippb.Batch, error) {
	batch, err := base.Q[mothership_db.Batch](s.db).F("name", req.Name).GetOrNil()
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

	page, nt, err := base.Q[mothership_db.Batch](s.db).GetPage(req, aipOptions)
	if err != nil {
		base.LogErrorf("failed to get batch page: %v", err)
		return nil, status.Error(codes.Internal, "failed to get batch page")
	}

	return &mothershippb.ListBatchesResponse{
		Batches:       base.SliceToPB[*mothershippb.Batch, *mothership_db.Batch](page),
		NextPageToken: nt,
	}, nil
}

func (s *Server) CreateBatch(ctx context.Context, req *mothershippb.CreateBatchRequest) (*mothershippb.Batch, error) {
	worker, err := s.getWorkerIdentity(ctx)
	if err != nil {
		return nil, err
	}

	batch := &mothership_db.Batch{
		Name:     base.NameGen("batches"),
		BatchID:  req.BatchId,
		WorkerID: worker.WorkerID,
	}

	if err := base.Q[mothership_db.Batch](s.db).Create(batch); err != nil {
		base.LogErrorf("failed to create batch: %v", err)
		return nil, status.Error(codes.Internal, "failed to create batch")
	}

	return batch.ToPB(), nil
}
