// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package mothership_rpc

import (
	"context"
	"database/sql"
	"github.com/openela/mothership/base"
	mothership_db "github.com/openela/mothership/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (s *Server) WorkerPing(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	worker, err := s.getWorkerIdentity(ctx)
	if err != nil {
		return nil, err
	}

	worker.LastCheckinTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if err := base.Q[mothership_db.Worker](s.db).U(worker); err != nil {
		return nil, status.Error(codes.Internal, "failed to update worker")
	}

	return &emptypb.Empty{}, nil
}
