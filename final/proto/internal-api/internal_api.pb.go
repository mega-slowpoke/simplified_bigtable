// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.24.4
// source: internal_api.proto

package __

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

// Messages for master-server interactions
type AssignTabletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName      string `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
	TabletStartRow string `protobuf:"bytes,2,opt,name=tablet_start_row,json=tabletStartRow,proto3" json:"tablet_start_row,omitempty"`
	TabletEndRow   string `protobuf:"bytes,3,opt,name=tablet_end_row,json=tabletEndRow,proto3" json:"tablet_end_row,omitempty"`
	ServerAddress  string `protobuf:"bytes,4,opt,name=server_address,json=serverAddress,proto3" json:"server_address,omitempty"`
}

func (x *AssignTabletRequest) Reset() {
	*x = AssignTabletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignTabletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignTabletRequest) ProtoMessage() {}

func (x *AssignTabletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignTabletRequest.ProtoReflect.Descriptor instead.
func (*AssignTabletRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{0}
}

func (x *AssignTabletRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

func (x *AssignTabletRequest) GetTabletStartRow() string {
	if x != nil {
		return x.TabletStartRow
	}
	return ""
}

func (x *AssignTabletRequest) GetTabletEndRow() string {
	if x != nil {
		return x.TabletEndRow
	}
	return ""
}

func (x *AssignTabletRequest) GetServerAddress() string {
	if x != nil {
		return x.ServerAddress
	}
	return ""
}

type AssignTabletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success      bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}

func (x *AssignTabletResponse) Reset() {
	*x = AssignTabletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignTabletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignTabletResponse) ProtoMessage() {}

func (x *AssignTabletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignTabletResponse.ProtoReflect.Descriptor instead.
func (*AssignTabletResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{1}
}

func (x *AssignTabletResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AssignTabletResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type UnassignTabletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName      string `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
	TabletStartRow string `protobuf:"bytes,2,opt,name=tablet_start_row,json=tabletStartRow,proto3" json:"tablet_start_row,omitempty"`
	TabletEndRow   string `protobuf:"bytes,3,opt,name=tablet_end_row,json=tabletEndRow,proto3" json:"tablet_end_row,omitempty"`
}

func (x *UnassignTabletRequest) Reset() {
	*x = UnassignTabletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnassignTabletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnassignTabletRequest) ProtoMessage() {}

func (x *UnassignTabletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnassignTabletRequest.ProtoReflect.Descriptor instead.
func (*UnassignTabletRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{2}
}

func (x *UnassignTabletRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

func (x *UnassignTabletRequest) GetTabletStartRow() string {
	if x != nil {
		return x.TabletStartRow
	}
	return ""
}

func (x *UnassignTabletRequest) GetTabletEndRow() string {
	if x != nil {
		return x.TabletEndRow
	}
	return ""
}

type UnassignTabletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success      bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}

func (x *UnassignTabletResponse) Reset() {
	*x = UnassignTabletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnassignTabletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnassignTabletResponse) ProtoMessage() {}

func (x *UnassignTabletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnassignTabletResponse.ProtoReflect.Descriptor instead.
func (*UnassignTabletResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{3}
}

func (x *UnassignTabletResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UnassignTabletResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type GetTabletAssignmentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName string `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
}

func (x *GetTabletAssignmentsRequest) Reset() {
	*x = GetTabletAssignmentsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTabletAssignmentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTabletAssignmentsRequest) ProtoMessage() {}

func (x *GetTabletAssignmentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTabletAssignmentsRequest.ProtoReflect.Descriptor instead.
func (*GetTabletAssignmentsRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{4}
}

func (x *GetTabletAssignmentsRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

type GetTabletAssignmentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TabletAssignments []*TabletAssignment `protobuf:"bytes,1,rep,name=tablet_assignments,json=tabletAssignments,proto3" json:"tablet_assignments,omitempty"`
}

func (x *GetTabletAssignmentsResponse) Reset() {
	*x = GetTabletAssignmentsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTabletAssignmentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTabletAssignmentsResponse) ProtoMessage() {}

func (x *GetTabletAssignmentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTabletAssignmentsResponse.ProtoReflect.Descriptor instead.
func (*GetTabletAssignmentsResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{5}
}

func (x *GetTabletAssignmentsResponse) GetTabletAssignments() []*TabletAssignment {
	if x != nil {
		return x.TabletAssignments
	}
	return nil
}

type TabletAssignment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TabletStartRow string `protobuf:"bytes,1,opt,name=tablet_start_row,json=tabletStartRow,proto3" json:"tablet_start_row,omitempty"`
	TabletEndRow   string `protobuf:"bytes,2,opt,name=tablet_end_row,json=tabletEndRow,proto3" json:"tablet_end_row,omitempty"`
	ServerAddress  string `protobuf:"bytes,3,opt,name=server_address,json=serverAddress,proto3" json:"server_address,omitempty"`
}

func (x *TabletAssignment) Reset() {
	*x = TabletAssignment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TabletAssignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TabletAssignment) ProtoMessage() {}

func (x *TabletAssignment) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TabletAssignment.ProtoReflect.Descriptor instead.
func (*TabletAssignment) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{6}
}

func (x *TabletAssignment) GetTabletStartRow() string {
	if x != nil {
		return x.TabletStartRow
	}
	return ""
}

func (x *TabletAssignment) GetTabletEndRow() string {
	if x != nil {
		return x.TabletEndRow
	}
	return ""
}

func (x *TabletAssignment) GetServerAddress() string {
	if x != nil {
		return x.ServerAddress
	}
	return ""
}

// OPTIONAL IN MVP
type ReportTabletStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName      string `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
	TabletStartRow string `protobuf:"bytes,2,opt,name=tablet_start_row,json=tabletStartRow,proto3" json:"tablet_start_row,omitempty"`
	TabletEndRow   string `protobuf:"bytes,3,opt,name=tablet_end_row,json=tabletEndRow,proto3" json:"tablet_end_row,omitempty"`
	ServerAddress  string `protobuf:"bytes,4,opt,name=server_address,json=serverAddress,proto3" json:"server_address,omitempty"`
	Status         string `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"` // e.g., "HEALTHY", "UNRESPONSIVE", etc.
}

func (x *ReportTabletStatusRequest) Reset() {
	*x = ReportTabletStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportTabletStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportTabletStatusRequest) ProtoMessage() {}

func (x *ReportTabletStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportTabletStatusRequest.ProtoReflect.Descriptor instead.
func (*ReportTabletStatusRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{7}
}

func (x *ReportTabletStatusRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

func (x *ReportTabletStatusRequest) GetTabletStartRow() string {
	if x != nil {
		return x.TabletStartRow
	}
	return ""
}

func (x *ReportTabletStatusRequest) GetTabletEndRow() string {
	if x != nil {
		return x.TabletEndRow
	}
	return ""
}

func (x *ReportTabletStatusRequest) GetServerAddress() string {
	if x != nil {
		return x.ServerAddress
	}
	return ""
}

func (x *ReportTabletStatusRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type ReportTabletStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success      bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}

func (x *ReportTabletStatusResponse) Reset() {
	*x = ReportTabletStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportTabletStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportTabletStatusResponse) ProtoMessage() {}

func (x *ReportTabletStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportTabletStatusResponse.ProtoReflect.Descriptor instead.
func (*ReportTabletStatusResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_proto_rawDescGZIP(), []int{8}
}

func (x *ReportTabletStatusResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ReportTabletStatusResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

var File_internal_api_proto protoreflect.FileDescriptor

var file_internal_api_proto_rawDesc = []byte{
	0x0a, 0x12, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x22, 0xab,
	0x01, 0x0a, 0x13, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x6f, 0x77, 0x12,
	0x24, 0x0a, 0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x6f,
	0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x45,
	0x6e, 0x64, 0x52, 0x6f, 0x77, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x55, 0x0a, 0x14,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x23,
	0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x86, 0x01, 0x0a, 0x15, 0x55, 0x6e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x72, 0x6f, 0x77,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x52, 0x6f, 0x77, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74,
	0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x45, 0x6e, 0x64, 0x52, 0x6f, 0x77, 0x22, 0x57, 0x0a, 0x16,
	0x55, 0x6e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3c, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x69, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x12, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x61, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x11, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x89,
	0x01, 0x0a, 0x10, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x6f, 0x77, 0x12, 0x24, 0x0a,
	0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x6f, 0x77, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x45, 0x6e, 0x64,
	0x52, 0x6f, 0x77, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0xc9, 0x01, 0x0a, 0x19, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x74, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x6f,
	0x77, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x65, 0x6e, 0x64, 0x5f,
	0x72, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x74, 0x45, 0x6e, 0x64, 0x52, 0x6f, 0x77, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x5b, 0x0a, 0x1a, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x23,
	0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x32, 0xfb, 0x02, 0x0a, 0x0d, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x0c, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x74, 0x12, 0x1d, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x0e, 0x55, 0x6e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x12, 0x1f, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x2e, 0x55, 0x6e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x2e, 0x55, 0x6e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x14, 0x47, 0x65, 0x74,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x25, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73, 0x73,
	0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x5f, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x62, 0x69,
	0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x04, 0x5a, 0x02, 0x2f, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_api_proto_rawDescOnce sync.Once
	file_internal_api_proto_rawDescData = file_internal_api_proto_rawDesc
)

func file_internal_api_proto_rawDescGZIP() []byte {
	file_internal_api_proto_rawDescOnce.Do(func() {
		file_internal_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_api_proto_rawDescData)
	})
	return file_internal_api_proto_rawDescData
}

var file_internal_api_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_internal_api_proto_goTypes = []any{
	(*AssignTabletRequest)(nil),          // 0: bigtable.AssignTabletRequest
	(*AssignTabletResponse)(nil),         // 1: bigtable.AssignTabletResponse
	(*UnassignTabletRequest)(nil),        // 2: bigtable.UnassignTabletRequest
	(*UnassignTabletResponse)(nil),       // 3: bigtable.UnassignTabletResponse
	(*GetTabletAssignmentsRequest)(nil),  // 4: bigtable.GetTabletAssignmentsRequest
	(*GetTabletAssignmentsResponse)(nil), // 5: bigtable.GetTabletAssignmentsResponse
	(*TabletAssignment)(nil),             // 6: bigtable.TabletAssignment
	(*ReportTabletStatusRequest)(nil),    // 7: bigtable.ReportTabletStatusRequest
	(*ReportTabletStatusResponse)(nil),   // 8: bigtable.ReportTabletStatusResponse
}
var file_internal_api_proto_depIdxs = []int32{
	6, // 0: bigtable.GetTabletAssignmentsResponse.tablet_assignments:type_name -> bigtable.TabletAssignment
	0, // 1: bigtable.MasterService.AssignTablet:input_type -> bigtable.AssignTabletRequest
	2, // 2: bigtable.MasterService.UnassignTablet:input_type -> bigtable.UnassignTabletRequest
	4, // 3: bigtable.MasterService.GetTabletAssignments:input_type -> bigtable.GetTabletAssignmentsRequest
	7, // 4: bigtable.MasterService.ReportTabletStatus:input_type -> bigtable.ReportTabletStatusRequest
	1, // 5: bigtable.MasterService.AssignTablet:output_type -> bigtable.AssignTabletResponse
	3, // 6: bigtable.MasterService.UnassignTablet:output_type -> bigtable.UnassignTabletResponse
	5, // 7: bigtable.MasterService.GetTabletAssignments:output_type -> bigtable.GetTabletAssignmentsResponse
	8, // 8: bigtable.MasterService.ReportTabletStatus:output_type -> bigtable.ReportTabletStatusResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_api_proto_init() }
func file_internal_api_proto_init() {
	if File_internal_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_api_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*AssignTabletRequest); i {
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
		file_internal_api_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AssignTabletResponse); i {
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
		file_internal_api_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*UnassignTabletRequest); i {
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
		file_internal_api_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UnassignTabletResponse); i {
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
		file_internal_api_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*GetTabletAssignmentsRequest); i {
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
		file_internal_api_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*GetTabletAssignmentsResponse); i {
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
		file_internal_api_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*TabletAssignment); i {
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
		file_internal_api_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*ReportTabletStatusRequest); i {
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
		file_internal_api_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*ReportTabletStatusResponse); i {
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
			RawDescriptor: file_internal_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_api_proto_goTypes,
		DependencyIndexes: file_internal_api_proto_depIdxs,
		MessageInfos:      file_internal_api_proto_msgTypes,
	}.Build()
	File_internal_api_proto = out.File
	file_internal_api_proto_rawDesc = nil
	file_internal_api_proto_goTypes = nil
	file_internal_api_proto_depIdxs = nil
}
