package nidhi

import (
	jsoniter "github.com/json-iterator/go"
)

type Document interface {
	Marshaler
	Unmarshaler

	SetDocumentId(id string)
	DocumentId() string
}

type Marshaler interface {
	MarshalDocument(w *jsoniter.Stream) error
}

type Unmarshaler interface {
	UnmarshalDocument(r *jsoniter.Iterator) error
}

func WriteString(w *jsoniter.Stream, field, value string) {
	if value != "" {
		w.WriteObjectField(field)
		w.WriteString(value)
	}
}

func WriteBool(w *jsoniter.Stream, field string, value bool) {
	if value {
		w.WriteObjectField(field)
		w.WriteBool(value)
	}
}

func WriteInt(w *jsoniter.Stream, field string, value int) {
	if value != 0 {
		w.WriteObjectField(field)
		w.WriteInt(value)
	}
}

func WriteFloat32(w *jsoniter.Stream, field string, value float32) {
	if value != 0 {
		w.WriteObjectField(field)
		w.WriteFloat32(value)
	}
}

func WriteFloat64(w *jsoniter.Stream, field string, value float64) {
	if value != 0 {
		w.WriteObjectField(field)
		w.WriteFloat64(value)
	}
}

func WriteInt32(w *jsoniter.Stream, field string, value int32) {
	if value != 0 {
		w.WriteObjectField(field)
		w.WriteInt32(value)
	}
}

func WriteInt64(w *jsoniter.Stream, field string, value int64) {
	if value != 0 {
		w.WriteObjectField(field)
		w.WriteInt64(value)
	}
}

func WriteMarshaler(w *jsoniter.Stream, field string, value Marshaler) {
	if value != nil {
		w.WriteObjectField(field)
		w.Error = value.MarshalDocument(w)
	}
}
