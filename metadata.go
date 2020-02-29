package nidhi

import (
	"errors"

	sq "github.com/elgris/sqrl"
	jsoniter "github.com/json-iterator/go"
)

type MetadataValue interface {
	Marshaler
	Unmarshaler
}

type Metadata map[string]MetadataValue

func (md Metadata) MarshalDocument(w *jsoniter.Stream) error {
	w.WriteObjectStart()
	l := len(md)
	for k, v := range md {
		w.WriteObjectField(k)
		v.MarshalDocument(w)
		if l--; l > 0 {
			w.WriteMore()
		}
	}
	w.WriteObjectEnd()

	return w.Error
}

func (md Metadata) UnmarshalDocument(r *jsoniter.Iterator) error {
	if md == nil {
		return errors.New("nidhi: emtpy metdata cannot be unmarshalled")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		vf, ok := factories[field]
		if ok {
			v := vf()
			v.UnmarshalDocument(r)
			md[field] = v
		} else {
			r.Skip()
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

type MetadataValueFactory func() MetadataValue

var factories = map[string]MetadataValueFactory{}

func RegisterMetadataValueFactory(key string, f MetadataValueFactory) {
	if _, ok := factories[key]; ok {
		panic("already registered Metadata factory with the same name")
	}

	factories[key] = f
}
