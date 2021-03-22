package nidhigen

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type StringField nidhi.OrderByString

type IntField = nidhi.OrderByInt

type FloatField = nidhi.OrderByFloat

type BoolField = UnorderedField

type TimeField = nidhi.OrderByTime

type UnorderedField string

type ProtoMarshaler struct{ proto.Message }

func (m ProtoMarshaler) MarshalDocument(w *jsoniter.Stream) error {
	buf, err := protojson.Marshal(m.Message)
	if err != nil {
		return err
	}

	w.WriteVal(json.RawMessage(buf))
	return nil
}
