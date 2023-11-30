package mothershipadmin_rpc

import (
	"context"
	"github.com/openela/mothership/base"
	mothership_db "github.com/openela/mothership/db"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	"go.ciq.dev/pika"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (s *Server) GetWorker(_ context.Context, req *mshipadminpb.GetWorkerRequest) (*mshipadminpb.Worker, error) {
	worker, err := base.Q[mothership_db.Worker](s.db).F("name", req.Name).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get worker: %v", err)
		return nil, status.Error(codes.Internal, "failed to get worker")
	}

	if worker == nil {
		return nil, status.Error(codes.NotFound, "worker not found")
	}

	return worker.ToPB(), nil
}

func (s *Server) ListWorkers(_ context.Context, req *mshipadminpb.ListWorkersRequest) (*mshipadminpb.ListWorkersResponse, error) {
	aipOptions := pika.ProtoReflect(&mshipadminpb.Worker{})

	page, nt, err := base.Q[mothership_db.Worker](s.db).GetPage(req, aipOptions)
	if err != nil {
		base.LogErrorf("failed to get worker page: %v", err)
		return nil, status.Error(codes.Internal, "failed to get worker page")
	}

	return &mshipadminpb.ListWorkersResponse{
		Workers:       base.SliceToPB[*mshipadminpb.Worker, *mothership_db.Worker](page),
		NextPageToken: nt,
	}, nil
}

func (s *Server) CreateWorker(_ context.Context, req *mshipadminpb.CreateWorkerRequest) (*mshipadminpb.Worker, error) {
	// Verify that the id is at least 4 characters long.
	if len(req.WorkerId) < 4 {
		return nil, status.Error(codes.InvalidArgument, "worker id must be at least 4 characters long")
	}

	// Create the worker.
	name := base.NameGen("workers")
	worker := &mothership_db.Worker{
		Name:      name,
		WorkerID:  req.WorkerId,
		ApiSecret: base.NameGen(name),
	}

	err := base.Q[mothership_db.Worker](s.db).Create(worker)
	if err != nil {
		// if the error is a duplicate key error, return an already exists error.
		if strings.Contains(err.Error(), "duplicate key") {
			st, _ := status.
				New(codes.AlreadyExists, "worker already exists").
				WithDetails(&errdetails.LocalizedMessage{
					Locale:  "en-US",
					Message: "Worker with given ID already exists.",
				})

			return nil, st.Err()
		}

		base.LogErrorf("failed to create worker: %v", err)
		return nil, status.Error(codes.Internal, "failed to create worker")
	}
	pb := worker.ToPB()
	pb.ApiSecret = worker.ApiSecret

	return pb, nil
}

func (s *Server) DeleteWorker(_ context.Context, req *mshipadminpb.DeleteWorkerRequest) (*emptypb.Empty, error) {
	worker, err := base.Q[mothership_db.Worker](s.db).F("name", req.Name).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get worker: %v", err)
		return nil, status.Error(codes.Internal, "failed to get worker")
	}

	if worker == nil {
		return nil, status.Error(codes.NotFound, "worker not found")
	}

	err = base.Q[mothership_db.Worker](s.db).D(worker)
	if err != nil {
		base.LogErrorf("failed to delete worker: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete worker")
	}

	return &emptypb.Empty{}, nil
}
