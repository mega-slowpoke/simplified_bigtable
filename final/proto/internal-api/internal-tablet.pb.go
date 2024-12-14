// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.24.4
// source: internal-tablet.proto

package proto

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

type CreateTableInternalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName      string                                     `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
	ColumnFamilies []*CreateTableInternalRequest_ColumnFamily `protobuf:"bytes,2,rep,name=column_families,json=columnFamilies,proto3" json:"column_families,omitempty"`
}

func (x *CreateTableInternalRequest) Reset() {
	*x = CreateTableInternalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tablet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTableInternalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTableInternalRequest) ProtoMessage() {}

func (x *CreateTableInternalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tablet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTableInternalRequest.ProtoReflect.Descriptor instead.
func (*CreateTableInternalRequest) Descriptor() ([]byte, []int) {
	return file_internal_tablet_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTableInternalRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

func (x *CreateTableInternalRequest) GetColumnFamilies() []*CreateTableInternalRequest_ColumnFamily {
	if x != nil {
		return x.ColumnFamilies
	}
	return nil
}

type CreateTableInternalResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *CreateTableInternalResponse) Reset() {
	*x = CreateTableInternalResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tablet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTableInternalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTableInternalResponse) ProtoMessage() {}

func (x *CreateTableInternalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tablet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTableInternalResponse.ProtoReflect.Descriptor instead.
func (*CreateTableInternalResponse) Descriptor() ([]byte, []int) {
	return file_internal_tablet_proto_rawDescGZIP(), []int{1}
}

func (x *CreateTableInternalResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CreateTableInternalResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type DeleteTableInternalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName string `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
}

func (x *DeleteTableInternalRequest) Reset() {
	*x = DeleteTableInternalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tablet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTableInternalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTableInternalRequest) ProtoMessage() {}

func (x *DeleteTableInternalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tablet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTableInternalRequest.ProtoReflect.Descriptor instead.
func (*DeleteTableInternalRequest) Descriptor() ([]byte, []int) {
	return file_internal_tablet_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteTableInternalRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

type DeleteTableInternalResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *DeleteTableInternalResponse) Reset() {
	*x = DeleteTableInternalResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tablet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTableInternalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTableInternalResponse) ProtoMessage() {}

func (x *DeleteTableInternalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tablet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTableInternalResponse.ProtoReflect.Descriptor instead.
func (*DeleteTableInternalResponse) Descriptor() ([]byte, []int) {
	return file_internal_tablet_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteTableInternalResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *DeleteTableInternalResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type CreateTableInternalRequest_ColumnFamily struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FamilyName string   `protobuf:"bytes,1,opt,name=family_name,json=familyName,proto3" json:"family_name,omitempty"`
	Columns    []string `protobuf:"bytes,2,rep,name=columns,proto3" json:"columns,omitempty"`
}

func (x *CreateTableInternalRequest_ColumnFamily) Reset() {
	*x = CreateTableInternalRequest_ColumnFamily{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tablet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTableInternalRequest_ColumnFamily) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTableInternalRequest_ColumnFamily) ProtoMessage() {}

func (x *CreateTableInternalRequest_ColumnFamily) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tablet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTableInternalRequest_ColumnFamily.ProtoReflect.Descriptor instead.
func (*CreateTableInternalRequest_ColumnFamily) Descriptor() ([]byte, []int) {
	return file_internal_tablet_proto_rawDescGZIP(), []int{0, 0}
}

func (x *CreateTableInternalRequest_ColumnFamily) GetFamilyName() string {
	if x != nil {
		return x.FamilyName
	}
	return ""
}

func (x *CreateTableInternalRequest_ColumnFamily) GetColumns() []string {
	if x != nil {
		return x.Columns
	}
	return nil
}

var File_internal_tablet_proto protoreflect.FileDescriptor

var file_internal_tablet_proto_rawDesc = []byte{
	0x0a, 0x15, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2d, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x22, 0xe2, 0x01, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x5a, 0x0a, 0x0f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x5f, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x69,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43,
	0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x52, 0x0e, 0x63, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x69, 0x65, 0x73, 0x1a, 0x49, 0x0a, 0x0c, 0x43,
	0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x66,
	0x61, 0x6d, 0x69, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x22, 0x51, 0x0a, 0x1b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3b, 0x0a, 0x1a, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x51, 0x0a, 0x1b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xcf, 0x01, 0x0a, 0x15, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x12, 0x24, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x5a, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x24,
	0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x66,
	0x69, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_internal_tablet_proto_rawDescOnce sync.Once
	file_internal_tablet_proto_rawDescData = file_internal_tablet_proto_rawDesc
)

func file_internal_tablet_proto_rawDescGZIP() []byte {
	file_internal_tablet_proto_rawDescOnce.Do(func() {
		file_internal_tablet_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_tablet_proto_rawDescData)
	})
	return file_internal_tablet_proto_rawDescData
}

var file_internal_tablet_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_internal_tablet_proto_goTypes = []any{
	(*CreateTableInternalRequest)(nil),              // 0: bigtable.CreateTableInternalRequest
	(*CreateTableInternalResponse)(nil),             // 1: bigtable.CreateTableInternalResponse
	(*DeleteTableInternalRequest)(nil),              // 2: bigtable.DeleteTableInternalRequest
	(*DeleteTableInternalResponse)(nil),             // 3: bigtable.DeleteTableInternalResponse
	(*CreateTableInternalRequest_ColumnFamily)(nil), // 4: bigtable.CreateTableInternalRequest.ColumnFamily
}
var file_internal_tablet_proto_depIdxs = []int32{
	4, // 0: bigtable.CreateTableInternalRequest.column_families:type_name -> bigtable.CreateTableInternalRequest.ColumnFamily
	0, // 1: bigtable.TabletInternalService.CreateTable:input_type -> bigtable.CreateTableInternalRequest
	2, // 2: bigtable.TabletInternalService.DeleteTable:input_type -> bigtable.DeleteTableInternalRequest
	1, // 3: bigtable.TabletInternalService.CreateTable:output_type -> bigtable.CreateTableInternalResponse
	3, // 4: bigtable.TabletInternalService.DeleteTable:output_type -> bigtable.DeleteTableInternalResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_tablet_proto_init() }
func file_internal_tablet_proto_init() {
	if File_internal_tablet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_tablet_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTableInternalRequest); i {
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
		file_internal_tablet_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTableInternalResponse); i {
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
		file_internal_tablet_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteTableInternalRequest); i {
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
		file_internal_tablet_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteTableInternalResponse); i {
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
		file_internal_tablet_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTableInternalRequest_ColumnFamily); i {
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
			RawDescriptor: file_internal_tablet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_tablet_proto_goTypes,
		DependencyIndexes: file_internal_tablet_proto_depIdxs,
		MessageInfos:      file_internal_tablet_proto_msgTypes,
	}.Build()
	File_internal_tablet_proto = out.File
	file_internal_tablet_proto_rawDesc = nil
	file_internal_tablet_proto_goTypes = nil
	file_internal_tablet_proto_depIdxs = nil
}
