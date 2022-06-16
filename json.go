package nidhi

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/protojsoniter"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type bufferer interface {
	Buffer() []byte
}

type buffer []byte

func (b buffer) Buffer() []byte { return b }

func getJson(v any) (bufferer, error) {
	switch v := v.(type) {
	case protojsoniter.Writer:
		w := jsoniter.ConfigDefault.BorrowStream(nil)
		return w, w.Error
	case proto.Message:
		jsonData, err := protojson.Marshal(v)
		if err != nil {
			return nil, err
		}
		return buffer(jsonData), nil
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return buffer(jsonData), nil
	}
}

func putJson(v bufferer) {
	switch v := v.(type) {
	case *jsoniter.Stream:
		jsoniter.ConfigDefault.ReturnStream(v)
	}
}

func unmarshalJSON(b []byte, v any) error {
	switch v := v.(type) {
	case protojsoniter.Reader:
		r := jsoniter.ConfigDefault.BorrowIterator(b)
		defer jsoniter.ConfigDefault.ReturnIterator(r)
		v.ReadJSON(r)
		return r.Error
	case proto.Message:
		return protojson.Unmarshal(b, v)
	default:
		return json.Unmarshal(b, v)
	}
}
