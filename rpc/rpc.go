// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package mothership_rpc

import (
	"github.com/openela/mothership/base"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/openela/mothership/third_party/googleapis/google/longrunning"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
)

type Server struct {
	base.GRPCServer

	mothershippb.UnimplementedSrpmArchiverServer
	longrunning.UnimplementedOperationsServer

	db       *base.DB
	temporal client.Client
}

func NewServer(db *base.DB, temporalClient client.Client, opts ...base.GRPCServerOption) (*Server, error) {
	grpcServer, err := base.NewGRPCServer(opts...)
	if err != nil {
		return nil, err
	}

	return &Server{
		GRPCServer: *grpcServer,
		db:         db,
		temporal:   temporalClient,
	}, nil
}

func (s *Server) Start() error {
	s.RegisterService(func(server *grpc.Server) {
		longrunning.RegisterOperationsServer(server, s)
		mothershippb.RegisterSrpmArchiverServer(server, s)
	})
	if err := s.GatewayEndpoints(
		longrunning.RegisterOperationsHandler,
		mothershippb.RegisterSrpmArchiverHandler,
	); err != nil {
		return err
	}

	return s.GRPCServer.Start()
}
