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
		v, ok := md[field]
		if ok {
			v.UnmarshalDocument(r)
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
