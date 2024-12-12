// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.24.4
// source: master.proto

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

// Messages for master-server interactions
// 现在我们先假设每个tablet server只能管理一个table，样register里面就可以直接放table name, table_start_row, table_end_row
// 如果一个tablet 还可以管理多个table，那么这里需要传输一个列表的table name, table_start_row, table_end_row
type RegisterTabletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName     string `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
	TableStartRow string `protobuf:"bytes,2,opt,name=table_start_row,json=tableStartRow,proto3" json:"table_start_row,omitempty"`
	TableEndRow   string `protobuf:"bytes,3,opt,name=table_end_row,json=tableEndRow,proto3" json:"table_end_row,omitempty"`
	TabletAddress string `protobuf:"bytes,4,opt,name=tablet_address,json=tabletAddress,proto3" json:"tablet_address,omitempty"`
}

func (x *RegisterTabletRequest) Reset() {
	*x = RegisterTabletRequest{}
	mi := &file_master_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterTabletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterTabletRequest) ProtoMessage() {}

func (x *RegisterTabletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterTabletRequest.ProtoReflect.Descriptor instead.
func (*RegisterTabletRequest) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterTabletRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

func (x *RegisterTabletRequest) GetTableStartRow() string {
	if x != nil {
		return x.TableStartRow
	}
	return ""
}

func (x *RegisterTabletRequest) GetTableEndRow() string {
	if x != nil {
		return x.TableEndRow
	}
	return ""
}

func (x *RegisterTabletRequest) GetTabletAddress() string {
	if x != nil {
		return x.TabletAddress
	}
	return ""
}

type RegisterTabletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success      bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}

func (x *RegisterTabletResponse) Reset() {
	*x = RegisterTabletResponse{}
	mi := &file_master_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterTabletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterTabletResponse) ProtoMessage() {}

func (x *RegisterTabletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterTabletResponse.ProtoReflect.Descriptor instead.
func (*RegisterTabletResponse) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterTabletResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *RegisterTabletResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type UnregisterTabletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName string `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
}

func (x *UnregisterTabletRequest) Reset() {
	*x = UnregisterTabletRequest{}
	mi := &file_master_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnregisterTabletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterTabletRequest) ProtoMessage() {}

func (x *UnregisterTabletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterTabletRequest.ProtoReflect.Descriptor instead.
func (*UnregisterTabletRequest) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{2}
}

func (x *UnregisterTabletRequest) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

type UnregisterTabletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success      bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}

func (x *UnregisterTabletResponse) Reset() {
	*x = UnregisterTabletResponse{}
	mi := &file_master_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnregisterTabletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterTabletResponse) ProtoMessage() {}

func (x *UnregisterTabletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterTabletResponse.ProtoReflect.Descriptor instead.
func (*UnregisterTabletResponse) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{3}
}

func (x *UnregisterTabletResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UnregisterTabletResponse) GetErrorMessage() string {
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
	mi := &file_master_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTabletAssignmentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTabletAssignmentsRequest) ProtoMessage() {}

func (x *GetTabletAssignmentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[4]
	if x != nil {
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
	return file_master_proto_rawDescGZIP(), []int{4}
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
	mi := &file_master_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTabletAssignmentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTabletAssignmentsResponse) ProtoMessage() {}

func (x *GetTabletAssignmentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[5]
	if x != nil {
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
	return file_master_proto_rawDescGZIP(), []int{5}
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
	mi := &file_master_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TabletAssignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TabletAssignment) ProtoMessage() {}

func (x *TabletAssignment) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[6]
	if x != nil {
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
	return file_master_proto_rawDescGZIP(), []int{6}
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
	mi := &file_master_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReportTabletStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportTabletStatusRequest) ProtoMessage() {}

func (x *ReportTabletStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[7]
	if x != nil {
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
	return file_master_proto_rawDescGZIP(), []int{7}
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
	mi := &file_master_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReportTabletStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportTabletStatusResponse) ProtoMessage() {}

func (x *ReportTabletStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[8]
	if x != nil {
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
	return file_master_proto_rawDescGZIP(), []int{8}
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

var File_master_proto protoreflect.FileDescriptor

var file_master_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x22, 0xa9, 0x01, 0x0a, 0x15, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x26, 0x0a, 0x0f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x72, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x6f, 0x77, 0x12, 0x22, 0x0a, 0x0d, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x45, 0x6e, 0x64, 0x52, 0x6f, 0x77, 0x12, 0x25, 0x0a,
	0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x22, 0x57, 0x0a, 0x16, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x38, 0x0a,
	0x17, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x59, 0x0a, 0x18, 0x55, 0x6e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x23, 0x0a,
	0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x3c, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x69, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73, 0x73,
	0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x49, 0x0a, 0x12, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x61, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x62,
	0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x11, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x89, 0x01, 0x0a, 0x10,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x72, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x6f, 0x77, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x74, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x45, 0x6e, 0x64, 0x52, 0x6f, 0x77,
	0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0xc9, 0x01, 0x0a, 0x19, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x6f, 0x77, 0x12, 0x24,
	0x0a, 0x0e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x6f, 0x77,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x45, 0x6e,
	0x64, 0x52, 0x6f, 0x77, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x5b, 0x0a, 0x1a, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0x86, 0x03, 0x0a, 0x0d, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x53, 0x0a, 0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x74, 0x12, 0x1f, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58, 0x0a, 0x10, 0x55, 0x6e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x12, 0x21, 0x2e, 0x62, 0x69,
	0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x65, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x25, 0x2e, 0x62, 0x69, 0x67, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x26, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x74, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x23,
	0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x62, 0x69, 0x67, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x66, 0x69, 0x6e,
	0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_master_proto_rawDescOnce sync.Once
	file_master_proto_rawDescData = file_master_proto_rawDesc
)

func file_master_proto_rawDescGZIP() []byte {
	file_master_proto_rawDescOnce.Do(func() {
		file_master_proto_rawDescData = protoimpl.X.CompressGZIP(file_master_proto_rawDescData)
	})
	return file_master_proto_rawDescData
}

var file_master_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_master_proto_goTypes = []any{
	(*RegisterTabletRequest)(nil),        // 0: bigtable.RegisterTabletRequest
	(*RegisterTabletResponse)(nil),       // 1: bigtable.RegisterTabletResponse
	(*UnregisterTabletRequest)(nil),      // 2: bigtable.UnregisterTabletRequest
	(*UnregisterTabletResponse)(nil),     // 3: bigtable.UnregisterTabletResponse
	(*GetTabletAssignmentsRequest)(nil),  // 4: bigtable.GetTabletAssignmentsRequest
	(*GetTabletAssignmentsResponse)(nil), // 5: bigtable.GetTabletAssignmentsResponse
	(*TabletAssignment)(nil),             // 6: bigtable.TabletAssignment
	(*ReportTabletStatusRequest)(nil),    // 7: bigtable.ReportTabletStatusRequest
	(*ReportTabletStatusResponse)(nil),   // 8: bigtable.ReportTabletStatusResponse
}
var file_master_proto_depIdxs = []int32{
	6, // 0: bigtable.GetTabletAssignmentsResponse.tablet_assignments:type_name -> bigtable.TabletAssignment
	0, // 1: bigtable.MasterService.RegisterTablet:input_type -> bigtable.RegisterTabletRequest
	2, // 2: bigtable.MasterService.UnregisterTablet:input_type -> bigtable.UnregisterTabletRequest
	4, // 3: bigtable.MasterService.GetTabletAssignments:input_type -> bigtable.GetTabletAssignmentsRequest
	7, // 4: bigtable.MasterService.ReportTabletStatus:input_type -> bigtable.ReportTabletStatusRequest
	1, // 5: bigtable.MasterService.RegisterTablet:output_type -> bigtable.RegisterTabletResponse
	2, // 6: bigtable.MasterService.UnregisterTablet:output_type -> bigtable.UnregisterTabletRequest
	5, // 7: bigtable.MasterService.GetTabletAssignments:output_type -> bigtable.GetTabletAssignmentsResponse
	8, // 8: bigtable.MasterService.ReportTabletStatus:output_type -> bigtable.ReportTabletStatusResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_master_proto_init() }
func file_master_proto_init() {
	if File_master_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_master_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_master_proto_goTypes,
		DependencyIndexes: file_master_proto_depIdxs,
		MessageInfos:      file_master_proto_msgTypes,
	}.Build()
	File_master_proto = out.File
	file_master_proto_rawDesc = nil
	file_master_proto_goTypes = nil
	file_master_proto_depIdxs = nil
}
