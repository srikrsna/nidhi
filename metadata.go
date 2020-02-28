package nidhi

import (
	"errors"

	sq "github.com/elgris/sqrl"
	jsoniter "github.com/json-iterator/go"
)

type Metadata map[string]interface {
	Marshaler
	Unmarshaler
}

func (m Metadata) Set(k string, v interface {
	Marshaler
	Unmarshaler
}) {
	m[k] = v
}

func (m Metadata) Get(k string) interface {
	Marshaler
	Unmarshaler
} {
	return m[k]
}

func (md Metadata) MarshalDocument(w *jsoniter.Stream) error {
	w.WriteObjectStart()
	for k, v := range md {
		w.WriteObjectField(k)
		v.MarshalDocument(w)
	}
	w.WriteObjectEnd()

	return w.Error
}

func (md Metadata) UnmarshalDocument(r *jsoniter.Iterator) error {
	if md == nil {
		return errors.New("nidhi: emtpy metdata cannot be unmarshalled")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		v := md[field]
		if v != nil {
			v.UnmarshalDocument(r)
		}
		return true
	})

	return r.Error
}

func updateMetadata(st *sq.UpdateBuilder, md Metadata, replace bool) *sq.UpdateBuilder {
	if md == nil {
		return st
	}

	if replace {
		return st.Set(metaCol, JSONB(md))
	}

	return st.Set(metaCol, sq.Expr(metaCol+" || "+"?", JSONB(md)))
}
