// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.3
// source: google/longrunning/operations.proto

package longrunning

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	status "google.golang.org/genproto/googleapis/rpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This resource represents a long-running operation that is the result of a
// network API call.
type Operation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The server-assigned name, which is only unique within the same service that
	// originally returns it. If you use the default HTTP mapping, the
	// `name` should be a resource name ending with `operations/{unique_id}`.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Service-specific metadata associated with the operation.  It typically
	// contains progress information and common metadata such as create time.
	// Some services might not provide such metadata.  Any method that returns a
	// long-running operation should document the metadata type, if any.
	Metadata *anypb.Any `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// If the value is `false`, it means the operation is still in progress.
	// If `true`, the operation is completed, and either `error` or `response` is
	// available.
	Done bool `protobuf:"varint,3,opt,name=done,proto3" json:"done,omitempty"`
	// The operation result, which can be either an `error` or a valid `response`.
	// If `done` == `false`, neither `error` nor `response` is set.
	// If `done` == `true`, exactly one of `error` or `response` is set.
	//
	// Types that are assignable to Result:
	//
	//	*Operation_Error
	//	*Operation_Response
	Result isOperation_Result `protobuf_oneof:"result"`
}

func (x *Operation) Reset() {
	*x = Operation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Operation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Operation) ProtoMessage() {}

func (x *Operation) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Operation.ProtoReflect.Descriptor instead.
func (*Operation) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{0}
}

func (x *Operation) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Operation) GetMetadata() *anypb.Any {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Operation) GetDone() bool {
	if x != nil {
		return x.Done
	}
	return false
}

func (m *Operation) GetResult() isOperation_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (x *Operation) GetError() *status.Status {
	if x, ok := x.GetResult().(*Operation_Error); ok {
		return x.Error
	}
	return nil
}

func (x *Operation) GetResponse() *anypb.Any {
	if x, ok := x.GetResult().(*Operation_Response); ok {
		return x.Response
	}
	return nil
}

type isOperation_Result interface {
	isOperation_Result()
}

type Operation_Error struct {
	// The error result of the operation in case of failure or cancellation.
	Error *status.Status `protobuf:"bytes,4,opt,name=error,proto3,oneof"`
}

type Operation_Response struct {
	// The normal response of the operation in case of success.  If the original
	// method returns no data on success, such as `Delete`, the response is
	// `google.protobuf.Empty`.  If the original method is standard
	// `Get`/`Create`/`Update`, the response should be the resource.  For other
	// methods, the response should have the type `XxxResponse`, where `Xxx`
	// is the original method name.  For example, if the original method name
	// is `TakeSnapshot()`, the inferred response type is
	// `TakeSnapshotResponse`.
	Response *anypb.Any `protobuf:"bytes,5,opt,name=response,proto3,oneof"`
}

func (*Operation_Error) isOperation_Result() {}

func (*Operation_Response) isOperation_Result() {}

// The request message for [Operations.GetOperation][google.longrunning.Operations.GetOperation].
type GetOperationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the operation resource.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetOperationRequest) Reset() {
	*x = GetOperationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetOperationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOperationRequest) ProtoMessage() {}

func (x *GetOperationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOperationRequest.ProtoReflect.Descriptor instead.
func (*GetOperationRequest) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{1}
}

func (x *GetOperationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The request message for [Operations.ListOperations][google.longrunning.Operations.ListOperations].
type ListOperationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the operation's parent resource.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// The standard list filter.
	Filter string `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	// The standard list page size.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The standard list page token.
	PageToken string `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListOperationsRequest) Reset() {
	*x = ListOperationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOperationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOperationsRequest) ProtoMessage() {}

func (x *ListOperationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOperationsRequest.ProtoReflect.Descriptor instead.
func (*ListOperationsRequest) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{2}
}

func (x *ListOperationsRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListOperationsRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

func (x *ListOperationsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListOperationsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

// The response message for [Operations.ListOperations][google.longrunning.Operations.ListOperations].
type ListOperationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of operations that matches the specified filter in the request.
	Operations []*Operation `protobuf:"bytes,1,rep,name=operations,proto3" json:"operations,omitempty"`
	// The standard List next-page token.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListOperationsResponse) Reset() {
	*x = ListOperationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOperationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOperationsResponse) ProtoMessage() {}

func (x *ListOperationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOperationsResponse.ProtoReflect.Descriptor instead.
func (*ListOperationsResponse) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{3}
}

func (x *ListOperationsResponse) GetOperations() []*Operation {
	if x != nil {
		return x.Operations
	}
	return nil
}

func (x *ListOperationsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// The request message for [Operations.CancelOperation][google.longrunning.Operations.CancelOperation].
type CancelOperationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the operation resource to be cancelled.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *CancelOperationRequest) Reset() {
	*x = CancelOperationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelOperationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelOperationRequest) ProtoMessage() {}

func (x *CancelOperationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelOperationRequest.ProtoReflect.Descriptor instead.
func (*CancelOperationRequest) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{4}
}

func (x *CancelOperationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The request message for [Operations.DeleteOperation][google.longrunning.Operations.DeleteOperation].
type DeleteOperationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the operation resource to be deleted.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DeleteOperationRequest) Reset() {
	*x = DeleteOperationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteOperationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteOperationRequest) ProtoMessage() {}

func (x *DeleteOperationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteOperationRequest.ProtoReflect.Descriptor instead.
func (*DeleteOperationRequest) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteOperationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The request message for [Operations.WaitOperation][google.longrunning.Operations.WaitOperation].
type WaitOperationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the operation resource to wait on.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The maximum duration to wait before timing out. If left blank, the wait
	// will be at most the time permitted by the underlying HTTP/RPC protocol.
	// If RPC context deadline is also specified, the shorter one will be used.
	Timeout *durationpb.Duration `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *WaitOperationRequest) Reset() {
	*x = WaitOperationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WaitOperationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WaitOperationRequest) ProtoMessage() {}

func (x *WaitOperationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WaitOperationRequest.ProtoReflect.Descriptor instead.
func (*WaitOperationRequest) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{6}
}

func (x *WaitOperationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *WaitOperationRequest) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

// A message representing the message types used by a long-running operation.
//
// Example:
//
//	rpc LongRunningRecognize(LongRunningRecognizeRequest)
//	    returns (google.longrunning.Operation) {
//	  option (google.longrunning.operation_info) = {
//	    response_type: "LongRunningRecognizeResponse"
//	    metadata_type: "LongRunningRecognizeMetadata"
//	  };
//	}
type OperationInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. The message name of the primary return type for this
	// long-running operation.
	// This type will be used to deserialize the LRO's response.
	//
	// If the response is in a different package from the rpc, a fully-qualified
	// message name must be used (e.g. `google.protobuf.Struct`).
	//
	// Note: Altering this value constitutes a breaking change.
	ResponseType string `protobuf:"bytes,1,opt,name=response_type,json=responseType,proto3" json:"response_type,omitempty"`
	// Required. The message name of the metadata type for this long-running
	// operation.
	//
	// If the response is in a different package from the rpc, a fully-qualified
	// message name must be used (e.g. `google.protobuf.Struct`).
	//
	// Note: Altering this value constitutes a breaking change.
	MetadataType string `protobuf:"bytes,2,opt,name=metadata_type,json=metadataType,proto3" json:"metadata_type,omitempty"`
}

func (x *OperationInfo) Reset() {
	*x = OperationInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_longrunning_operations_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OperationInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperationInfo) ProtoMessage() {}

func (x *OperationInfo) ProtoReflect() protoreflect.Message {
	mi := &file_google_longrunning_operations_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperationInfo.ProtoReflect.Descriptor instead.
func (*OperationInfo) Descriptor() ([]byte, []int) {
	return file_google_longrunning_operations_proto_rawDescGZIP(), []int{7}
}

func (x *OperationInfo) GetResponseType() string {
	if x != nil {
		return x.ResponseType
	}
	return ""
}

func (x *OperationInfo) GetMetadataType() string {
	if x != nil {
		return x.MetadataType
	}
	return ""
}

var file_google_longrunning_operations_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*OperationInfo)(nil),
		Field:         1049,
		Name:          "google.longrunning.operation_info",
		Tag:           "bytes,1049,opt,name=operation_info",
		Filename:      "google/longrunning/operations.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// Additional information regarding long-running operations.
	// In particular, this specifies the types that are returned from
	// long-running operations.
	//
	// Required for methods that return `google.longrunning.Operation`; invalid
	// otherwise.
	//
	// optional google.longrunning.OperationInfo operation_info = 1049;
	E_OperationInfo = &file_google_longrunning_operations_proto_extTypes[0]
)

var File_google_longrunning_operations_proto protoreflect.FileDescriptor

var file_google_longrunning_operations_proto_rawDesc = []byte{
	0x0a, 0x23, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x2f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f,
	0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x01, 0x0a, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x48, 0x00,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x32, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x48,
	0x00, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x29, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0x7f, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x7f, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x2c, 0x0a, 0x16, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0x2c, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x5f,
	0x0a, 0x14, 0x57, 0x61, 0x69, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22,
	0x59, 0x0a, 0x0d, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x32, 0xaa, 0x05, 0x0a, 0x0a, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x94, 0x01, 0x0a, 0x0e, 0x4c, 0x69,
	0x73, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x29, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e,
	0x67, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x76, 0x31,
	0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x7d, 0xda, 0x41, 0x0b, 0x6e, 0x61, 0x6d, 0x65, 0x2c, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x12, 0x7f, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x27, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75,
	0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a,
	0x12, 0x18, 0x2f, 0x76, 0x31, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x2a, 0x2a, 0x7d, 0xda, 0x41, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x7e, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f,
	0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a,
	0x2a, 0x18, 0x2f, 0x76, 0x31, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x2a, 0x2a, 0x7d, 0xda, 0x41, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x88, 0x01, 0x0a, 0x0f, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c,
	0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65,
	0x6c, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x31, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x24, 0x22, 0x1f, 0x2f, 0x76, 0x31, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x2a, 0x2a, 0x7d, 0x3a, 0x63, 0x61, 0x6e, 0x63,
	0x65, 0x6c, 0x3a, 0x01, 0x2a, 0xda, 0x41, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x5a, 0x0a, 0x0d,
	0x57, 0x61, 0x69, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69,
	0x6e, 0x67, 0x2e, 0x57, 0x61, 0x69, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x1a, 0x1d, 0xca, 0x41, 0x1a, 0x6c, 0x6f, 0x6e,
	0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61,
	0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x3a, 0x69, 0x0a, 0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x99, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75,
	0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x42, 0xa8, 0x01, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x42, 0x0f, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65,
	0x6e, 0x65, 0x6c, 0x61, 0x2f, 0x6d, 0x73, 0x68, 0x69, 0x70, 0x2f, 0x74, 0x68, 0x69, 0x72, 0x64,
	0x5f, 0x70, 0x61, 0x72, 0x74, 0x79, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x3b, 0x6c, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67,
	0xf8, 0x01, 0x01, 0xaa, 0x02, 0x12, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x4c, 0x6f, 0x6e,
	0x67, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0xca, 0x02, 0x12, 0x47, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x5c, 0x4c, 0x6f, 0x6e, 0x67, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_google_longrunning_operations_proto_rawDescOnce sync.Once
	file_google_longrunning_operations_proto_rawDescData = file_google_longrunning_operations_proto_rawDesc
)

func file_google_longrunning_operations_proto_rawDescGZIP() []byte {
	file_google_longrunning_operations_proto_rawDescOnce.Do(func() {
		file_google_longrunning_operations_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_longrunning_operations_proto_rawDescData)
	})
	return file_google_longrunning_operations_proto_rawDescData
}

var file_google_longrunning_operations_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_google_longrunning_operations_proto_goTypes = []interface{}{
	(*Operation)(nil),                  // 0: google.longrunning.Operation
	(*GetOperationRequest)(nil),        // 1: google.longrunning.GetOperationRequest
	(*ListOperationsRequest)(nil),      // 2: google.longrunning.ListOperationsRequest
	(*ListOperationsResponse)(nil),     // 3: google.longrunning.ListOperationsResponse
	(*CancelOperationRequest)(nil),     // 4: google.longrunning.CancelOperationRequest
	(*DeleteOperationRequest)(nil),     // 5: google.longrunning.DeleteOperationRequest
	(*WaitOperationRequest)(nil),       // 6: google.longrunning.WaitOperationRequest
	(*OperationInfo)(nil),              // 7: google.longrunning.OperationInfo
	(*anypb.Any)(nil),                  // 8: google.protobuf.Any
	(*status.Status)(nil),              // 9: google.rpc.Status
	(*durationpb.Duration)(nil),        // 10: google.protobuf.Duration
	(*descriptorpb.MethodOptions)(nil), // 11: google.protobuf.MethodOptions
	(*emptypb.Empty)(nil),              // 12: google.protobuf.Empty
}
var file_google_longrunning_operations_proto_depIdxs = []int32{
	8,  // 0: google.longrunning.Operation.metadata:type_name -> google.protobuf.Any
	9,  // 1: google.longrunning.Operation.error:type_name -> google.rpc.Status
	8,  // 2: google.longrunning.Operation.response:type_name -> google.protobuf.Any
	0,  // 3: google.longrunning.ListOperationsResponse.operations:type_name -> google.longrunning.Operation
	10, // 4: google.longrunning.WaitOperationRequest.timeout:type_name -> google.protobuf.Duration
	11, // 5: google.longrunning.operation_info:extendee -> google.protobuf.MethodOptions
	7,  // 6: google.longrunning.operation_info:type_name -> google.longrunning.OperationInfo
	2,  // 7: google.longrunning.Operations.ListOperations:input_type -> google.longrunning.ListOperationsRequest
	1,  // 8: google.longrunning.Operations.GetOperation:input_type -> google.longrunning.GetOperationRequest
	5,  // 9: google.longrunning.Operations.DeleteOperation:input_type -> google.longrunning.DeleteOperationRequest
	4,  // 10: google.longrunning.Operations.CancelOperation:input_type -> google.longrunning.CancelOperationRequest
	6,  // 11: google.longrunning.Operations.WaitOperation:input_type -> google.longrunning.WaitOperationRequest
	3,  // 12: google.longrunning.Operations.ListOperations:output_type -> google.longrunning.ListOperationsResponse
	0,  // 13: google.longrunning.Operations.GetOperation:output_type -> google.longrunning.Operation
	12, // 14: google.longrunning.Operations.DeleteOperation:output_type -> google.protobuf.Empty
	12, // 15: google.longrunning.Operations.CancelOperation:output_type -> google.protobuf.Empty
	0,  // 16: google.longrunning.Operations.WaitOperation:output_type -> google.longrunning.Operation
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	6,  // [6:7] is the sub-list for extension type_name
	5,  // [5:6] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_google_longrunning_operations_proto_init() }
func file_google_longrunning_operations_proto_init() {
	if File_google_longrunning_operations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_google_longrunning_operations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Operation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_longrunning_operations_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetOperationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_longrunning_operations_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListOperationsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_longrunning_operations_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListOperationsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_longrunning_operations_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelOperationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_longrunning_operations_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteOperationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_longrunning_operations_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WaitOperationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_longrunning_operations_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OperationInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_google_longrunning_operations_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Operation_Error)(nil),
		(*Operation_Response)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_longrunning_operations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 1,
			NumServices:   1,
		},
		GoTypes:           file_google_longrunning_operations_proto_goTypes,
		DependencyIndexes: file_google_longrunning_operations_proto_depIdxs,
		MessageInfos:      file_google_longrunning_operations_proto_msgTypes,
		ExtensionInfos:    file_google_longrunning_operations_proto_extTypes,
	}.Build()
	File_google_longrunning_operations_proto = out.File
	file_google_longrunning_operations_proto_rawDesc = nil
	file_google_longrunning_operations_proto_goTypes = nil
	file_google_longrunning_operations_proto_depIdxs = nil
}
