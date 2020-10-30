// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: pkg/protoc-gen-nidhi/test_data/simple.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
	_ "github.com/srikrsna/nidhi/nidhi"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Book struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title     string  `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Author    *Author `protobuf:"bytes,3,opt,name=author,proto3" json:"author,omitempty"`
	PageCount int32   `protobuf:"varint,4,opt,name=page_count,json=pageCount,proto3" json:"page_count,omitempty"`
}

func (x *Book) Reset() {
	*x = Book{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book) ProtoMessage() {}

func (x *Book) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book.ProtoReflect.Descriptor instead.
func (*Book) Descriptor() ([]byte, []int) {
	return file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescGZIP(), []int{0}
}

func (x *Book) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Book) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Book) GetAuthor() *Author {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *Book) GetPageCount() int32 {
	if x != nil {
		return x.PageCount
	}
	return 0
}

type Author struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Bio  string `protobuf:"bytes,2,opt,name=bio,proto3" json:"bio,omitempty"`
}

func (x *Author) Reset() {
	*x = Author{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Author) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Author) ProtoMessage() {}

func (x *Author) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Author.ProtoReflect.Descriptor instead.
func (*Author) Descriptor() ([]byte, []int) {
	return file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescGZIP(), []int{1}
}

func (x *Author) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Author) GetBio() string {
	if x != nil {
		return x.Bio
	}
	return ""
}

type Page struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number  int32  `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Page) Reset() {
	*x = Page{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Page) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Page) ProtoMessage() {}

func (x *Page) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Page.ProtoReflect.Descriptor instead.
func (*Page) Descriptor() ([]byte, []int) {
	return file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescGZIP(), []int{2}
}

func (x *Page) GetNumber() int32 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Page) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_pkg_protoc_gen_nidhi_test_data_simple_proto protoreflect.FileDescriptor

var file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74,
	0x65, 0x73, 0x74, 0x1a, 0x11, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x6e, 0x69, 0x64, 0x68, 0x69,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x04, 0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x70, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x3a, 0x07, 0xea, 0xfc, 0x09, 0x03, 0x62,
	0x6f, 0x6b, 0x22, 0x2e, 0x0a, 0x06, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x62, 0x69, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x62,
	0x69, 0x6f, 0x22, 0x38, 0x0a, 0x04, 0x50, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x42, 0x3d, 0x5a, 0x3b,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x72, 0x69, 0x6b, 0x72,
	0x73, 0x6e, 0x61, 0x2f, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x74,
	0x65, 0x73, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescOnce sync.Once
	file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescData = file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDesc
)

func file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescGZIP() []byte {
	file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescOnce.Do(func() {
		file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescData)
	})
	return file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDescData
}

var file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_protoc_gen_nidhi_test_data_simple_proto_goTypes = []interface{}{
	(*Book)(nil),   // 0: test.Book
	(*Author)(nil), // 1: test.Author
	(*Page)(nil),   // 2: test.Page
}
var file_pkg_protoc_gen_nidhi_test_data_simple_proto_depIdxs = []int32{
	1, // 0: test.Book.author:type_name -> test.Author
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_protoc_gen_nidhi_test_data_simple_proto_init() }
func file_pkg_protoc_gen_nidhi_test_data_simple_proto_init() {
	if File_pkg_protoc_gen_nidhi_test_data_simple_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Book); i {
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
		file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Author); i {
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
		file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Page); i {
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
			RawDescriptor: file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_protoc_gen_nidhi_test_data_simple_proto_goTypes,
		DependencyIndexes: file_pkg_protoc_gen_nidhi_test_data_simple_proto_depIdxs,
		MessageInfos:      file_pkg_protoc_gen_nidhi_test_data_simple_proto_msgTypes,
	}.Build()
	File_pkg_protoc_gen_nidhi_test_data_simple_proto = out.File
	file_pkg_protoc_gen_nidhi_test_data_simple_proto_rawDesc = nil
	file_pkg_protoc_gen_nidhi_test_data_simple_proto_goTypes = nil
	file_pkg_protoc_gen_nidhi_test_data_simple_proto_depIdxs = nil
}
