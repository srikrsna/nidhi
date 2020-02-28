package nidhi

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type Document interface {
	Marshaler
	Unmarshaler

	DocumentId() string
}

type Marshaler interface {
	MarshalDocument(w *jsoniter.Stream) error
}

type Unmarshaler interface {
	UnmarshalDocument(r *jsoniter.Iterator) error
}

type jsonb struct {
	v interface {
		Marshaler
		Unmarshaler
	}
}

func JSONB(v interface {
	Marshaler
	Unmarshaler
}) interface {
	driver.Valuer
	sql.Scanner
} {
	return &jsonb{v}
}

func (j jsonb) Scan(src interface{}) error {
	dat, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("nidhi: seems to be a bug, error while scanning jsonb: expected []byte got %T", src)
	}

	iter := jsoniter.ConfigDefault.BorrowIterator(dat)
	defer jsoniter.ConfigDefault.ReturnIterator(iter)

	return j.v.UnmarshalDocument(iter)
}

func (j jsonb) Value() (driver.Value, error) {
	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	if err := j.v.MarshalDocument(stream); err != nil {
		return nil, err
	}

	data := make([]byte, stream.Buffered())
	copy(data, stream.Buffer())

	return data, nil
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
