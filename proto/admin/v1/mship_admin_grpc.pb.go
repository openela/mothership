// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: proto/admin/v1/mship_admin.proto

package mshipadminpb

import (
	context "context"
	longrunning "github.com/openela/mothership/third_party/googleapis/google/longrunning"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MshipAdminClient is the client API for MshipAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MshipAdminClient interface {
	// Gets a worker
	GetWorker(ctx context.Context, in *GetWorkerRequest, opts ...grpc.CallOption) (*Worker, error)
	// Lists the workers registered
	ListWorkers(ctx context.Context, in *ListWorkersRequest, opts ...grpc.CallOption) (*ListWorkersResponse, error)
	// (-- api-linter: core::0133::http-body=disabled
	//
	//	aip.dev/not-precedent: See below in the CreateWorkerRequest. We only allow worker_id --)
	//
	// Creates a worker
	CreateWorker(ctx context.Context, in *CreateWorkerRequest, opts ...grpc.CallOption) (*Worker, error)
	// Deletes a worker
	// Worker cannot be deleted if it has created an entry.
	DeleteWorker(ctx context.Context, in *DeleteWorkerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Rescue an entry import attempt
	// This should be called after fixing patches that caused the import to fail.
	// This will re-run the import attempt.
	RescueEntryImport(ctx context.Context, in *RescueEntryImportRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Retract the entry
	// To be able to retract an entry, the entry must be in the `ARCHIVED` state.
	// This will allow an NVR to be re-imported.
	RetractEntry(ctx context.Context, in *RetractEntryRequest, opts ...grpc.CallOption) (*longrunning.Operation, error)
}

type mshipAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewMshipAdminClient(cc grpc.ClientConnInterface) MshipAdminClient {
	return &mshipAdminClient{cc}
}

func (c *mshipAdminClient) GetWorker(ctx context.Context, in *GetWorkerRequest, opts ...grpc.CallOption) (*Worker, error) {
	out := new(Worker)
	err := c.cc.Invoke(ctx, "/mothership.admin.v1.MshipAdmin/GetWorker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mshipAdminClient) ListWorkers(ctx context.Context, in *ListWorkersRequest, opts ...grpc.CallOption) (*ListWorkersResponse, error) {
	out := new(ListWorkersResponse)
	err := c.cc.Invoke(ctx, "/mothership.admin.v1.MshipAdmin/ListWorkers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mshipAdminClient) CreateWorker(ctx context.Context, in *CreateWorkerRequest, opts ...grpc.CallOption) (*Worker, error) {
	out := new(Worker)
	err := c.cc.Invoke(ctx, "/mothership.admin.v1.MshipAdmin/CreateWorker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mshipAdminClient) DeleteWorker(ctx context.Context, in *DeleteWorkerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/mothership.admin.v1.MshipAdmin/DeleteWorker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mshipAdminClient) RescueEntryImport(ctx context.Context, in *RescueEntryImportRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/mothership.admin.v1.MshipAdmin/RescueEntryImport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mshipAdminClient) RetractEntry(ctx context.Context, in *RetractEntryRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	out := new(longrunning.Operation)
	err := c.cc.Invoke(ctx, "/mothership.admin.v1.MshipAdmin/RetractEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MshipAdminServer is the server API for MshipAdmin service.
// All implementations must embed UnimplementedMshipAdminServer
// for forward compatibility
type MshipAdminServer interface {
	// Gets a worker
	GetWorker(context.Context, *GetWorkerRequest) (*Worker, error)
	// Lists the workers registered
	ListWorkers(context.Context, *ListWorkersRequest) (*ListWorkersResponse, error)
	// (-- api-linter: core::0133::http-body=disabled
	//
	//	aip.dev/not-precedent: See below in the CreateWorkerRequest. We only allow worker_id --)
	//
	// Creates a worker
	CreateWorker(context.Context, *CreateWorkerRequest) (*Worker, error)
	// Deletes a worker
	// Worker cannot be deleted if it has created an entry.
	DeleteWorker(context.Context, *DeleteWorkerRequest) (*emptypb.Empty, error)
	// Rescue an entry import attempt
	// This should be called after fixing patches that caused the import to fail.
	// This will re-run the import attempt.
	RescueEntryImport(context.Context, *RescueEntryImportRequest) (*emptypb.Empty, error)
	// Retract the entry
	// To be able to retract an entry, the entry must be in the `ARCHIVED` state.
	// This will allow an NVR to be re-imported.
	RetractEntry(context.Context, *RetractEntryRequest) (*longrunning.Operation, error)
	mustEmbedUnimplementedMshipAdminServer()
}

// UnimplementedMshipAdminServer must be embedded to have forward compatible implementations.
type UnimplementedMshipAdminServer struct {
}

func (UnimplementedMshipAdminServer) GetWorker(context.Context, *GetWorkerRequest) (*Worker, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorker not implemented")
}
func (UnimplementedMshipAdminServer) ListWorkers(context.Context, *ListWorkersRequest) (*ListWorkersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWorkers not implemented")
}
func (UnimplementedMshipAdminServer) CreateWorker(context.Context, *CreateWorkerRequest) (*Worker, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWorker not implemented")
}
func (UnimplementedMshipAdminServer) DeleteWorker(context.Context, *DeleteWorkerRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteWorker not implemented")
}
func (UnimplementedMshipAdminServer) RescueEntryImport(context.Context, *RescueEntryImportRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RescueEntryImport not implemented")
}
func (UnimplementedMshipAdminServer) RetractEntry(context.Context, *RetractEntryRequest) (*longrunning.Operation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetractEntry not implemented")
}
func (UnimplementedMshipAdminServer) mustEmbedUnimplementedMshipAdminServer() {}

// UnsafeMshipAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MshipAdminServer will
// result in compilation errors.
type UnsafeMshipAdminServer interface {
	mustEmbedUnimplementedMshipAdminServer()
}

func RegisterMshipAdminServer(s grpc.ServiceRegistrar, srv MshipAdminServer) {
	s.RegisterService(&MshipAdmin_ServiceDesc, srv)
}

func _MshipAdmin_GetWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWorkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MshipAdminServer).GetWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.admin.v1.MshipAdmin/GetWorker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MshipAdminServer).GetWorker(ctx, req.(*GetWorkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MshipAdmin_ListWorkers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWorkersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MshipAdminServer).ListWorkers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.admin.v1.MshipAdmin/ListWorkers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MshipAdminServer).ListWorkers(ctx, req.(*ListWorkersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MshipAdmin_CreateWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWorkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MshipAdminServer).CreateWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.admin.v1.MshipAdmin/CreateWorker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MshipAdminServer).CreateWorker(ctx, req.(*CreateWorkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MshipAdmin_DeleteWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteWorkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MshipAdminServer).DeleteWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.admin.v1.MshipAdmin/DeleteWorker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MshipAdminServer).DeleteWorker(ctx, req.(*DeleteWorkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MshipAdmin_RescueEntryImport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RescueEntryImportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MshipAdminServer).RescueEntryImport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.admin.v1.MshipAdmin/RescueEntryImport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MshipAdminServer).RescueEntryImport(ctx, req.(*RescueEntryImportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MshipAdmin_RetractEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetractEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MshipAdminServer).RetractEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.admin.v1.MshipAdmin/RetractEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MshipAdminServer).RetractEntry(ctx, req.(*RetractEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MshipAdmin_ServiceDesc is the grpc.ServiceDesc for MshipAdmin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MshipAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mothership.admin.v1.MshipAdmin",
	HandlerType: (*MshipAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWorker",
			Handler:    _MshipAdmin_GetWorker_Handler,
		},
		{
			MethodName: "ListWorkers",
			Handler:    _MshipAdmin_ListWorkers_Handler,
		},
		{
			MethodName: "CreateWorker",
			Handler:    _MshipAdmin_CreateWorker_Handler,
		},
		{
			MethodName: "DeleteWorker",
			Handler:    _MshipAdmin_DeleteWorker_Handler,
		},
		{
			MethodName: "RescueEntryImport",
			Handler:    _MshipAdmin_RescueEntryImport_Handler,
		},
		{
			MethodName: "RetractEntry",
			Handler:    _MshipAdmin_RetractEntry_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/admin/v1/mship_admin.proto",
}
