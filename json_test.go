package nidhi_test

import (
	"testing"

	"github.com/akshayjshah/attest"
	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
)

type jsonStdLib jsoniterReadWriter

type jsoniterReadWriter struct {
	Id string `json:"id,omitempty"`
}

func (v *jsoniterReadWriter) WriteJSON(w *jsoniter.Stream) {
	if v == nil {
		w.WriteEmptyObject()
		return
	}
	w.WriteObjectStart()
	if v.Id != "" {
		w.WriteObjectField("name")
		w.WriteString(v.Id)
	}
	w.WriteObjectEnd()
}

func (v *jsoniterReadWriter) ReadJSON(r *jsoniter.Iterator) {
	for field := r.ReadObject(); field != ""; field = r.ReadObject() {
		switch field {
		case "name":
			v.Id = r.ReadString()
		default:
			r.Skip()
		}
	}
}

func TestGetJson(t *testing.T) {
	t.Parallel()
	t.Run("metadata", func(t *testing.T) {
		t.Parallel()
		md := nidhi.Metadata{"part": &metadataPart{Value: "some"}, "second": &metadataPart{}}
		w, err := nidhi.GetJson(md)
		attest.Ok(t, err)
		attest.True(t, string(w.Buffer()) == `{"part":{"value":"some"},"second":{}}` || string(w.Buffer()) == `{"second":{},"part":{"value":"some"}}`)
		emd := nidhi.Metadata{"part": &metadataPart{}, "second": &metadataPart{}}
		attest.Ok(t, nidhi.UnmarshalJson(w.Buffer(), emd))
		attest.Equal(t, emd, md)
	})
	t.Run("jsoniter", func(t *testing.T) {
		t.Parallel()
		rw := &jsoniterReadWriter{Id: "some"}
		w, err := nidhi.GetJson(rw)
		attest.Ok(t, err)
		attest.Equal(t, string(w.Buffer()), `{"name":"some"}`)
		var grw jsoniterReadWriter
		attest.Ok(t, nidhi.UnmarshalJson(w.Buffer(), &grw))
		attest.Equal(t, &grw, rw)
	})
	t.Run("proto", func(t *testing.T) {
		t.Parallel()
		pt, err := anypb.New(&anypb.Any{})
		attest.Ok(t, err)
		w, err := nidhi.GetJson(pt)
		attest.Ok(t, err)
		var pany anypb.Any
		protojson.Unmarshal(w.Buffer(), &pany)
		attest.Equal(t, &pany, pt, attest.Cmp(protocmp.Transform()))
		var gpt anypb.Any
		attest.Ok(t, nidhi.UnmarshalJson(w.Buffer(), &gpt))
		attest.Equal(t, &gpt, pt, attest.Cmp(protocmp.Transform()))
	})
	t.Run("stdlib", func(t *testing.T) {
		t.Parallel()
		rw := &jsonStdLib{Id: "some"}
		w, err := nidhi.GetJson(rw)
		attest.Ok(t, err)
		attest.Equal(t, string(w.Buffer()), `{"id":"some"}`)
		var grw jsonStdLib
		attest.Ok(t, nidhi.UnmarshalJson(w.Buffer(), &grw))
		attest.Equal(t, &grw, rw)
	})
}
