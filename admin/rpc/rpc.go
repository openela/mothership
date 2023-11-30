package mothershipadmin_rpc

import (
	"github.com/openela/mothership/base"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	"github.com/openela/mothership/third_party/googleapis/google/longrunning"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
)

type Server struct {
	base.GRPCServer

	mshipadminpb.UnimplementedMshipAdminServer
	longrunning.UnimplementedOperationsServer

	db       *base.DB
	temporal client.Client
}

func NewServer(db *base.DB, temporalClient client.Client, oidcInterceptorDetails *base.OidcInterceptorDetails, opts ...base.GRPCServerOption) (*Server, error) {
	oidcInterceptor, err := base.OidcGrpcInterceptor(oidcInterceptorDetails)
	if err != nil {
		return nil, err
	}

	opts = append(opts, base.WithUnaryInterceptors(oidcInterceptor))
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
		mshipadminpb.RegisterMshipAdminServer(server, s)
	})
	if err := s.GatewayEndpoints(
		longrunning.RegisterOperationsHandler,
		mshipadminpb.RegisterMshipAdminHandler,
	); err != nil {
		return err
	}

	return s.GRPCServer.Start()
}
