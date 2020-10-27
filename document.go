package nidhi

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
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
