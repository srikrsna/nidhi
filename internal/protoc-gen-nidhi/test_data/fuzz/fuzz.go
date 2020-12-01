package fuzz

var ff = [...]interface{}{
	FuzzAll,
	FuzzAll_StringOneOf,
	FuzzAll_Int32OneOf,
	FuzzAll_Int64OneOf,
	FuzzAll_Uint32OneOf,
	FuzzAll_Uint64OneOf,
	FuzzAll_FloatOneOf,
	FuzzAll_DoubleOneOf,
	FuzzAll_BoolOneOf,
	FuzzAll_BytesOneOf,
	FuzzAll_SimpleObjectOneOf,
	FuzzSimple,
	FuzzNestedOne,
	FuzzNestedTwo,
	FuzzNestedThree,
}

func FuzzFuncs() []interface{} {
	return ff[:]
}
