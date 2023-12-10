// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: src/internal/core/server/transport/proto/response.proto

package transport

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

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  uint32            `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Id      string            `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Pid     string            `protobuf:"bytes,3,opt,name=pid,proto3" json:"pid,omitempty"`
	Body    []byte            `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	Headers map[string]string `protobuf:"bytes,6,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_internal_core_server_transport_proto_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_src_internal_core_server_transport_proto_response_proto_msgTypes[0]
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
	return file_src_internal_core_server_transport_proto_response_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Response) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Response) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *Response) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *Response) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

var File_src_internal_core_server_transport_proto_response_proto protoreflect.FileDescriptor

var file_src_internal_core_server_transport_proto_response_proto_rawDesc = []byte{
	0x0a, 0x37, 0x73, 0x72, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x22, 0xd0, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12,
	0x3a, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x24, 0x5a, 0x22, 0x73, 0x72, 0x63, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_src_internal_core_server_transport_proto_response_proto_rawDescOnce sync.Once
	file_src_internal_core_server_transport_proto_response_proto_rawDescData = file_src_internal_core_server_transport_proto_response_proto_rawDesc
)

func file_src_internal_core_server_transport_proto_response_proto_rawDescGZIP() []byte {
	file_src_internal_core_server_transport_proto_response_proto_rawDescOnce.Do(func() {
		file_src_internal_core_server_transport_proto_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_internal_core_server_transport_proto_response_proto_rawDescData)
	})
	return file_src_internal_core_server_transport_proto_response_proto_rawDescData
}

var file_src_internal_core_server_transport_proto_response_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_src_internal_core_server_transport_proto_response_proto_goTypes = []interface{}{
	(*Response)(nil), // 0: transport.Response
	nil,              // 1: transport.Response.HeadersEntry
}
var file_src_internal_core_server_transport_proto_response_proto_depIdxs = []int32{
	1, // 0: transport.Response.headers:type_name -> transport.Response.HeadersEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_src_internal_core_server_transport_proto_response_proto_init() }
func file_src_internal_core_server_transport_proto_response_proto_init() {
	if File_src_internal_core_server_transport_proto_response_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_src_internal_core_server_transport_proto_response_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_src_internal_core_server_transport_proto_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_src_internal_core_server_transport_proto_response_proto_goTypes,
		DependencyIndexes: file_src_internal_core_server_transport_proto_response_proto_depIdxs,
		MessageInfos:      file_src_internal_core_server_transport_proto_response_proto_msgTypes,
	}.Build()
	File_src_internal_core_server_transport_proto_response_proto = out.File
	file_src_internal_core_server_transport_proto_response_proto_rawDesc = nil
	file_src_internal_core_server_transport_proto_response_proto_goTypes = nil
	file_src_internal_core_server_transport_proto_response_proto_depIdxs = nil
}
