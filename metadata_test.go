package nidhi_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi"
)

func init() {
	f := func() nidhi.MetadataValue {
		return &testDoc{}
	}
	nidhi.RegisterMetadataValueFactory("a", f)
	nidhi.RegisterMetadataValueFactory("b", f)
	nidhi.RegisterMetadataValueFactory("one", f)
}

func TestMetadata_MarshalDocument(t *testing.T) {
	tests := []struct {
		name    string
		md      nidhi.Metadata
		wantErr bool
	}{
		{
			"One Field",
			nidhi.Metadata{
				"one": &testDoc{Id: "id"},
			},
			false,
		},
		{
			"Two Fields",
			nidhi.Metadata{
				"a": &testDoc{Id: "id"},
				"b": &testDoc{Id: "id"},
			},
			false,
		},
		{
			"Nil",
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := jsoniter.ConfigDefault.BorrowStream(nil)
			defer jsoniter.ConfigDefault.ReturnStream(stream)

			if err := tt.md.MarshalDocument(stream); (err != nil) != tt.wantErr {
				t.Errorf("Metadata.MarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
			}

			exp, err := json.Marshal(tt.md)
			if err != nil {
				panic(err)
			}
			if tt.md == nil {
				exp = []byte("{}")
			}

			act := stream.Buffer()

			if !tt.wantErr && !bytes.Equal(exp, act) {
				t.Errorf("Metadata.MarshalDocument() exp: %s, act: %s", exp, act)
			}
		})
	}
}

func TestMetadata_UnmarshalDocument(t *testing.T) {
	tests := []struct {
		name    string
		md      nidhi.Metadata
		wantErr bool
	}{
		{
			"One Field",
			nidhi.Metadata{
				"one": &testDoc{Id: "id"},
			},
			false,
		},
		{
			"Two Fields",
			nidhi.Metadata{
				"a": &testDoc{Id: "id"},
				"b": &testDoc{Id: "id"},
			},
			false,
		},
		{
			"Nil",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.md)
			if err != nil {
				panic(err)
			}

			iter := jsoniter.ConfigDefault.BorrowIterator(data)
			defer jsoniter.ConfigDefault.ReturnIterator(iter)

			act := nidhi.Metadata{}
			if tt.md == nil {
				act = nil
			}
			if err := act.UnmarshalDocument(iter); (err != nil) != tt.wantErr {
				t.Errorf("Metadata.UnmarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(act, tt.md) {
				t.Errorf("Metadata.UnmarshalDocument() exp: %v, act: %v", tt.md, act)
			}
		})
	}
}
