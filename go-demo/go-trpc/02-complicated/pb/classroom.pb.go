// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: classroom.proto

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomId int32 `protobuf:"varint,1,opt,name=roomId,proto3" json:"roomId,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classroom_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_classroom_proto_msgTypes[0]
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
	return file_classroom_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetRoomId() int32 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Classroom *Classroom `protobuf:"bytes,1,opt,name=classroom,proto3" json:"classroom,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classroom_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_classroom_proto_msgTypes[1]
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
	return file_classroom_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetClassroom() *Classroom {
	if x != nil {
		return x.Classroom
	}
	return nil
}

// 定义教室struct
type Classroom struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int32      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Address  string     `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	Students []*Student `protobuf:"bytes,4,rep,name=students,proto3" json:"students,omitempty"`
}

func (x *Classroom) Reset() {
	*x = Classroom{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classroom_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Classroom) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Classroom) ProtoMessage() {}

func (x *Classroom) ProtoReflect() protoreflect.Message {
	mi := &file_classroom_proto_msgTypes[2]
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
	return file_classroom_proto_rawDescGZIP(), []int{2}
}

func (x *Classroom) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Classroom) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Classroom) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Classroom) GetStudents() []*Student {
	if x != nil {
		return x.Students
	}
	return nil
}

var File_classroom_proto protoreflect.FileDescriptor

var file_classroom_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x72, 0x6f, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x64, 0x1a, 0x0d, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x21, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72,
	0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x22, 0x45, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x39, 0x0a, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x72, 0x6f, 0x6f, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x72, 0x6f, 0x6f,
	0x6d, 0x52, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x72, 0x6f, 0x6f, 0x6d, 0x22, 0x80, 0x01, 0x0a,
	0x09, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x72, 0x6f, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x35, 0x0a, 0x08, 0x73, 0x74, 0x75, 0x64,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x74, 0x72, 0x70,
	0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x53, 0x74,
	0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x32,
	0x56, 0x0a, 0x10, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x72, 0x6f, 0x6f, 0x6d, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19,
	0x2e, 0x74, 0x72, 0x70, 0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x64, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x74, 0x72, 0x70, 0x63,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2c, 0x5a, 0x2a, 0x7a, 0x69, 0x79, 0x69, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x74,
	0x72, 0x70, 0x63, 0x2f, 0x30, 0x32, 0x2d, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_classroom_proto_rawDescOnce sync.Once
	file_classroom_proto_rawDescData = file_classroom_proto_rawDesc
)

func file_classroom_proto_rawDescGZIP() []byte {
	file_classroom_proto_rawDescOnce.Do(func() {
		file_classroom_proto_rawDescData = protoimpl.X.CompressGZIP(file_classroom_proto_rawDescData)
	})
	return file_classroom_proto_rawDescData
}

var file_classroom_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_classroom_proto_goTypes = []any{
	(*Request)(nil),   // 0: trpc.complicated.Request
	(*Response)(nil),  // 1: trpc.complicated.Response
	(*Classroom)(nil), // 2: trpc.complicated.Classroom
	(*Student)(nil),   // 3: trpc.complicated.Student
}
var file_classroom_proto_depIdxs = []int32{
	2, // 0: trpc.complicated.Response.classroom:type_name -> trpc.complicated.Classroom
	3, // 1: trpc.complicated.Classroom.students:type_name -> trpc.complicated.Student
	0, // 2: trpc.complicated.ClassroomService.GetInfo:input_type -> trpc.complicated.Request
	1, // 3: trpc.complicated.ClassroomService.GetInfo:output_type -> trpc.complicated.Response
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_classroom_proto_init() }
func file_classroom_proto_init() {
	if File_classroom_proto != nil {
		return
	}
	file_student_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_classroom_proto_msgTypes[0].Exporter = func(v any, i int) any {
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
		file_classroom_proto_msgTypes[1].Exporter = func(v any, i int) any {
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
		file_classroom_proto_msgTypes[2].Exporter = func(v any, i int) any {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_classroom_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_classroom_proto_goTypes,
		DependencyIndexes: file_classroom_proto_depIdxs,
		MessageInfos:      file_classroom_proto_msgTypes,
	}.Build()
	File_classroom_proto = out.File
	file_classroom_proto_rawDesc = nil
	file_classroom_proto_goTypes = nil
	file_classroom_proto_depIdxs = nil
}
