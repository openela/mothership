// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.3
// source: proto/admin/v1/bugtracker.proto

package mshipadminpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Supported bug trackers.
type BugTrackerConfig_Type int32

const (
	// Unknown bug tracker.
	BugTrackerConfig_UNKNOWN BugTrackerConfig_Type = 0
	// MantisBT bug tracker.
	BugTrackerConfig_MANTIS BugTrackerConfig_Type = 1
)

// Enum value maps for BugTrackerConfig_Type.
var (
	BugTrackerConfig_Type_name = map[int32]string{
		0: "UNKNOWN",
		1: "MANTIS",
	}
	BugTrackerConfig_Type_value = map[string]int32{
		"UNKNOWN": 0,
		"MANTIS":  1,
	}
)

func (x BugTrackerConfig_Type) Enum() *BugTrackerConfig_Type {
	p := new(BugTrackerConfig_Type)
	*p = x
	return p
}

func (x BugTrackerConfig_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BugTrackerConfig_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_admin_v1_bugtracker_proto_enumTypes[0].Descriptor()
}

func (BugTrackerConfig_Type) Type() protoreflect.EnumType {
	return &file_proto_admin_v1_bugtracker_proto_enumTypes[0]
}

func (x BugTrackerConfig_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BugTrackerConfig_Type.Descriptor instead.
func (BugTrackerConfig_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_admin_v1_bugtracker_proto_rawDescGZIP(), []int{0, 0}
}

// BugTrackerConfig is the configuration for a bug tracker.
// Usually, the bug tracker is a third-party service, such as Mantis or Track
// The configuration is used to track import batches
type BugTrackerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Type of the bug tracker.
	Type BugTrackerConfig_Type `protobuf:"varint,1,opt,name=type,proto3,enum=mothership.admin.v1.BugTrackerConfig_Type" json:"type,omitempty"`
	// URI of the bug tracker.
	Uri string `protobuf:"bytes,2,opt,name=uri,proto3" json:"uri,omitempty"`
	// Configuration for the bug tracker.
	//
	// Types that are assignable to Config:
	//
	//	*BugTrackerConfig_Mantis
	Config isBugTrackerConfig_Config `protobuf_oneof:"config"`
}

func (x *BugTrackerConfig) Reset() {
	*x = BugTrackerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_admin_v1_bugtracker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BugTrackerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BugTrackerConfig) ProtoMessage() {}

func (x *BugTrackerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_proto_admin_v1_bugtracker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BugTrackerConfig.ProtoReflect.Descriptor instead.
func (*BugTrackerConfig) Descriptor() ([]byte, []int) {
	return file_proto_admin_v1_bugtracker_proto_rawDescGZIP(), []int{0}
}

func (x *BugTrackerConfig) GetType() BugTrackerConfig_Type {
	if x != nil {
		return x.Type
	}
	return BugTrackerConfig_UNKNOWN
}

func (x *BugTrackerConfig) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

func (m *BugTrackerConfig) GetConfig() isBugTrackerConfig_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (x *BugTrackerConfig) GetMantis() *BugTrackerConfig_MantisConfig {
	if x, ok := x.GetConfig().(*BugTrackerConfig_Mantis); ok {
		return x.Mantis
	}
	return nil
}

type isBugTrackerConfig_Config interface {
	isBugTrackerConfig_Config()
}

type BugTrackerConfig_Mantis struct {
	// User-defined configuration for MantisBT.
	Mantis *BugTrackerConfig_MantisConfig `protobuf:"bytes,3,opt,name=mantis,proto3,oneof"`
}

func (*BugTrackerConfig_Mantis) isBugTrackerConfig_Config() {}

// Configuration options for MantisBT
type BugTrackerConfig_MantisConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// API key for the bug tracker.
	ApiKey string `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	// Project ID mapping.
	// Maps major version to project ID.
	ProjectIds map[int32]int64 `protobuf:"bytes,2,rep,name=project_ids,json=projectIds,proto3" json:"project_ids,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *BugTrackerConfig_MantisConfig) Reset() {
	*x = BugTrackerConfig_MantisConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_admin_v1_bugtracker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BugTrackerConfig_MantisConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BugTrackerConfig_MantisConfig) ProtoMessage() {}

func (x *BugTrackerConfig_MantisConfig) ProtoReflect() protoreflect.Message {
	mi := &file_proto_admin_v1_bugtracker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BugTrackerConfig_MantisConfig.ProtoReflect.Descriptor instead.
func (*BugTrackerConfig_MantisConfig) Descriptor() ([]byte, []int) {
	return file_proto_admin_v1_bugtracker_proto_rawDescGZIP(), []int{0, 0}
}

func (x *BugTrackerConfig_MantisConfig) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *BugTrackerConfig_MantisConfig) GetProjectIds() map[int32]int64 {
	if x != nil {
		return x.ProjectIds
	}
	return nil
}

var File_proto_admin_v1_bugtracker_proto protoreflect.FileDescriptor

var file_proto_admin_v1_bugtracker_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31,
	0x2f, 0x62, 0x75, 0x67, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x13, 0x6d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x22, 0xab, 0x03, 0x0a, 0x10, 0x42, 0x75, 0x67, 0x54, 0x72,
	0x61, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3e, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x6d, 0x6f, 0x74, 0x68,
	0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x42, 0x75, 0x67, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x12, 0x4c, 0x0a,
	0x06, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e,
	0x6d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x67, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4d, 0x61, 0x6e, 0x74, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x48, 0x00, 0x52, 0x06, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x73, 0x1a, 0xcb, 0x01, 0x0a, 0x0c,
	0x4d, 0x61, 0x6e, 0x74, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x17, 0x0a, 0x07,
	0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x70, 0x69, 0x4b, 0x65, 0x79, 0x12, 0x63, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x6d, 0x6f, 0x74,
	0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x42, 0x75, 0x67, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x4d, 0x61, 0x6e, 0x74, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x73, 0x1a, 0x3d, 0x0a, 0x0f, 0x50, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x1f, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x4d, 0x41, 0x4e, 0x54, 0x49, 0x53, 0x10, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x42, 0x6a, 0x0a, 0x1f, 0x6f, 0x72, 0x67, 0x2e, 0x6f, 0x70, 0x65, 0x6e,
	0x65, 0x6c, 0x61, 0x2e, 0x6d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x42, 0x75, 0x67, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x65, 0x6c, 0x61, 0x2f, 0x6d,
	0x73, 0x68, 0x69, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2f, 0x76, 0x31, 0x3b, 0x6d, 0x73, 0x68, 0x69, 0x70, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_admin_v1_bugtracker_proto_rawDescOnce sync.Once
	file_proto_admin_v1_bugtracker_proto_rawDescData = file_proto_admin_v1_bugtracker_proto_rawDesc
)

func file_proto_admin_v1_bugtracker_proto_rawDescGZIP() []byte {
	file_proto_admin_v1_bugtracker_proto_rawDescOnce.Do(func() {
		file_proto_admin_v1_bugtracker_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_admin_v1_bugtracker_proto_rawDescData)
	})
	return file_proto_admin_v1_bugtracker_proto_rawDescData
}

var file_proto_admin_v1_bugtracker_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_admin_v1_bugtracker_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_admin_v1_bugtracker_proto_goTypes = []interface{}{
	(BugTrackerConfig_Type)(0),            // 0: mothership.admin.v1.BugTrackerConfig.Type
	(*BugTrackerConfig)(nil),              // 1: mothership.admin.v1.BugTrackerConfig
	(*BugTrackerConfig_MantisConfig)(nil), // 2: mothership.admin.v1.BugTrackerConfig.MantisConfig
	nil,                                   // 3: mothership.admin.v1.BugTrackerConfig.MantisConfig.ProjectIdsEntry
}
var file_proto_admin_v1_bugtracker_proto_depIdxs = []int32{
	0, // 0: mothership.admin.v1.BugTrackerConfig.type:type_name -> mothership.admin.v1.BugTrackerConfig.Type
	2, // 1: mothership.admin.v1.BugTrackerConfig.mantis:type_name -> mothership.admin.v1.BugTrackerConfig.MantisConfig
	3, // 2: mothership.admin.v1.BugTrackerConfig.MantisConfig.project_ids:type_name -> mothership.admin.v1.BugTrackerConfig.MantisConfig.ProjectIdsEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_admin_v1_bugtracker_proto_init() }
func file_proto_admin_v1_bugtracker_proto_init() {
	if File_proto_admin_v1_bugtracker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_admin_v1_bugtracker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BugTrackerConfig); i {
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
		file_proto_admin_v1_bugtracker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BugTrackerConfig_MantisConfig); i {
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
	file_proto_admin_v1_bugtracker_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*BugTrackerConfig_Mantis)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_admin_v1_bugtracker_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_admin_v1_bugtracker_proto_goTypes,
		DependencyIndexes: file_proto_admin_v1_bugtracker_proto_depIdxs,
		EnumInfos:         file_proto_admin_v1_bugtracker_proto_enumTypes,
		MessageInfos:      file_proto_admin_v1_bugtracker_proto_msgTypes,
	}.Build()
	File_proto_admin_v1_bugtracker_proto = out.File
	file_proto_admin_v1_bugtracker_proto_rawDesc = nil
	file_proto_admin_v1_bugtracker_proto_goTypes = nil
	file_proto_admin_v1_bugtracker_proto_depIdxs = nil
}
