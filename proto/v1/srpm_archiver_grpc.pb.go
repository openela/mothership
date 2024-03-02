// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: proto/v1/srpm_archiver.proto

package mothershippb

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

// SrpmArchiverClient is the client API for SrpmArchiver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SrpmArchiverClient interface {
	// Returns a batch
	GetBatch(ctx context.Context, in *GetBatchRequest, opts ...grpc.CallOption) (*Batch, error)
	// Returns a list of batches that match the filter criteria.
	ListBatches(ctx context.Context, in *ListBatchesRequest, opts ...grpc.CallOption) (*ListBatchesResponse, error)
	// Creates a batch.
	// Only worker credentials can create a batch.
	CreateBatch(ctx context.Context, in *CreateBatchRequest, opts ...grpc.CallOption) (*Batch, error)
	// Returns an entry
	GetEntry(ctx context.Context, in *GetEntryRequest, opts ...grpc.CallOption) (*Entry, error)
	// Returns a list of entries that match the filter criteria.
	ListEntries(ctx context.Context, in *ListEntriesRequest, opts ...grpc.CallOption) (*ListEntriesResponse, error)
	// Submits an SRPM to be archived.
	// A worker can call this method to submit an SRPM to be archived.
	// The call can occur even before uploading the SRPM to the object storage
	// that way it can be ensured that a certain hash is "leased" by the worker.
	// Other workers will still keep the hash in their backlog until the SRPM is
	// verified processed.
	// Until they can query an entry with `sha256_sum=X` matching the hash of the
	// SRPM, it will not be deleted from the backlog.
	// If after 2 hours the SRPM is not processed, the worker can assume that
	// the SRPM is lost and can be re-uploaded. It that case, the entry will be
	// re-assigned to the worker.
	// If a checksum can't be leased because it's already being processed,
	// AlreadyExists error will be returned.
	// The worker MUST stop processing the SRPM in that case.
	SubmitEntry(ctx context.Context, in *SubmitEntryRequest, opts ...grpc.CallOption) (*longrunning.Operation, error)
	// WorkerUploadObject is used by workers to upload objects to the
	// object storage service.
	// Returns AlreadyExists if the SRPM already exists.
	// This doesn't necessarily mean that the worker should stop processing,
	// especially if it acquired a lease to process this particular SRPM.
	WorkerUploadObject(ctx context.Context, opts ...grpc.CallOption) (SrpmArchiver_WorkerUploadObjectClient, error)
	// WorkerPing is used by workers to ping the server.
	// This is used to check if the worker is still alive.
	WorkerPing(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type srpmArchiverClient struct {
	cc grpc.ClientConnInterface
}

func NewSrpmArchiverClient(cc grpc.ClientConnInterface) SrpmArchiverClient {
	return &srpmArchiverClient{cc}
}

func (c *srpmArchiverClient) GetBatch(ctx context.Context, in *GetBatchRequest, opts ...grpc.CallOption) (*Batch, error) {
	out := new(Batch)
	err := c.cc.Invoke(ctx, "/mothership.v1.SrpmArchiver/GetBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srpmArchiverClient) ListBatches(ctx context.Context, in *ListBatchesRequest, opts ...grpc.CallOption) (*ListBatchesResponse, error) {
	out := new(ListBatchesResponse)
	err := c.cc.Invoke(ctx, "/mothership.v1.SrpmArchiver/ListBatches", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srpmArchiverClient) CreateBatch(ctx context.Context, in *CreateBatchRequest, opts ...grpc.CallOption) (*Batch, error) {
	out := new(Batch)
	err := c.cc.Invoke(ctx, "/mothership.v1.SrpmArchiver/CreateBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srpmArchiverClient) GetEntry(ctx context.Context, in *GetEntryRequest, opts ...grpc.CallOption) (*Entry, error) {
	out := new(Entry)
	err := c.cc.Invoke(ctx, "/mothership.v1.SrpmArchiver/GetEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srpmArchiverClient) ListEntries(ctx context.Context, in *ListEntriesRequest, opts ...grpc.CallOption) (*ListEntriesResponse, error) {
	out := new(ListEntriesResponse)
	err := c.cc.Invoke(ctx, "/mothership.v1.SrpmArchiver/ListEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srpmArchiverClient) SubmitEntry(ctx context.Context, in *SubmitEntryRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	out := new(longrunning.Operation)
	err := c.cc.Invoke(ctx, "/mothership.v1.SrpmArchiver/SubmitEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *srpmArchiverClient) WorkerUploadObject(ctx context.Context, opts ...grpc.CallOption) (SrpmArchiver_WorkerUploadObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &SrpmArchiver_ServiceDesc.Streams[0], "/mothership.v1.SrpmArchiver/WorkerUploadObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &srpmArchiverWorkerUploadObjectClient{stream}
	return x, nil
}

type SrpmArchiver_WorkerUploadObjectClient interface {
	Send(*WorkerUploadObjectRequest) error
	CloseAndRecv() (*WorkerUploadObjectResponse, error)
	grpc.ClientStream
}

type srpmArchiverWorkerUploadObjectClient struct {
	grpc.ClientStream
}

func (x *srpmArchiverWorkerUploadObjectClient) Send(m *WorkerUploadObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *srpmArchiverWorkerUploadObjectClient) CloseAndRecv() (*WorkerUploadObjectResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(WorkerUploadObjectResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *srpmArchiverClient) WorkerPing(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/mothership.v1.SrpmArchiver/WorkerPing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SrpmArchiverServer is the server API for SrpmArchiver service.
// All implementations must embed UnimplementedSrpmArchiverServer
// for forward compatibility
type SrpmArchiverServer interface {
	// Returns a batch
	GetBatch(context.Context, *GetBatchRequest) (*Batch, error)
	// Returns a list of batches that match the filter criteria.
	ListBatches(context.Context, *ListBatchesRequest) (*ListBatchesResponse, error)
	// Creates a batch.
	// Only worker credentials can create a batch.
	CreateBatch(context.Context, *CreateBatchRequest) (*Batch, error)
	// Returns an entry
	GetEntry(context.Context, *GetEntryRequest) (*Entry, error)
	// Returns a list of entries that match the filter criteria.
	ListEntries(context.Context, *ListEntriesRequest) (*ListEntriesResponse, error)
	// Submits an SRPM to be archived.
	// A worker can call this method to submit an SRPM to be archived.
	// The call can occur even before uploading the SRPM to the object storage
	// that way it can be ensured that a certain hash is "leased" by the worker.
	// Other workers will still keep the hash in their backlog until the SRPM is
	// verified processed.
	// Until they can query an entry with `sha256_sum=X` matching the hash of the
	// SRPM, it will not be deleted from the backlog.
	// If after 2 hours the SRPM is not processed, the worker can assume that
	// the SRPM is lost and can be re-uploaded. It that case, the entry will be
	// re-assigned to the worker.
	// If a checksum can't be leased because it's already being processed,
	// AlreadyExists error will be returned.
	// The worker MUST stop processing the SRPM in that case.
	SubmitEntry(context.Context, *SubmitEntryRequest) (*longrunning.Operation, error)
	// WorkerUploadObject is used by workers to upload objects to the
	// object storage service.
	// Returns AlreadyExists if the SRPM already exists.
	// This doesn't necessarily mean that the worker should stop processing,
	// especially if it acquired a lease to process this particular SRPM.
	WorkerUploadObject(SrpmArchiver_WorkerUploadObjectServer) error
	// WorkerPing is used by workers to ping the server.
	// This is used to check if the worker is still alive.
	WorkerPing(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedSrpmArchiverServer()
}

// UnimplementedSrpmArchiverServer must be embedded to have forward compatible implementations.
type UnimplementedSrpmArchiverServer struct {
}

func (UnimplementedSrpmArchiverServer) GetBatch(context.Context, *GetBatchRequest) (*Batch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBatch not implemented")
}
func (UnimplementedSrpmArchiverServer) ListBatches(context.Context, *ListBatchesRequest) (*ListBatchesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBatches not implemented")
}
func (UnimplementedSrpmArchiverServer) CreateBatch(context.Context, *CreateBatchRequest) (*Batch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBatch not implemented")
}
func (UnimplementedSrpmArchiverServer) GetEntry(context.Context, *GetEntryRequest) (*Entry, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEntry not implemented")
}
func (UnimplementedSrpmArchiverServer) ListEntries(context.Context, *ListEntriesRequest) (*ListEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEntries not implemented")
}
func (UnimplementedSrpmArchiverServer) SubmitEntry(context.Context, *SubmitEntryRequest) (*longrunning.Operation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitEntry not implemented")
}
func (UnimplementedSrpmArchiverServer) WorkerUploadObject(SrpmArchiver_WorkerUploadObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method WorkerUploadObject not implemented")
}
func (UnimplementedSrpmArchiverServer) WorkerPing(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WorkerPing not implemented")
}
func (UnimplementedSrpmArchiverServer) mustEmbedUnimplementedSrpmArchiverServer() {}

// UnsafeSrpmArchiverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SrpmArchiverServer will
// result in compilation errors.
type UnsafeSrpmArchiverServer interface {
	mustEmbedUnimplementedSrpmArchiverServer()
}

func RegisterSrpmArchiverServer(s grpc.ServiceRegistrar, srv SrpmArchiverServer) {
	s.RegisterService(&SrpmArchiver_ServiceDesc, srv)
}

func _SrpmArchiver_GetBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrpmArchiverServer).GetBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.v1.SrpmArchiver/GetBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrpmArchiverServer).GetBatch(ctx, req.(*GetBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrpmArchiver_ListBatches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBatchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrpmArchiverServer).ListBatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.v1.SrpmArchiver/ListBatches",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrpmArchiverServer).ListBatches(ctx, req.(*ListBatchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrpmArchiver_CreateBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrpmArchiverServer).CreateBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.v1.SrpmArchiver/CreateBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrpmArchiverServer).CreateBatch(ctx, req.(*CreateBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrpmArchiver_GetEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrpmArchiverServer).GetEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.v1.SrpmArchiver/GetEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrpmArchiverServer).GetEntry(ctx, req.(*GetEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrpmArchiver_ListEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrpmArchiverServer).ListEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.v1.SrpmArchiver/ListEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrpmArchiverServer).ListEntries(ctx, req.(*ListEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrpmArchiver_SubmitEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrpmArchiverServer).SubmitEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.v1.SrpmArchiver/SubmitEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrpmArchiverServer).SubmitEntry(ctx, req.(*SubmitEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SrpmArchiver_WorkerUploadObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SrpmArchiverServer).WorkerUploadObject(&srpmArchiverWorkerUploadObjectServer{stream})
}

type SrpmArchiver_WorkerUploadObjectServer interface {
	SendAndClose(*WorkerUploadObjectResponse) error
	Recv() (*WorkerUploadObjectRequest, error)
	grpc.ServerStream
}

type srpmArchiverWorkerUploadObjectServer struct {
	grpc.ServerStream
}

func (x *srpmArchiverWorkerUploadObjectServer) SendAndClose(m *WorkerUploadObjectResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *srpmArchiverWorkerUploadObjectServer) Recv() (*WorkerUploadObjectRequest, error) {
	m := new(WorkerUploadObjectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SrpmArchiver_WorkerPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SrpmArchiverServer).WorkerPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mothership.v1.SrpmArchiver/WorkerPing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SrpmArchiverServer).WorkerPing(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SrpmArchiver_ServiceDesc is the grpc.ServiceDesc for SrpmArchiver service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SrpmArchiver_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mothership.v1.SrpmArchiver",
	HandlerType: (*SrpmArchiverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBatch",
			Handler:    _SrpmArchiver_GetBatch_Handler,
		},
		{
			MethodName: "ListBatches",
			Handler:    _SrpmArchiver_ListBatches_Handler,
		},
		{
			MethodName: "CreateBatch",
			Handler:    _SrpmArchiver_CreateBatch_Handler,
		},
		{
			MethodName: "GetEntry",
			Handler:    _SrpmArchiver_GetEntry_Handler,
		},
		{
			MethodName: "ListEntries",
			Handler:    _SrpmArchiver_ListEntries_Handler,
		},
		{
			MethodName: "SubmitEntry",
			Handler:    _SrpmArchiver_SubmitEntry_Handler,
		},
		{
			MethodName: "WorkerPing",
			Handler:    _SrpmArchiver_WorkerPing_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WorkerUploadObject",
			Handler:       _SrpmArchiver_WorkerUploadObject_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/v1/srpm_archiver.proto",
}