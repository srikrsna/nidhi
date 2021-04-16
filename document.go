package nidhi

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
)

type Document interface {
	Marshaler
	Unmarshaler

	DocumentId() string
	SetDocumentId(id string)
}

type Marshaler interface {
	MarshalDocument(w *jsoniter.Stream) error
}

type Unmarshaler interface {
	UnmarshalDocument(r *jsoniter.Iterator) error
}

type Jsonb struct {
	V interface{}
}

func (j Jsonb) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return jsoniter.Unmarshal([]byte(v), j.V)
	case []byte:
		return jsoniter.Unmarshal(v, j.V)
	}

	return fmt.Errorf("nidhi: error while scanning jsonb, expected []byte got %T", src)
}

func (j Jsonb) Value() (driver.Value, error) {
	return jsoniter.Marshal(j.V)
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
	if j.v == nil {
		return errors.New("nidhi: nil passed into JSONB. cannot scan into nil")
	}

	dat, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("nidhi: error while scanning jsonb, expected []byte got %T", src)
	}

	iter := jsoniter.ConfigDefault.BorrowIterator(dat)
	defer jsoniter.ConfigDefault.ReturnIterator(iter)

	return j.v.UnmarshalDocument(iter)
}

func (j jsonb) Value() (driver.Value, error) {
	if j.v == nil {
		return []byte("null"), nil
	}

	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	if err := j.v.MarshalDocument(stream); err != nil {
		return nil, err
	}

	data := make([]byte, stream.Buffered())
	copy(data, stream.Buffer())

	return data, nil
}

type NoopUnmarshaler struct {
	Marshaler
}

func (NoopUnmarshaler) UnmarshalDocument(_ *jsoniter.Iterator) error {
	return nil
}

type NoopMarshaler struct {
	Unmarshaler
}

func (NoopMarshaler) MarshalDocument(w *jsoniter.Stream) error {
	return nil
}

type JsonbArray []Marshaler

func (j JsonbArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	var sb strings.Builder
	stream := jsoniter.ConfigDefault.BorrowStream(&sb)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	ba := make(pq.StringArray, 0, len(j))
	for _, v := range j {
		sb.Reset()
		stream.Reset(&sb)
		if err := v.MarshalDocument(stream); err != nil {
			return nil, err
		}
		stream.Flush()

		ba = append(ba, sb.String())
	}

	return ba.Value()
}
