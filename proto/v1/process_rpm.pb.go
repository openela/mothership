// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.3
// source: proto/v1/process_rpm.proto

package mothershippb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ProcessRPMRequest is the request message for the ProcessRPM workflow
type ProcessRPMRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// URI of the RPM to process
	// e.g. gs://bucket/path/to/rpm.rpm
	// The server must have read access to the RPM and WILL error if it does not
	RpmUri string `protobuf:"bytes,1,opt,name=rpm_uri,json=rpmUri,proto3" json:"rpm_uri,omitempty"`
	// OS Release of the RPM
	// e.g. Red Hat Enterprise Linux release 8.8 (Ootpa)
	OsRelease string `protobuf:"bytes,2,opt,name=os_release,json=osRelease,proto3" json:"os_release,omitempty"`
	// Self reported checksum of the RPM
	// Must be a SHA256 checksum and match the RPM
	Checksum string `protobuf:"bytes,3,opt,name=checksum,proto3" json:"checksum,omitempty"`
	// Self reported repository of the RPM
	// e.g. BaseOS
	Repository string `protobuf:"bytes,4,opt,name=repository,proto3" json:"repository,omitempty"`
	// Batch to associate the RPM with
	Batch string `protobuf:"bytes,5,opt,name=batch,proto3" json:"batch,omitempty"`
}

func (x *ProcessRPMRequest) Reset() {
	*x = ProcessRPMRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_process_rpm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRPMRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRPMRequest) ProtoMessage() {}

func (x *ProcessRPMRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_process_rpm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRPMRequest.ProtoReflect.Descriptor instead.
func (*ProcessRPMRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_process_rpm_proto_rawDescGZIP(), []int{0}
}

func (x *ProcessRPMRequest) GetRpmUri() string {
	if x != nil {
		return x.RpmUri
	}
	return ""
}

func (x *ProcessRPMRequest) GetOsRelease() string {
	if x != nil {
		return x.OsRelease
	}
	return ""
}

func (x *ProcessRPMRequest) GetChecksum() string {
	if x != nil {
		return x.Checksum
	}
	return ""
}

func (x *ProcessRPMRequest) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

func (x *ProcessRPMRequest) GetBatch() string {
	if x != nil {
		return x.Batch
	}
	return ""
}

// ProcessRPMInternalRequest is the request message that the Server
// uses in its call to the ProcessRPM workflow
type ProcessRPMInternalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Worker ID of the worker processing the RPM
	WorkerId string `protobuf:"bytes,1,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
}

func (x *ProcessRPMInternalRequest) Reset() {
	*x = ProcessRPMInternalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_process_rpm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRPMInternalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRPMInternalRequest) ProtoMessage() {}

func (x *ProcessRPMInternalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_process_rpm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRPMInternalRequest.ProtoReflect.Descriptor instead.
func (*ProcessRPMInternalRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_process_rpm_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessRPMInternalRequest) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

// ProcessRPMArgs is the arguments for the ProcessRPM workflow
type ProcessRPMArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Public request
	Request *ProcessRPMRequest `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	// Internal request
	InternalRequest *ProcessRPMInternalRequest `protobuf:"bytes,2,opt,name=internal_request,json=internalRequest,proto3" json:"internal_request,omitempty"`
}

func (x *ProcessRPMArgs) Reset() {
	*x = ProcessRPMArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_process_rpm_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRPMArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRPMArgs) ProtoMessage() {}

func (x *ProcessRPMArgs) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_process_rpm_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRPMArgs.ProtoReflect.Descriptor instead.
func (*ProcessRPMArgs) Descriptor() ([]byte, []int) {
	return file_proto_v1_process_rpm_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessRPMArgs) GetRequest() *ProcessRPMRequest {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *ProcessRPMArgs) GetInternalRequest() *ProcessRPMInternalRequest {
	if x != nil {
		return x.InternalRequest
	}
	return nil
}

// ProcessRPMMetadata is the metadata for the ProcessRPM workflow
type ProcessRPMMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The time at which the workflow started
	StartTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// The time at which the workflow finished
	EndTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *ProcessRPMMetadata) Reset() {
	*x = ProcessRPMMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_process_rpm_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRPMMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRPMMetadata) ProtoMessage() {}

func (x *ProcessRPMMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_process_rpm_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRPMMetadata.ProtoReflect.Descriptor instead.
func (*ProcessRPMMetadata) Descriptor() ([]byte, []int) {
	return file_proto_v1_process_rpm_proto_rawDescGZIP(), []int{3}
}

func (x *ProcessRPMMetadata) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ProcessRPMMetadata) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

// ProcessRPMResponse is the response message for the ProcessRPM workflow
type ProcessRPMResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The entry created for the RPM
	Entry *Entry `protobuf:"bytes,1,opt,name=entry,proto3" json:"entry,omitempty"`
}

func (x *ProcessRPMResponse) Reset() {
	*x = ProcessRPMResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_process_rpm_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRPMResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRPMResponse) ProtoMessage() {}

func (x *ProcessRPMResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_process_rpm_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRPMResponse.ProtoReflect.Descriptor instead.
func (*ProcessRPMResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_process_rpm_proto_rawDescGZIP(), []int{4}
}

func (x *ProcessRPMResponse) GetEntry() *Entry {
	if x != nil {
		return x.Entry
	}
	return nil
}

// ImportRPMResponse is the response message for the ImportRPM activity
type ImportRPMResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Commit hash of the imported RPM
	// e.g. 1234567890abcdef1234567890abcdef12345678
	CommitHash string `protobuf:"bytes,1,opt,name=commit_hash,json=commitHash,proto3" json:"commit_hash,omitempty"`
	// Commit URI of the imported RPM
	CommitUri string `protobuf:"bytes,2,opt,name=commit_uri,json=commitUri,proto3" json:"commit_uri,omitempty"`
	// Commit branch of the imported RPM
	CommitBranch string `protobuf:"bytes,3,opt,name=commit_branch,json=commitBranch,proto3" json:"commit_branch,omitempty"`
	// Commit tag of the imported RPM
	CommitTag string `protobuf:"bytes,4,opt,name=commit_tag,json=commitTag,proto3" json:"commit_tag,omitempty"`
	// NEVRA of the imported RPM
	// e.g. rpm-1.0.0-1.el8.x86_64
	Nevra string `protobuf:"bytes,5,opt,name=nevra,proto3" json:"nevra,omitempty"`
	// Package name of the imported RPM
	// e.g. rpm
	Pkg string `protobuf:"bytes,6,opt,name=pkg,proto3" json:"pkg,omitempty"`
}

func (x *ImportRPMResponse) Reset() {
	*x = ImportRPMResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_process_rpm_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportRPMResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportRPMResponse) ProtoMessage() {}

func (x *ImportRPMResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_process_rpm_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportRPMResponse.ProtoReflect.Descriptor instead.
func (*ImportRPMResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_process_rpm_proto_rawDescGZIP(), []int{5}
}

func (x *ImportRPMResponse) GetCommitHash() string {
	if x != nil {
		return x.CommitHash
	}
	return ""
}

func (x *ImportRPMResponse) GetCommitUri() string {
	if x != nil {
		return x.CommitUri
	}
	return ""
}

func (x *ImportRPMResponse) GetCommitBranch() string {
	if x != nil {
		return x.CommitBranch
	}
	return ""
}

func (x *ImportRPMResponse) GetCommitTag() string {
	if x != nil {
		return x.CommitTag
	}
	return ""
}

func (x *ImportRPMResponse) GetNevra() string {
	if x != nil {
		return x.Nevra
	}
	return ""
}

func (x *ImportRPMResponse) GetPkg() string {
	if x != nil {
		return x.Pkg
	}
	return ""
}

var File_proto_v1_process_rpm_proto protoreflect.FileDescriptor

var file_proto_v1_process_rpm_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x5f, 0x72, 0x70, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x6d, 0x6f,
	0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65,
	0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xb1, 0x01, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52,
	0x50, 0x4d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x07, 0x72, 0x70, 0x6d,
	0x5f, 0x75, 0x72, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52,
	0x06, 0x72, 0x70, 0x6d, 0x55, 0x72, 0x69, 0x12, 0x22, 0x0a, 0x0a, 0x6f, 0x73, 0x5f, 0x72, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02,
	0x52, 0x09, 0x6f, 0x73, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x08, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0,
	0x41, 0x02, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x12, 0x23, 0x0a, 0x0a,
	0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x22, 0x3d, 0x0a, 0x19, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x52, 0x50, 0x4d, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x08, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x22, 0xab, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x52, 0x50, 0x4d, 0x41, 0x72, 0x67, 0x73, 0x12, 0x3f, 0x0a, 0x07, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6d, 0x6f, 0x74,
	0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x52, 0x50, 0x4d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x03, 0xe0, 0x41,
	0x02, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x58, 0x0a, 0x10, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69,
	0x70, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x50, 0x4d, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x03,
	0xe0, 0x41, 0x02, 0x52, 0x0f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x86, 0x01, 0x0a, 0x12, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x52, 0x50, 0x4d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x39, 0x0a, 0x0a, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x40, 0x0a,
	0x12, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x50, 0x4d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x22,
	0xdd, 0x01, 0x0a, 0x11, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x50, 0x4d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52,
	0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x22, 0x0a, 0x0a, 0x63,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x03, 0xe0, 0x41, 0x02, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x55, 0x72, 0x69, 0x12,
	0x28, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x0c, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x22, 0x0a, 0x0a, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x5f, 0x74, 0x61, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0,
	0x41, 0x02, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x54, 0x61, 0x67, 0x12, 0x19, 0x0a,
	0x05, 0x6e, 0x65, 0x76, 0x72, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41,
	0x02, 0x52, 0x05, 0x6e, 0x65, 0x76, 0x72, 0x61, 0x12, 0x15, 0x0a, 0x03, 0x70, 0x6b, 0x67, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x03, 0x70, 0x6b, 0x67, 0x42,
	0x5e, 0x0a, 0x19, 0x6f, 0x72, 0x67, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x65, 0x6c, 0x61, 0x2e, 0x6d,
	0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x70, 0x6d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e,
	0x65, 0x6c, 0x61, 0x2f, 0x6d, 0x73, 0x68, 0x69, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x76, 0x31, 0x3b, 0x6d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_process_rpm_proto_rawDescOnce sync.Once
	file_proto_v1_process_rpm_proto_rawDescData = file_proto_v1_process_rpm_proto_rawDesc
)

func file_proto_v1_process_rpm_proto_rawDescGZIP() []byte {
	file_proto_v1_process_rpm_proto_rawDescOnce.Do(func() {
		file_proto_v1_process_rpm_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_process_rpm_proto_rawDescData)
	})
	return file_proto_v1_process_rpm_proto_rawDescData
}

var file_proto_v1_process_rpm_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_v1_process_rpm_proto_goTypes = []interface{}{
	(*ProcessRPMRequest)(nil),         // 0: mothership.v1.ProcessRPMRequest
	(*ProcessRPMInternalRequest)(nil), // 1: mothership.v1.ProcessRPMInternalRequest
	(*ProcessRPMArgs)(nil),            // 2: mothership.v1.ProcessRPMArgs
	(*ProcessRPMMetadata)(nil),        // 3: mothership.v1.ProcessRPMMetadata
	(*ProcessRPMResponse)(nil),        // 4: mothership.v1.ProcessRPMResponse
	(*ImportRPMResponse)(nil),         // 5: mothership.v1.ImportRPMResponse
	(*timestamppb.Timestamp)(nil),     // 6: google.protobuf.Timestamp
	(*Entry)(nil),                     // 7: mothership.v1.Entry
}
var file_proto_v1_process_rpm_proto_depIdxs = []int32{
	0, // 0: mothership.v1.ProcessRPMArgs.request:type_name -> mothership.v1.ProcessRPMRequest
	1, // 1: mothership.v1.ProcessRPMArgs.internal_request:type_name -> mothership.v1.ProcessRPMInternalRequest
	6, // 2: mothership.v1.ProcessRPMMetadata.start_time:type_name -> google.protobuf.Timestamp
	6, // 3: mothership.v1.ProcessRPMMetadata.end_time:type_name -> google.protobuf.Timestamp
	7, // 4: mothership.v1.ProcessRPMResponse.entry:type_name -> mothership.v1.Entry
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_v1_process_rpm_proto_init() }
func file_proto_v1_process_rpm_proto_init() {
	if File_proto_v1_process_rpm_proto != nil {
		return
	}
	file_proto_v1_entry_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_process_rpm_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRPMRequest); i {
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
		file_proto_v1_process_rpm_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRPMInternalRequest); i {
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
		file_proto_v1_process_rpm_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRPMArgs); i {
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
		file_proto_v1_process_rpm_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRPMMetadata); i {
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
		file_proto_v1_process_rpm_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRPMResponse); i {
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
		file_proto_v1_process_rpm_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImportRPMResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_v1_process_rpm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_v1_process_rpm_proto_goTypes,
		DependencyIndexes: file_proto_v1_process_rpm_proto_depIdxs,
		MessageInfos:      file_proto_v1_process_rpm_proto_msgTypes,
	}.Build()
	File_proto_v1_process_rpm_proto = out.File
	file_proto_v1_process_rpm_proto_rawDesc = nil
	file_proto_v1_process_rpm_proto_goTypes = nil
	file_proto_v1_process_rpm_proto_depIdxs = nil
}
