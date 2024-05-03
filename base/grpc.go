// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GRPCServer struct {
	server     *grpc.Server
	gatewayMux *runtime.ServeMux

	gatewayClientConn *grpc.ClientConn

	serverOptions      []grpc.ServerOption
	dialOptions        []grpc.DialOption
	muxOptions         []runtime.ServeMuxOption
	outgoingHeaders    []string
	incomingHeaders    []string
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	timeout            time.Duration
	grpcPort           int
	gatewayPort        int
	noGrpcGateway      bool
	noMetrics          bool

	// ServeMuxOptions
	additionalHeaders map[string]bool
}

type GrpcEndpointRegister func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
type GRPCServerOption func(*GRPCServer)

// WithServerOptions sets the gRPC server options. (Append)
func WithServerOptions(opts ...grpc.ServerOption) GRPCServerOption {
	return func(g *GRPCServer) {
		g.serverOptions = append(g.serverOptions, opts...)
	}
}

// WithDialOptions sets the gRPC dial options. (Append)
func WithDialOptions(opts ...grpc.DialOption) GRPCServerOption {
	return func(g *GRPCServer) {
		g.dialOptions = append(g.dialOptions, opts...)
	}
}

// WithMuxOptions sets the gRPC-gateway mux options. (Append)
func WithMuxOptions(opts ...runtime.ServeMuxOption) GRPCServerOption {
	return func(g *GRPCServer) {
		g.muxOptions = append(g.muxOptions, opts...)
	}
}

// WithOutgoingHeaders sets the outgoing headers for the gRPC-gateway. (Append)
func WithOutgoingHeaders(headers ...string) GRPCServerOption {
	return func(g *GRPCServer) {
		g.outgoingHeaders = append(g.outgoingHeaders, headers...)
	}
}

// WithIncomingHeaders sets the incoming headers for the gRPC-gateway. (Append)
func WithIncomingHeaders(headers ...string) GRPCServerOption {
	return func(g *GRPCServer) {
		g.incomingHeaders = append(g.incomingHeaders, headers...)
	}
}

// WithUnaryInterceptors sets the gRPC unary interceptors. (Append)
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) GRPCServerOption {
	return func(g *GRPCServer) {
		g.unaryInterceptors = append(g.unaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors sets the gRPC stream interceptors. (Append)
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) GRPCServerOption {
	return func(g *GRPCServer) {
		g.streamInterceptors = append(g.streamInterceptors, interceptors...)
	}
}

// WithTimeout sets the timeout for the gRPC server.
func WithTimeout(timeout time.Duration) GRPCServerOption {
	return func(g *GRPCServer) {
		g.timeout = timeout
	}
}

// WithGRPCPort sets the gRPC port for the gRPC server.
func WithGRPCPort(port int) GRPCServerOption {
	return func(g *GRPCServer) {
		g.grpcPort = port
	}
}

// WithGatewayPort sets the gRPC-gateway port for the gRPC server.
func WithGatewayPort(port int) GRPCServerOption {
	return func(g *GRPCServer) {
		g.gatewayPort = port
	}
}

// WithNoGRPCGateway disables the gRPC-gateway for the gRPC server.
func WithNoGRPCGateway() GRPCServerOption {
	return func(g *GRPCServer) {
		g.noGrpcGateway = true
	}
}

// WithNoMetrics disables the Prometheus metrics for the gRPC server.
func WithNoMetrics() GRPCServerOption {
	return func(g *GRPCServer) {
		g.noMetrics = true
	}
}

// WithServeMuxAdditionalHeaders sets additional headers for the gRPC-gateway. (incoming)
func WithServeMuxAdditionalHeaders(headers ...string) GRPCServerOption {
	return func(opts *GRPCServer) {
		if opts.additionalHeaders == nil {
			opts.additionalHeaders = make(map[string]bool)
		}
		for _, header := range headers {
			opts.additionalHeaders[header] = true
		}
	}
}

func DefaultServeMuxOptions(s ...GRPCServer) []runtime.ServeMuxOption {
	additionalHeaders := make(map[string]bool)

	// We only make the arg a spread operator so that we can use the same function
	// for both the default and the user-defined options.
	// It can only ever be one GRPCServer
	if len(s) > 0 {
		if len(s) > 1 {
			slog.Error("DefaultServeMuxOptions: too many arguments")
			os.Exit(1)
		}

		opts := s[0]
		if opts.additionalHeaders != nil {
			for header := range opts.additionalHeaders {
				additionalHeaders[header] = true
			}
		}
	}

	return []runtime.ServeMuxOption{
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			switch strings.ToLower(s) {
			case "authorization",
				"cookie",
				"x-mship-worker-secret":
				return s, true
			}

			if additionalHeaders[s] {
				return s, true
			}

			if strings.ToLower(s) == "content-type" {
				return "original-content-type", true
			}

			return s, false
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}
}

// NewGRPCServer creates a new gRPC-server with gRPC-gateway, default interceptors
// and exposed Prometheus metrics.
func NewGRPCServer(opts ...GRPCServerOption) (*GRPCServer, error) {
	g := &GRPCServer{
		serverOptions:      []grpc.ServerOption{},
		dialOptions:        []grpc.DialOption{},
		muxOptions:         []runtime.ServeMuxOption{},
		outgoingHeaders:    []string{},
		incomingHeaders:    []string{},
		unaryInterceptors:  []grpc.UnaryServerInterceptor{},
		streamInterceptors: []grpc.StreamServerInterceptor{},
	}

	// Apply options first
	for _, opt := range opts {
		opt(g)
	}

	// Set defaults
	if g.timeout == 0 {
		g.timeout = 10 * time.Second
	}
	if g.grpcPort == 0 {
		g.grpcPort = 8080
	}
	if g.gatewayPort == 0 {
		g.gatewayPort = g.grpcPort + 1
	}
	if len(g.muxOptions) == 0 {
		g.muxOptions = DefaultServeMuxOptions()
	}

	// Always prepend the insecure dial option
	// RESF deploys with Istio, which handles mTLS
	g.dialOptions = append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, g.dialOptions...)

	// Set default interceptors
	if g.unaryInterceptors == nil {
		g.unaryInterceptors = []grpc.UnaryServerInterceptor{}
	}
	if g.streamInterceptors == nil {
		g.streamInterceptors = []grpc.StreamServerInterceptor{}
	}

	// Always prepend the prometheus interceptor
	g.unaryInterceptors = append([]grpc.UnaryServerInterceptor{grpc_prometheus.UnaryServerInterceptor}, g.unaryInterceptors...)
	g.streamInterceptors = append([]grpc.StreamServerInterceptor{grpc_prometheus.StreamServerInterceptor}, g.streamInterceptors...)

	// Chain the interceptors
	g.serverOptions = append(g.serverOptions, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(g.unaryInterceptors...)))
	g.serverOptions = append(g.serverOptions, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(g.streamInterceptors...)))

	g.server = grpc.NewServer(g.serverOptions...)

	g.gatewayMux = runtime.NewServeMux(g.muxOptions...)

	// Create gateway client connection
	var err error
	g.gatewayClientConn, err = grpc.Dial("localhost:"+strconv.Itoa(g.grpcPort), g.dialOptions...)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *GRPCServer) RegisterService(register func(*grpc.Server)) {
	register(g.server)
}

func (g *GRPCServer) GatewayEndpoints(registerEndpoints ...GrpcEndpointRegister) error {
	for _, register := range registerEndpoints {
		if err := register(context.Background(), g.gatewayMux, g.gatewayClientConn); err != nil {
			return err
		}
	}

	return nil
}

func (g *GRPCServer) GatewayMux() *runtime.ServeMux {
	return g.gatewayMux
}

func (g *GRPCServer) Start() error {
	// Create gRPC listener
	grpcLis, err := net.Listen("tcp", ":"+strconv.Itoa(g.grpcPort))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	// First start the gRPC server
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		slog.Info("gRPC server running", "port", g.grpcPort)
		grpc_prometheus.Register(g.server)

		err := g.server.Serve(grpcLis)
		if err != nil {
			slog.Error("gRPC server failed to serve", "error", err.Error())
			os.Exit(1)
		}

		slog.Info("gRPC server stopped")
	}(&wg)

	// Then start the gRPC-gateway
	if !g.noGrpcGateway {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			slog.Info("gRPC-gateway running", "port", g.gatewayPort)
			err := http.ListenAndServe(":"+strconv.Itoa(g.gatewayPort), g.gatewayMux)
			if err != nil {
				slog.Error("gRPC-gateway failed to serve", "error", err.Error())
				os.Exit(1)
			}

			slog.Info("gRPC-gateway stopped")
		}(&wg)
	}

	// Serve proxmux
	if !g.noMetrics {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			promMux := http.NewServeMux()
			promMux.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(":7332", promMux)
			if err != nil {
				slog.Error("Prometheus mux failed to serve", "error", err.Error())
			}
		}(&wg)
	}

	wg.Wait()

	return nil
}
