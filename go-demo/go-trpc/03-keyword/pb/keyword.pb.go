// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: keyword.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ResponseCode int32

const (
	ResponseCode_OK            ResponseCode = 0
	ResponseCode_FAIL          ResponseCode = 1
	ResponseCode_INVALID_PARAM ResponseCode = 2
)

// Enum value maps for ResponseCode.
var (
	ResponseCode_name = map[int32]string{
		0: "OK",
		1: "FAIL",
		2: "INVALID_PARAM",
	}
	ResponseCode_value = map[string]int32{
		"OK":            0,
		"FAIL":          1,
		"INVALID_PARAM": 2,
	}
)

func (x ResponseCode) Enum() *ResponseCode {
	p := new(ResponseCode)
	*p = x
	return p
}

func (x ResponseCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseCode) Descriptor() protoreflect.EnumDescriptor {
	return file_keyword_proto_enumTypes[0].Descriptor()
}

func (ResponseCode) Type() protoreflect.EnumType {
	return &file_keyword_proto_enumTypes[0]
}

func (x ResponseCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseCode.Descriptor instead.
func (ResponseCode) EnumDescriptor() ([]byte, []int) {
	return file_keyword_proto_rawDescGZIP(), []int{0}
}

// message： 定义结构体，类比go中的type
type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// optional： 可选字段
	ReqCreateTime *string           `protobuf:"bytes,1,opt,name=reqCreateTime,proto3,oneof" json:"reqCreateTime,omitempty"`
	ReqInfo       map[string]string `protobuf:"bytes,2,rep,name=reqInfo,proto3" json:"reqInfo,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_keyword_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_keyword_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_keyword_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetReqCreateTime() string {
	if x != nil && x.ReqCreateTime != nil {
		return *x.ReqCreateTime
	}
	return ""
}

func (x *Request) GetReqInfo() map[string]string {
	if x != nil {
		return x.ReqInfo
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RspInfo map[string]string `protobuf:"bytes,1,rep,name=rspInfo,proto3" json:"rspInfo,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_keyword_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_keyword_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_keyword_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetRspInfo() map[string]string {
	if x != nil {
		return x.RspInfo
	}
	return nil
}

type Classroom struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// repeated 列表（切片）
	StudentIds []int32 `protobuf:"varint,2,rep,packed,name=studentIds,proto3" json:"studentIds,omitempty"`
}

func (x *Classroom) Reset() {
	*x = Classroom{}
	if protoimpl.UnsafeEnabled {
		mi := &file_keyword_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Classroom) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Classroom) ProtoMessage() {}

func (x *Classroom) ProtoReflect() protoreflect.Message {
	mi := &file_keyword_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Classroom.ProtoReflect.Descriptor instead.
func (*Classroom) Descriptor() ([]byte, []int) {
	return file_keyword_proto_rawDescGZIP(), []int{2}
}

func (x *Classroom) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Classroom) GetStudentIds() []int32 {
	if x != nil {
		return x.StudentIds
	}
	return nil
}

var File_keyword_proto protoreflect.FileDescriptor

var file_keyword_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0c, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x22, 0xc0, 0x01,
	0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x0d, 0x72, 0x65, 0x71,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x3c, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x49, 0x6e, 0x66, 0x6f, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x72, 0x65, 0x71, 0x49, 0x6e,
	0x66, 0x6f, 0x1a, 0x3a, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x10,
	0x0a, 0x0e, 0x5f, 0x72, 0x65, 0x71, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x22, 0x85, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a,
	0x07, 0x72, 0x73, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x73, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x07, 0x72, 0x73, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x3a, 0x0a, 0x0c,
	0x52, 0x73, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x3f, 0x0a, 0x09, 0x43, 0x6c, 0x61, 0x73,
	0x73, 0x72, 0x6f, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0a, 0x73,
	0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x2a, 0x33, 0x0a, 0x0c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10,
	0x00, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x49,
	0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x10, 0x02, 0x32, 0x4d,
	0x0a, 0x0e, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x3b, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x15,
	0x2e, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x6b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x28, 0x5a,
	0x26, 0x7a, 0x69, 0x79, 0x69, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x64, 0x65, 0x6d,
	0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x74, 0x72, 0x70, 0x63, 0x2f, 0x30, 0x33, 0x2d, 0x6b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_keyword_proto_rawDescOnce sync.Once
	file_keyword_proto_rawDescData = file_keyword_proto_rawDesc
)

func file_keyword_proto_rawDescGZIP() []byte {
	file_keyword_proto_rawDescOnce.Do(func() {
		file_keyword_proto_rawDescData = protoimpl.X.CompressGZIP(file_keyword_proto_rawDescData)
	})
	return file_keyword_proto_rawDescData
}

var file_keyword_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_keyword_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_keyword_proto_goTypes = []any{
	(ResponseCode)(0), // 0: trpc.keyword.ResponseCode
	(*Request)(nil),   // 1: trpc.keyword.Request
	(*Response)(nil),  // 2: trpc.keyword.Response
	(*Classroom)(nil), // 3: trpc.keyword.Classroom
	nil,               // 4: trpc.keyword.Request.ReqInfoEntry
	nil,               // 5: trpc.keyword.Response.RspInfoEntry
}
var file_keyword_proto_depIdxs = []int32{
	4, // 0: trpc.keyword.Request.reqInfo:type_name -> trpc.keyword.Request.ReqInfoEntry
	5, // 1: trpc.keyword.Response.rspInfo:type_name -> trpc.keyword.Response.RspInfoEntry
	1, // 2: trpc.keyword.KeywordService.GetKeyword:input_type -> trpc.keyword.Request
	2, // 3: trpc.keyword.KeywordService.GetKeyword:output_type -> trpc.keyword.Response
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_keyword_proto_init() }
func file_keyword_proto_init() {
	if File_keyword_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_keyword_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Request); i {
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
		file_keyword_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Response); i {
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
		file_keyword_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Classroom); i {
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
	file_keyword_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_keyword_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_keyword_proto_goTypes,
		DependencyIndexes: file_keyword_proto_depIdxs,
		EnumInfos:         file_keyword_proto_enumTypes,
		MessageInfos:      file_keyword_proto_msgTypes,
	}.Build()
	File_keyword_proto = out.File
	file_keyword_proto_rawDesc = nil
	file_keyword_proto_goTypes = nil
	file_keyword_proto_depIdxs = nil
}