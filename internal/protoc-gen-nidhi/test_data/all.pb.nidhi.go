// Code generated by protoc-gen-nidhi. DO NOT EDIT.
// source: internal/protoc-gen-nidhi/test_data/all.proto

package pb

import (
	"context"
	"database/sql"
	"errors"

	jsoniter "github.com/json-iterator/go"

	"github.com/srikrsna/nidhi"
	nidhigen "github.com/srikrsna/nidhi/nidhigen"
)

var (
	_ = context.Background
	_ = (*sql.DB)(nil)
	_ = errors.New
	_ = jsoniter.Marshal
	_ = nidhi.JSONB
	_ = nidhigen.WriteString
)

func (doc *All) DocumentId() string {
	return doc.Id
}

func (doc *All) SetDocumentId(id string) {
	doc.Id = id
}

type AllCollection struct {
	*allCollection

	ogCol *nidhi.Collection
}

func OpenAllCollection(ctx context.Context, db *sql.DB) (*AllCollection, error) {
	col, err := nidhi.OpenCollection(ctx, db, "pb", "alls", nidhi.CollectionOptions{
		Fields: []string{"id", "stringField", "int32Field", "int64Field", "uint32Field", "uint64Field", "floatField", "doubleField", "boolField", "bytesField", "primitiveRepeated", "stringOneOf", "int32OneOf", "int64OneOf", "uint32OneOf", "uint64OneOf", "floatOneOf", "doubleOneOf", "boolOneOf", "bytesOneOf", "simpleObjectOneOf", "simpleObjectField", "simpleRepeated", "nestedOne"},
	})
	if err != nil {
		return nil, err
	}
	return &AllCollection{
		&allCollection{col: col},
		col,
	}, nil
}

func (st *AllCollection) BeginTx(ctx context.Context, opt *sql.TxOptions) (*AllTxCollection, error) {
	txCol, err := st.ogCol.BeginTx(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &AllTxCollection{&allCollection{txCol}, txCol}, nil
}

func (st *AllCollection) WithTransaction(tx *nidhi.TxToken) *AllTxCollection {
	txCol := st.ogCol.WithTransaction(tx)
	return &AllTxCollection{&allCollection{txCol}, txCol}
}

type AllTxCollection struct {
	*allCollection
	txCol *nidhi.TxCollection
}

func (tx *AllTxCollection) Rollback() error {
	return tx.txCol.Rollback()
}

func (tx *AllTxCollection) Commit() error {
	return tx.txCol.Commit()
}

func (tx *AllTxCollection) TxToken() *nidhi.TxToken {
	return nidhi.NewTxToken(tx.txCol)
}

type allCollection struct {
	col nidhigen.Collection
}

func (st *allCollection) CreateAll(ctx context.Context, b *All, ops ...nidhi.CreateOption) (string, error) {
	return st.col.Create(ctx, b, ops)
}

func (st *allCollection) QueryAlls(ctx context.Context, f isAllQuery, ops ...nidhi.QueryOption) ([]*All, error) {
	var ee []*All
	ctr := func() nidhi.Document {
		var e All
		ee = append(ee, &e)
		return &e
	}

	return ee, st.col.Query(ctx, f, ctr, ops)
}

func (st *allCollection) ReplaceAll(ctx context.Context, b *All, ops ...nidhi.ReplaceOption) error {
	return st.col.Replace(ctx, b, ops)
}

func (st *allCollection) DeleteAll(ctx context.Context, id string, ops ...nidhi.DeleteOption) error {
	return st.col.Delete(ctx, id, ops)
}

func (st *allCollection) GetAll(ctx context.Context, id string, ops ...nidhi.GetOption) (*All, error) {
	var entity All
	return &entity, st.col.Get(ctx, id, &entity, ops)
}

func (st *allCollection) UpdateAlls(ctx context.Context, b *All, f isAllQuery, ops ...nidhi.UpdateOption) error {
	return st.col.Update(ctx, b, f, ops)
}

func (st *allCollection) DeleteAlls(ctx context.Context, f isAllQuery, ops ...nidhi.DeleteOption) error {
	return st.col.DeleteMany(ctx, f, ops)
}

func GetAllQuery() AllQuery {
	return (*imp_AllQuery)(nidhi.GetQuery())
}

func PutAllQuery(q AllQuery) {
	nidhi.PutQuery((*nidhi.Query)(q.(*imp_AllQuery)))
}

type AllQuery interface {
	Id(*nidhi.StringQuery) AllConj
	StringField(*nidhi.StringQuery) AllConj
	Int32Field(*nidhi.IntQuery) AllConj
	Int64Field(*nidhi.IntQuery) AllConj
	Uint32Field(*nidhi.IntQuery) AllConj
	Uint64Field(*nidhi.IntQuery) AllConj
	FloatField(*nidhi.FloatQuery) AllConj
	DoubleField(*nidhi.FloatQuery) AllConj
	BoolField(*nidhi.BoolQuery) AllConj
	PrimitiveRepeated(nidhi.SliceOptions, ...string) AllConj
	StringOneOf(*nidhi.StringQuery) AllConj
	Int32OneOf(*nidhi.IntQuery) AllConj
	Int64OneOf(*nidhi.IntQuery) AllConj
	Uint32OneOf(*nidhi.IntQuery) AllConj
	Uint64OneOf(*nidhi.IntQuery) AllConj
	FloatOneOf(*nidhi.FloatQuery) AllConj
	DoubleOneOf(*nidhi.FloatQuery) AllConj
	BoolOneOf(*nidhi.BoolQuery) AllConj
	SimpleObjectOneOf() AllSimpleObjectOneOfQuery
	SimpleObjectField() AllSimpleObjectFieldQuery
	SimpleRepeated(...*Simple) AllConj
	NestedOne() AllNestedOneQuery

	// Generic With Type Safety
	Paren(iq isAllQuery) AllConj
	Where(query string, args ...interface{}) AllConj
	Not() AllQuery
	ReplaceArgs(args ...interface{}) error
}

type AllConj interface {
	And() AllQuery
	Or() AllQuery
	isAllQuery
}

type isAllQuery interface {
	imp_AllQuery()
	nidhi.Sqlizer
}

type imp_AllQuery nidhi.Query

func (q *imp_AllQuery) imp_AllQuery() {}

func (q *imp_AllQuery) Id(f *nidhi.StringQuery) AllConj {
	(*nidhi.Query)(q).Id(f)
	return q
}
func (q *imp_AllQuery) StringField(f *nidhi.StringQuery) AllConj {
	(*nidhi.Query)(q).Field(" "+nidhi.ColDoc+"->>'stringField'", f)
	return q
}

func (q *imp_AllQuery) Int32Field(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'int32Field')::bigint", f)
	return q
}

func (q *imp_AllQuery) Int64Field(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'int64Field')::bigint", f)
	return q
}

func (q *imp_AllQuery) Uint32Field(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'uint32Field')::bigint", f)
	return q
}

func (q *imp_AllQuery) Uint64Field(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'uint64Field')::bigint", f)
	return q
}

func (q *imp_AllQuery) FloatField(f *nidhi.FloatQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'floatField')::double precision", f)
	return q
}

func (q *imp_AllQuery) DoubleField(f *nidhi.FloatQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'doubleField')::double precision", f)
	return q
}

func (q *imp_AllQuery) BoolField(f *nidhi.BoolQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'boolField')::bool", f)
	return q
}

func (q *imp_AllQuery) PrimitiveRepeated(opt nidhi.SliceOptions, arr ...string) AllConj {
	(*nidhi.Query)(q).Field(
		" "+nidhi.ColDoc+"->'primitiveRepeated'",
		&nidhi.SliceQuery{
			Slice:   arr,
			Options: opt,
		},
	)
	return q
}

func (q *imp_AllQuery) StringOneOf(f *nidhi.StringQuery) AllConj {
	(*nidhi.Query)(q).Field(" "+nidhi.ColDoc+"->>'stringOneOf'", f)
	return q
}

func (q *imp_AllQuery) Int32OneOf(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'int32OneOf')::bigint", f)
	return q
}

func (q *imp_AllQuery) Int64OneOf(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'int64OneOf')::bigint", f)
	return q
}

func (q *imp_AllQuery) Uint32OneOf(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'uint32OneOf')::bigint", f)
	return q
}

func (q *imp_AllQuery) Uint64OneOf(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'uint64OneOf')::bigint", f)
	return q
}

func (q *imp_AllQuery) FloatOneOf(f *nidhi.FloatQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'floatOneOf')::double precision", f)
	return q
}

func (q *imp_AllQuery) DoubleOneOf(f *nidhi.FloatQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'doubleOneOf')::double precision", f)
	return q
}

func (q *imp_AllQuery) BoolOneOf(f *nidhi.BoolQuery) AllConj {
	(*nidhi.Query)(q).Field(" ("+nidhi.ColDoc+"->'boolOneOf')::bool", f)
	return q
}

func (q *imp_AllQuery) SimpleObjectOneOf() AllSimpleObjectOneOfQuery {
	(*nidhi.Query)(q).Prefix(" (" + nidhi.ColDoc + "->'simpleObjectOneOf'")
	return (*imp_AllSimpleObjectOneOfQuery)(q)
}
func (q *imp_AllQuery) SimpleObjectField() AllSimpleObjectFieldQuery {
	(*nidhi.Query)(q).Prefix(" (" + nidhi.ColDoc + "->'simpleObjectField'")
	return (*imp_AllSimpleObjectFieldQuery)(q)
}
func (q *imp_AllQuery) SimpleRepeated(arr ...*Simple) AllConj {
	(*nidhi.Query)(q).Field(
		" "+nidhi.ColDoc+"->'simpleRepeated'",
		nidhi.MarshalerQuery{
			Marshaler: SimpleSlice(arr),
		},
	)
	return q
}

func (q *imp_AllQuery) NestedOne() AllNestedOneQuery {
	(*nidhi.Query)(q).Prefix(" (" + nidhi.ColDoc + "->'nestedOne'")
	return (*imp_AllNestedOneQuery)(q)
}

func (q *imp_AllQuery) Paren(iq isAllQuery) AllConj {
	(*nidhi.Query)(q).Paren(iq)
	return q
}

func (q *imp_AllQuery) Where(query string, args ...interface{}) AllConj {
	(*nidhi.Query)(q).Where(query, args...)
	return q
}

func (q *imp_AllQuery) Not() AllQuery {
	(*nidhi.Query)(q).Not()
	return q
}

func (q *imp_AllQuery) And() AllQuery {
	(*nidhi.Query)(q).And()
	return q
}

func (q *imp_AllQuery) Or() AllQuery {
	(*nidhi.Query)(q).Or()
	return q
}

func (q *imp_AllQuery) ReplaceArgs(args ...interface{}) error {
	return (*nidhi.Query)(q).ReplaceArgs()
}

func (q *imp_AllQuery) ToSql() (string, []interface{}, error) {
	return (*nidhi.Query)(q).ToSql()
}

type AllSimpleObjectOneOfQuery interface {
	StringField(*nidhi.StringQuery) AllConj
}

type imp_AllSimpleObjectOneOfQuery nidhi.Query

func (q *imp_AllSimpleObjectOneOfQuery) StringField(f *nidhi.StringQuery) AllConj {
	(*nidhi.Query)(q).Field("->>'stringField')", f)
	return (*imp_AllQuery)(q)
}

type AllSimpleObjectFieldQuery interface {
	StringField(*nidhi.StringQuery) AllConj
}

type imp_AllSimpleObjectFieldQuery nidhi.Query

func (q *imp_AllSimpleObjectFieldQuery) StringField(f *nidhi.StringQuery) AllConj {
	(*nidhi.Query)(q).Field("->>'stringField')", f)
	return (*imp_AllQuery)(q)
}

type AllNestedOneQuery interface {
	NestetedInt(*nidhi.IntQuery) AllConj
	Nested() AllNestedOneNestedQuery
}

type imp_AllNestedOneQuery nidhi.Query

func (q *imp_AllNestedOneQuery) NestetedInt(f *nidhi.IntQuery) AllConj {
	(*nidhi.Query)(q).Field("->'nestetedInt')::bigint", f)
	return (*imp_AllQuery)(q)
}

func (q *imp_AllNestedOneQuery) Nested() AllNestedOneNestedQuery {
	(*nidhi.Query)(q).Prefix("->'nested'")
	return (*imp_AllNestedOneNestedQuery)(q)
}

type AllNestedOneNestedQuery interface {
	SomeField(*nidhi.StringQuery) AllConj
	Nested(...*NestedThree) AllConj
}

type imp_AllNestedOneNestedQuery nidhi.Query

func (q *imp_AllNestedOneNestedQuery) SomeField(f *nidhi.StringQuery) AllConj {
	(*nidhi.Query)(q).Field("->>'someField')", f)
	return (*imp_AllQuery)(q)
}

func (q *imp_AllNestedOneNestedQuery) Nested(arr ...*NestedThree) AllConj {
	(*nidhi.Query)(q).Field(
		"->'nested'",
		nidhi.MarshalerQuery{
			Marshaler: NestedThreeSlice(arr),
		},
	)
	return (*imp_AllQuery)(q)
}

func (doc *All) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteString(w, "id", doc.Id, first)
	first = nidhigen.WriteString(w, "stringField", doc.StringField, first)
	first = nidhigen.WriteInt32(w, "int32Field", doc.Int32Field, first)
	first = nidhigen.WriteInt64(w, "int64Field", doc.Int64Field, first)
	first = nidhigen.WriteUint32(w, "uint32Field", doc.Uint32Field, first)
	first = nidhigen.WriteUint64(w, "uint64Field", doc.Uint64Field, first)
	first = nidhigen.WriteFloat32(w, "floatField", doc.FloatField, first)
	first = nidhigen.WriteFloat64(w, "doubleField", doc.DoubleField, first)
	first = nidhigen.WriteBool(w, "boolField", doc.BoolField, first)
	first = nidhigen.WriteBytes(w, "bytesField", doc.BytesField, first)

	first = nidhigen.WriteStringSlice(w, "primitiveRepeated", doc.PrimitiveRepeated, first)
	first = nidhigen.WriteMarshaler(w, "simpleObjectField", doc.SimpleObjectField, first)
	first = nidhigen.WriteMarshaler(w, "simpleRepeated", SimpleSlice(doc.SimpleRepeated), first)
	first = nidhigen.WriteMarshaler(w, "nestedOne", doc.NestedOne, first)
	first = nidhigen.WriteOneOf(w, doc.OneOf, first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *All) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "id":
			doc.Id = r.ReadString()
		case "stringField":
			doc.StringField = r.ReadString()
		case "int32Field":
			doc.Int32Field = r.ReadInt32()
		case "int64Field":
			doc.Int64Field = r.ReadInt64()
		case "uint32Field":
			doc.Uint32Field = r.ReadUint32()
		case "uint64Field":
			doc.Uint64Field = r.ReadUint64()
		case "floatField":
			doc.FloatField = r.ReadFloat32()
		case "doubleField":
			doc.DoubleField = r.ReadFloat64()
		case "boolField":
			doc.BoolField = r.ReadBool()
		case "bytesField":
			doc.BytesField = nidhigen.ReadByteSlice(r)
		case "primitiveRepeated":
			doc.PrimitiveRepeated = []string{}
			r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
				e := r.ReadString()
				doc.PrimitiveRepeated = append(doc.PrimitiveRepeated, e)
				return true
			})
		case "simpleObjectField":
			doc.SimpleObjectField = &Simple{}
			r.Error = doc.SimpleObjectField.UnmarshalDocument(r)
		case "simpleRepeated":
			doc.SimpleRepeated = []*Simple{}
			r.Error = (*SimpleSlice)(&doc.SimpleRepeated).UnmarshalDocument(r)
		case "nestedOne":
			doc.NestedOne = &NestedOne{}
			r.Error = doc.NestedOne.UnmarshalDocument(r)

		case "stringOneOf":
			var f All_StringOneOf
			f.StringOneOf = r.ReadString()

			doc.OneOf = &f
		case "int32OneOf":
			var f All_Int32OneOf
			f.Int32OneOf = r.ReadInt32()

			doc.OneOf = &f
		case "int64OneOf":
			var f All_Int64OneOf
			f.Int64OneOf = r.ReadInt64()

			doc.OneOf = &f
		case "uint32OneOf":
			var f All_Uint32OneOf
			f.Uint32OneOf = r.ReadUint32()

			doc.OneOf = &f
		case "uint64OneOf":
			var f All_Uint64OneOf
			f.Uint64OneOf = r.ReadUint64()

			doc.OneOf = &f
		case "floatOneOf":
			var f All_FloatOneOf
			f.FloatOneOf = r.ReadFloat32()

			doc.OneOf = &f
		case "doubleOneOf":
			var f All_DoubleOneOf
			f.DoubleOneOf = r.ReadFloat64()

			doc.OneOf = &f
		case "boolOneOf":
			var f All_BoolOneOf
			f.BoolOneOf = r.ReadBool()

			doc.OneOf = &f
		case "bytesOneOf":
			var f All_BytesOneOf
			f.BytesOneOf = nidhigen.ReadByteSlice(r)

			doc.OneOf = &f
		case "simpleObjectOneOf":
			var f All_SimpleObjectOneOf
			f.SimpleObjectOneOf = &Simple{}
			r.Error = f.SimpleObjectOneOf.UnmarshalDocument(r)

			doc.OneOf = &f
		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

func (of *All_StringOneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteStringOneOf(w, "stringOneOf", of.StringOneOf)
	return w.Error
}

func (of *All_Int32OneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteInt32OneOf(w, "int32OneOf", of.Int32OneOf)
	return w.Error
}

func (of *All_Int64OneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteInt64OneOf(w, "int64OneOf", of.Int64OneOf)
	return w.Error
}

func (of *All_Uint32OneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteUint32OneOf(w, "uint32OneOf", of.Uint32OneOf)
	return w.Error
}

func (of *All_Uint64OneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteUint64OneOf(w, "uint64OneOf", of.Uint64OneOf)
	return w.Error
}

func (of *All_FloatOneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteFloat32OneOf(w, "floatOneOf", of.FloatOneOf)
	return w.Error
}

func (of *All_DoubleOneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteFloat64OneOf(w, "doubleOneOf", of.DoubleOneOf)
	return w.Error
}

func (of *All_BoolOneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteBoolOneOf(w, "boolOneOf", of.BoolOneOf)
	return w.Error
}

func (of *All_BytesOneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteBytesOneOf(w, "bytesOneOf", of.BytesOneOf)

	return w.Error
}

func (of *All_SimpleObjectOneOf) MarshalDocument(w *jsoniter.Stream) error {
	nidhigen.WriteMarshalerOneOf(w, "simpleObjectOneOf", of.SimpleObjectOneOf)
	return w.Error
}

type AllSlice []*All

func (s AllSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *AllSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e All
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}

func (doc *Simple) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteString(w, "stringField", doc.StringField, first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *Simple) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "stringField":
			doc.StringField = r.ReadString()

		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type SimpleSlice []*Simple

func (s SimpleSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *SimpleSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e Simple
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}

func (doc *NestedOne) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteInt32(w, "nestetedInt", doc.NestetedInt, first)
	first = nidhigen.WriteMarshaler(w, "nested", doc.Nested, first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *NestedOne) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "nestetedInt":
			doc.NestetedInt = r.ReadInt32()
		case "nested":
			doc.Nested = &NestedTwo{}
			r.Error = doc.Nested.UnmarshalDocument(r)

		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type NestedOneSlice []*NestedOne

func (s NestedOneSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *NestedOneSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e NestedOne
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}

func (doc *NestedTwo) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteString(w, "someField", doc.SomeField, first)
	first = nidhigen.WriteMarshaler(w, "nested", NestedThreeSlice(doc.Nested), first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *NestedTwo) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "someField":
			doc.SomeField = r.ReadString()
		case "nested":
			doc.Nested = []*NestedThree{}
			r.Error = (*NestedThreeSlice)(&doc.Nested).UnmarshalDocument(r)

		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type NestedTwoSlice []*NestedTwo

func (s NestedTwoSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *NestedTwoSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e NestedTwo
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}

func (doc *NestedThree) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteString(w, "some", doc.Some, first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *NestedThree) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "some":
			doc.Some = r.ReadString()

		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type NestedThreeSlice []*NestedThree

func (s NestedThreeSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *NestedThreeSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e NestedThree
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}
