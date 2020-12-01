// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.13.0
// source: internal/protoc-gen-nidhi/test_data/all.proto

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

type All struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	StringField       string   `protobuf:"bytes,2,opt,name=string_field,json=stringField,proto3" json:"string_field,omitempty"`
	Int32Field        int32    `protobuf:"varint,3,opt,name=int32_field,json=int32Field,proto3" json:"int32_field,omitempty"`
	Int64Field        int64    `protobuf:"varint,4,opt,name=int64_field,json=int64Field,proto3" json:"int64_field,omitempty"`
	Uint32Field       uint32   `protobuf:"varint,5,opt,name=uint32_field,json=uint32Field,proto3" json:"uint32_field,omitempty"`
	Uint64Field       uint64   `protobuf:"varint,6,opt,name=uint64_field,json=uint64Field,proto3" json:"uint64_field,omitempty"`
	FloatField        float32  `protobuf:"fixed32,7,opt,name=float_field,json=floatField,proto3" json:"float_field,omitempty"`
	DoubleField       float64  `protobuf:"fixed64,8,opt,name=double_field,json=doubleField,proto3" json:"double_field,omitempty"`
	BoolField         bool     `protobuf:"varint,9,opt,name=bool_field,json=boolField,proto3" json:"bool_field,omitempty"`
	BytesField        []byte   `protobuf:"bytes,10,opt,name=bytes_field,json=bytesField,proto3" json:"bytes_field,omitempty"`
	PrimitiveRepeated []string `protobuf:"bytes,11,rep,name=primitive_repeated,json=primitiveRepeated,proto3" json:"primitive_repeated,omitempty"`
	// Types that are assignable to OneOf:
	//	*All_StringOneOf
	//	*All_Int32OneOf
	//	*All_Int64OneOf
	//	*All_Uint32OneOf
	//	*All_Uint64OneOf
	//	*All_FloatOneOf
	//	*All_DoubleOneOf
	//	*All_BoolOneOf
	//	*All_BytesOneOf
	//	*All_SimpleObjectOneOf
	OneOf             isAll_OneOf `protobuf_oneof:"one_of"`
	SimpleObjectField *Simple     `protobuf:"bytes,22,opt,name=simple_object_field,json=simpleObjectField,proto3" json:"simple_object_field,omitempty"`
	SimpleRepeated    []*Simple   `protobuf:"bytes,23,rep,name=simple_repeated,json=simpleRepeated,proto3" json:"simple_repeated,omitempty"`
	NestedOne         *NestedOne  `protobuf:"bytes,24,opt,name=nested_one,json=nestedOne,proto3" json:"nested_one,omitempty"`
}

func (x *All) Reset() {
	*x = All{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *All) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*All) ProtoMessage() {}

func (x *All) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use All.ProtoReflect.Descriptor instead.
func (*All) Descriptor() ([]byte, []int) {
	return file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescGZIP(), []int{0}
}

func (x *All) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *All) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

func (x *All) GetInt32Field() int32 {
	if x != nil {
		return x.Int32Field
	}
	return 0
}

func (x *All) GetInt64Field() int64 {
	if x != nil {
		return x.Int64Field
	}
	return 0
}

func (x *All) GetUint32Field() uint32 {
	if x != nil {
		return x.Uint32Field
	}
	return 0
}

func (x *All) GetUint64Field() uint64 {
	if x != nil {
		return x.Uint64Field
	}
	return 0
}

func (x *All) GetFloatField() float32 {
	if x != nil {
		return x.FloatField
	}
	return 0
}

func (x *All) GetDoubleField() float64 {
	if x != nil {
		return x.DoubleField
	}
	return 0
}

func (x *All) GetBoolField() bool {
	if x != nil {
		return x.BoolField
	}
	return false
}

func (x *All) GetBytesField() []byte {
	if x != nil {
		return x.BytesField
	}
	return nil
}

func (x *All) GetPrimitiveRepeated() []string {
	if x != nil {
		return x.PrimitiveRepeated
	}
	return nil
}

func (m *All) GetOneOf() isAll_OneOf {
	if m != nil {
		return m.OneOf
	}
	return nil
}

func (x *All) GetStringOneOf() string {
	if x, ok := x.GetOneOf().(*All_StringOneOf); ok {
		return x.StringOneOf
	}
	return ""
}

func (x *All) GetInt32OneOf() int32 {
	if x, ok := x.GetOneOf().(*All_Int32OneOf); ok {
		return x.Int32OneOf
	}
	return 0
}

func (x *All) GetInt64OneOf() int64 {
	if x, ok := x.GetOneOf().(*All_Int64OneOf); ok {
		return x.Int64OneOf
	}
	return 0
}

func (x *All) GetUint32OneOf() uint32 {
	if x, ok := x.GetOneOf().(*All_Uint32OneOf); ok {
		return x.Uint32OneOf
	}
	return 0
}

func (x *All) GetUint64OneOf() uint64 {
	if x, ok := x.GetOneOf().(*All_Uint64OneOf); ok {
		return x.Uint64OneOf
	}
	return 0
}

func (x *All) GetFloatOneOf() float32 {
	if x, ok := x.GetOneOf().(*All_FloatOneOf); ok {
		return x.FloatOneOf
	}
	return 0
}

func (x *All) GetDoubleOneOf() float64 {
	if x, ok := x.GetOneOf().(*All_DoubleOneOf); ok {
		return x.DoubleOneOf
	}
	return 0
}

func (x *All) GetBoolOneOf() bool {
	if x, ok := x.GetOneOf().(*All_BoolOneOf); ok {
		return x.BoolOneOf
	}
	return false
}

func (x *All) GetBytesOneOf() []byte {
	if x, ok := x.GetOneOf().(*All_BytesOneOf); ok {
		return x.BytesOneOf
	}
	return nil
}

func (x *All) GetSimpleObjectOneOf() *Simple {
	if x, ok := x.GetOneOf().(*All_SimpleObjectOneOf); ok {
		return x.SimpleObjectOneOf
	}
	return nil
}

func (x *All) GetSimpleObjectField() *Simple {
	if x != nil {
		return x.SimpleObjectField
	}
	return nil
}

func (x *All) GetSimpleRepeated() []*Simple {
	if x != nil {
		return x.SimpleRepeated
	}
	return nil
}

func (x *All) GetNestedOne() *NestedOne {
	if x != nil {
		return x.NestedOne
	}
	return nil
}

type isAll_OneOf interface {
	isAll_OneOf()
}

type All_StringOneOf struct {
	StringOneOf string `protobuf:"bytes,12,opt,name=string_one_of,json=stringOneOf,proto3,oneof"`
}

type All_Int32OneOf struct {
	Int32OneOf int32 `protobuf:"varint,13,opt,name=int32_one_of,json=int32OneOf,proto3,oneof"`
}

type All_Int64OneOf struct {
	Int64OneOf int64 `protobuf:"varint,14,opt,name=int64_one_of,json=int64OneOf,proto3,oneof"`
}

type All_Uint32OneOf struct {
	Uint32OneOf uint32 `protobuf:"varint,15,opt,name=uint32_one_of,json=uint32OneOf,proto3,oneof"`
}

type All_Uint64OneOf struct {
	Uint64OneOf uint64 `protobuf:"varint,16,opt,name=uint64_one_of,json=uint64OneOf,proto3,oneof"`
}

type All_FloatOneOf struct {
	FloatOneOf float32 `protobuf:"fixed32,17,opt,name=float_one_of,json=floatOneOf,proto3,oneof"`
}

type All_DoubleOneOf struct {
	DoubleOneOf float64 `protobuf:"fixed64,18,opt,name=double_one_of,json=doubleOneOf,proto3,oneof"`
}

type All_BoolOneOf struct {
	BoolOneOf bool `protobuf:"varint,19,opt,name=bool_one_of,json=boolOneOf,proto3,oneof"`
}

type All_BytesOneOf struct {
	BytesOneOf []byte `protobuf:"bytes,20,opt,name=bytes_one_of,json=bytesOneOf,proto3,oneof"`
}

type All_SimpleObjectOneOf struct {
	SimpleObjectOneOf *Simple `protobuf:"bytes,21,opt,name=simple_object_one_of,json=simpleObjectOneOf,proto3,oneof"`
}

func (*All_StringOneOf) isAll_OneOf() {}

func (*All_Int32OneOf) isAll_OneOf() {}

func (*All_Int64OneOf) isAll_OneOf() {}

func (*All_Uint32OneOf) isAll_OneOf() {}

func (*All_Uint64OneOf) isAll_OneOf() {}

func (*All_FloatOneOf) isAll_OneOf() {}

func (*All_DoubleOneOf) isAll_OneOf() {}

func (*All_BoolOneOf) isAll_OneOf() {}

func (*All_BytesOneOf) isAll_OneOf() {}

func (*All_SimpleObjectOneOf) isAll_OneOf() {}

type Simple struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StringField string `protobuf:"bytes,1,opt,name=string_field,json=stringField,proto3" json:"string_field,omitempty"`
}

func (x *Simple) Reset() {
	*x = Simple{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Simple) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Simple) ProtoMessage() {}

func (x *Simple) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Simple.ProtoReflect.Descriptor instead.
func (*Simple) Descriptor() ([]byte, []int) {
	return file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescGZIP(), []int{1}
}

func (x *Simple) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

type NestedOne struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NestetedInt int32      `protobuf:"varint,1,opt,name=nesteted_int,json=nestetedInt,proto3" json:"nesteted_int,omitempty"`
	Nested      *NestedTwo `protobuf:"bytes,2,opt,name=nested,proto3" json:"nested,omitempty"`
}

func (x *NestedOne) Reset() {
	*x = NestedOne{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NestedOne) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NestedOne) ProtoMessage() {}

func (x *NestedOne) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NestedOne.ProtoReflect.Descriptor instead.
func (*NestedOne) Descriptor() ([]byte, []int) {
	return file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescGZIP(), []int{2}
}

func (x *NestedOne) GetNestetedInt() int32 {
	if x != nil {
		return x.NestetedInt
	}
	return 0
}

func (x *NestedOne) GetNested() *NestedTwo {
	if x != nil {
		return x.Nested
	}
	return nil
}

type NestedTwo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SomeField string         `protobuf:"bytes,1,opt,name=some_field,json=someField,proto3" json:"some_field,omitempty"`
	Nested    []*NestedThree `protobuf:"bytes,2,rep,name=nested,proto3" json:"nested,omitempty"`
}

func (x *NestedTwo) Reset() {
	*x = NestedTwo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NestedTwo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NestedTwo) ProtoMessage() {}

func (x *NestedTwo) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NestedTwo.ProtoReflect.Descriptor instead.
func (*NestedTwo) Descriptor() ([]byte, []int) {
	return file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescGZIP(), []int{3}
}

func (x *NestedTwo) GetSomeField() string {
	if x != nil {
		return x.SomeField
	}
	return ""
}

func (x *NestedTwo) GetNested() []*NestedThree {
	if x != nil {
		return x.Nested
	}
	return nil
}

type NestedThree struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Some string `protobuf:"bytes,1,opt,name=some,proto3" json:"some,omitempty"`
}

func (x *NestedThree) Reset() {
	*x = NestedThree{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NestedThree) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NestedThree) ProtoMessage() {}

func (x *NestedThree) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NestedThree.ProtoReflect.Descriptor instead.
func (*NestedThree) Descriptor() ([]byte, []int) {
	return file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescGZIP(), []int{4}
}

func (x *NestedThree) GetSome() string {
	if x != nil {
		return x.Some
	}
	return ""
}

var File_internal_protoc_gen_nidhi_test_data_all_proto protoreflect.FileDescriptor

var file_internal_protoc_gen_nidhi_test_data_all_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x61, 0x6c, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x74, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x6e, 0x69, 0x64,
	0x68, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb6, 0x07, 0x0a, 0x03, 0x41, 0x6c, 0x6c,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x36, 0x34,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x75, 0x69, 0x6e,
	0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x66,
	0x6c, 0x6f, 0x61, 0x74, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x0a, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x21, 0x0a, 0x0c,
	0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x62, 0x6f, 0x6f, 0x6c, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x73, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12,
	0x2d, 0x0a, 0x12, 0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x72, 0x65, 0x70,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x09, 0x52, 0x11, 0x70, 0x72, 0x69,
	0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x24,
	0x0a, 0x0d, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6f, 0x6e, 0x65, 0x5f, 0x6f, 0x66, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4f,
	0x6e, 0x65, 0x4f, 0x66, 0x12, 0x22, 0x0a, 0x0c, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x6f, 0x6e,
	0x65, 0x5f, 0x6f, 0x66, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e,
	0x74, 0x33, 0x32, 0x4f, 0x6e, 0x65, 0x4f, 0x66, 0x12, 0x22, 0x0a, 0x0c, 0x69, 0x6e, 0x74, 0x36,
	0x34, 0x5f, 0x6f, 0x6e, 0x65, 0x5f, 0x6f, 0x66, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00,
	0x52, 0x0a, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x4f, 0x6e, 0x65, 0x4f, 0x66, 0x12, 0x24, 0x0a, 0x0d,
	0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x6f, 0x6e, 0x65, 0x5f, 0x6f, 0x66, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x0b, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x4f, 0x6e, 0x65,
	0x4f, 0x66, 0x12, 0x24, 0x0a, 0x0d, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x6f, 0x6e, 0x65,
	0x5f, 0x6f, 0x66, 0x18, 0x10, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x0b, 0x75, 0x69, 0x6e,
	0x74, 0x36, 0x34, 0x4f, 0x6e, 0x65, 0x4f, 0x66, 0x12, 0x22, 0x0a, 0x0c, 0x66, 0x6c, 0x6f, 0x61,
	0x74, 0x5f, 0x6f, 0x6e, 0x65, 0x5f, 0x6f, 0x66, 0x18, 0x11, 0x20, 0x01, 0x28, 0x02, 0x48, 0x00,
	0x52, 0x0a, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x4f, 0x6e, 0x65, 0x4f, 0x66, 0x12, 0x24, 0x0a, 0x0d,
	0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5f, 0x6f, 0x6e, 0x65, 0x5f, 0x6f, 0x66, 0x18, 0x12, 0x20,
	0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x4f, 0x6e, 0x65,
	0x4f, 0x66, 0x12, 0x20, 0x0a, 0x0b, 0x62, 0x6f, 0x6f, 0x6c, 0x5f, 0x6f, 0x6e, 0x65, 0x5f, 0x6f,
	0x66, 0x18, 0x13, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x4f,
	0x6e, 0x65, 0x4f, 0x66, 0x12, 0x22, 0x0a, 0x0c, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x6f, 0x6e,
	0x65, 0x5f, 0x6f, 0x66, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0a, 0x62, 0x79,
	0x74, 0x65, 0x73, 0x4f, 0x6e, 0x65, 0x4f, 0x66, 0x12, 0x3f, 0x0a, 0x14, 0x73, 0x69, 0x6d, 0x70,
	0x6c, 0x65, 0x5f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6f, 0x6e, 0x65, 0x5f, 0x6f, 0x66,
	0x18, 0x15, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x69,
	0x6d, 0x70, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x11, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x4f, 0x6e, 0x65, 0x4f, 0x66, 0x12, 0x3c, 0x0a, 0x13, 0x73, 0x69, 0x6d,
	0x70, 0x6c, 0x65, 0x5f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x69,
	0x6d, 0x70, 0x6c, 0x65, 0x52, 0x11, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x35, 0x0a, 0x0f, 0x73, 0x69, 0x6d, 0x70, 0x6c,
	0x65, 0x5f, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x17, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x0e,
	0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x2e,
	0x0a, 0x0a, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x65, 0x18, 0x18, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x4f, 0x6e, 0x65, 0x52, 0x09, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4f, 0x6e, 0x65, 0x3a, 0x07,
	0xea, 0xfc, 0x09, 0x03, 0x61, 0x6c, 0x6c, 0x42, 0x08, 0x0a, 0x06, 0x6f, 0x6e, 0x65, 0x5f, 0x6f,
	0x66, 0x22, 0x2b, 0x0a, 0x06, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x22, 0x57,
	0x0a, 0x09, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4f, 0x6e, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6e,
	0x65, 0x73, 0x74, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x74, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x12, 0x27,
	0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x54, 0x77, 0x6f, 0x52,
	0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x22, 0x55, 0x0a, 0x09, 0x4e, 0x65, 0x73, 0x74, 0x65,
	0x64, 0x54, 0x77, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x6f, 0x6d, 0x65, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x6f, 0x6d, 0x65, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x12, 0x29, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65,
	0x64, 0x54, 0x68, 0x72, 0x65, 0x65, 0x52, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x22, 0x21,
	0x0a, 0x0b, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x54, 0x68, 0x72, 0x65, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x6f, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x6d,
	0x65, 0x42, 0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x73, 0x72, 0x69, 0x6b, 0x72, 0x73, 0x6e, 0x61, 0x2f, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67,
	0x65, 0x6e, 0x2d, 0x6e, 0x69, 0x64, 0x68, 0x69, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescOnce sync.Once
	file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescData = file_internal_protoc_gen_nidhi_test_data_all_proto_rawDesc
)

func file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescGZIP() []byte {
	file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescOnce.Do(func() {
		file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescData)
	})
	return file_internal_protoc_gen_nidhi_test_data_all_proto_rawDescData
}

var file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_internal_protoc_gen_nidhi_test_data_all_proto_goTypes = []interface{}{
	(*All)(nil),         // 0: test.All
	(*Simple)(nil),      // 1: test.Simple
	(*NestedOne)(nil),   // 2: test.NestedOne
	(*NestedTwo)(nil),   // 3: test.NestedTwo
	(*NestedThree)(nil), // 4: test.NestedThree
}
var file_internal_protoc_gen_nidhi_test_data_all_proto_depIdxs = []int32{
	1, // 0: test.All.simple_object_one_of:type_name -> test.Simple
	1, // 1: test.All.simple_object_field:type_name -> test.Simple
	1, // 2: test.All.simple_repeated:type_name -> test.Simple
	2, // 3: test.All.nested_one:type_name -> test.NestedOne
	3, // 4: test.NestedOne.nested:type_name -> test.NestedTwo
	4, // 5: test.NestedTwo.nested:type_name -> test.NestedThree
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_internal_protoc_gen_nidhi_test_data_all_proto_init() }
func file_internal_protoc_gen_nidhi_test_data_all_proto_init() {
	if File_internal_protoc_gen_nidhi_test_data_all_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*All); i {
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
		file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Simple); i {
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
		file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NestedOne); i {
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
		file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NestedTwo); i {
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
		file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NestedThree); i {
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
	file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*All_StringOneOf)(nil),
		(*All_Int32OneOf)(nil),
		(*All_Int64OneOf)(nil),
		(*All_Uint32OneOf)(nil),
		(*All_Uint64OneOf)(nil),
		(*All_FloatOneOf)(nil),
		(*All_DoubleOneOf)(nil),
		(*All_BoolOneOf)(nil),
		(*All_BytesOneOf)(nil),
		(*All_SimpleObjectOneOf)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_protoc_gen_nidhi_test_data_all_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_protoc_gen_nidhi_test_data_all_proto_goTypes,
		DependencyIndexes: file_internal_protoc_gen_nidhi_test_data_all_proto_depIdxs,
		MessageInfos:      file_internal_protoc_gen_nidhi_test_data_all_proto_msgTypes,
	}.Build()
	File_internal_protoc_gen_nidhi_test_data_all_proto = out.File
	file_internal_protoc_gen_nidhi_test_data_all_proto_rawDesc = nil
	file_internal_protoc_gen_nidhi_test_data_all_proto_goTypes = nil
	file_internal_protoc_gen_nidhi_test_data_all_proto_depIdxs = nil
}
