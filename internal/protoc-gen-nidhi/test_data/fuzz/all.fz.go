// Code generated by protoc-gen-nidhi. DO NOT EDIT.
// source: internal/protoc-gen-nidhi/test_data/all.proto

package fuzz

import (
	fuzz "github.com/google/gofuzz"
	pb "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi/test_data"
)

func FuzzAll(msg *pb.All, c fuzz.Continue) {
    c.Fuzz(&msg.Id)
    c.Fuzz(&msg.StringField)
    c.Fuzz(&msg.Int32Field)
    c.Fuzz(&msg.Int64Field)
    c.Fuzz(&msg.Uint32Field)
    c.Fuzz(&msg.Uint64Field)
    c.Fuzz(&msg.FloatField)
    c.Fuzz(&msg.DoubleField)
    c.Fuzz(&msg.BoolField)
    c.Fuzz(&msg.BytesField)
    c.Fuzz(&msg.PrimitiveRepeated)
    c.Fuzz(&msg.SimpleObjectField)
    c.Fuzz(&msg.SimpleRepeated)
    c.Fuzz(&msg.NestedOne)
    c.Fuzz(&msg.Timestamp)
    c.Fuzz(&msg.AnyField)
    switch c.Int31n(10) {
    case 0:
        var f pb.All_StringOneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 1:
        var f pb.All_Int32OneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 2:
        var f pb.All_Int64OneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 3:
        var f pb.All_Uint32OneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 4:
        var f pb.All_Uint64OneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 5:
        var f pb.All_FloatOneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 6:
        var f pb.All_DoubleOneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 7:
        var f pb.All_BoolOneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 8:
        var f pb.All_BytesOneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    case 9:
        var f pb.All_SimpleObjectOneOf
        c.Fuzz(&f)
        msg.OneOf = &f
    }
}
 
func FuzzAll_StringOneOf(msg *pb.All_StringOneOf, c fuzz.Continue) {
    c.Fuzz(&msg.StringOneOf)
}
 
func FuzzAll_Int32OneOf(msg *pb.All_Int32OneOf, c fuzz.Continue) {
    c.Fuzz(&msg.Int32OneOf)
}
 
func FuzzAll_Int64OneOf(msg *pb.All_Int64OneOf, c fuzz.Continue) {
    c.Fuzz(&msg.Int64OneOf)
}
 
func FuzzAll_Uint32OneOf(msg *pb.All_Uint32OneOf, c fuzz.Continue) {
    c.Fuzz(&msg.Uint32OneOf)
}
 
func FuzzAll_Uint64OneOf(msg *pb.All_Uint64OneOf, c fuzz.Continue) {
    c.Fuzz(&msg.Uint64OneOf)
}
 
func FuzzAll_FloatOneOf(msg *pb.All_FloatOneOf, c fuzz.Continue) {
    c.Fuzz(&msg.FloatOneOf)
}
 
func FuzzAll_DoubleOneOf(msg *pb.All_DoubleOneOf, c fuzz.Continue) {
    c.Fuzz(&msg.DoubleOneOf)
}
 
func FuzzAll_BoolOneOf(msg *pb.All_BoolOneOf, c fuzz.Continue) {
    c.Fuzz(&msg.BoolOneOf)
}
 
func FuzzAll_BytesOneOf(msg *pb.All_BytesOneOf, c fuzz.Continue) {
    c.Fuzz(&msg.BytesOneOf)
}
 
func FuzzAll_SimpleObjectOneOf(msg *pb.All_SimpleObjectOneOf, c fuzz.Continue) {
    c.Fuzz(&msg.SimpleObjectOneOf)
}

func FuzzSimple(msg *pb.Simple, c fuzz.Continue) {
    c.Fuzz(&msg.StringField)
}

func FuzzNestedOne(msg *pb.NestedOne, c fuzz.Continue) {
    c.Fuzz(&msg.NestetedInt)
    c.Fuzz(&msg.Nested)
    c.Fuzz(&msg.T)
    c.Fuzz(&msg.A)
}

func FuzzNestedTwo(msg *pb.NestedTwo, c fuzz.Continue) {
    c.Fuzz(&msg.SomeField)
    c.Fuzz(&msg.Nested)
}

func FuzzNestedThree(msg *pb.NestedThree, c fuzz.Continue) {
    c.Fuzz(&msg.Some)
}
