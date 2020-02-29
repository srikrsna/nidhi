package nidhi_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi"
)

type testDoc struct {
	Id string `json:"Id,omitempty"`
}

func (doc *testDoc) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}
	w.WriteObjectStart()
	if doc.Id != "" {
		w.WriteObjectField("Id")
		w.WriteString(doc.Id)
	}
	w.WriteObjectEnd()
	return w.Error
}

func (doc *testDoc) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("nil document passed for unmarshal")
	}
	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "Id":
			doc.Id = r.ReadString()
		default:
			r.Skip()
		}
		return true
	})
	return r.Error
}

func Test_jsonb_Scan(t *testing.T) {
	must := func(d []byte, err error) []byte {
		if err != nil {
			panic("err must not be nil")
		}

		return d
	}
	type fields struct {
		v interface {
			nidhi.Marshaler
			nidhi.Unmarshaler
		}
	}
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"One Filled",
			fields{&testDoc{}},
			args{&testDoc{Id: "This is an identifier"}},
			false,
		},
		{
			"Empty",
			fields{&testDoc{}},
			args{&testDoc{}},
			false,
		},
		{
			"Nil Unmarshal",
			fields{nil},
			args{&testDoc{Id: "This is an identifier"}},
			true,
		},
		{
			"Document Nil Unmarshal",
			fields{(*testDoc)(nil)},
			args{&testDoc{Id: "This is an identifier"}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := nidhi.JSONB(tt.fields.v)

			if err := j.Scan(must(json.Marshal(tt.args.src))); (err != nil) != tt.wantErr {
				t.Errorf("jsonb.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.fields.v, tt.args.src) {
				t.Errorf("jsnob Scan() act: %v, exp: %v", tt.fields.v, tt.args.src)
			}
		})
	}
}

func Test_jsonb_Value(t *testing.T) {
	must := func(d []byte, err error) []byte {
		if err != nil {
			panic("err must not be nil")
		}

		return d
	}
	type fields struct {
		v interface {
			nidhi.Marshaler
			nidhi.Unmarshaler
		}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"One Filled",
			fields{&testDoc{Id: "This is an ID"}},
			false,
		},
		{
			"Empty",
			fields{&testDoc{}},
			false,
		},
		{
			"Nil Marshal",
			fields{nil},
			false,
		},
		{
			"Document Nil Marshal",
			fields{(*testDoc)(nil)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := nidhi.JSONB(tt.fields.v)

			data, err := j.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonb.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			act, ok := data.([]byte)
			if !tt.wantErr && !ok {
				t.Error("jsonb.Value() must always return a []byte")
			}

			exp := must(json.Marshal(tt.fields.v))

			if !tt.wantErr && !bytes.Equal(act, exp) {
				t.Errorf("jsnob Scan() act: %s, exp: %s", act, exp)
			}
		})
	}
}
