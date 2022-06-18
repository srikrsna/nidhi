package nidhi_test

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi"
)

var (
	_ nidhi.MetadataPart = (*metadataPart)(nil)
)

type metadataPart struct {
	Value string `json:"value,omitempty"`
}

// MarshalMD marshals the metadata.
func (m *metadataPart) MarshalMDP(w *jsoniter.Stream) {
	w.WriteObjectStart()
	if m.Value != "" {
		w.WriteObjectField("value")
		w.WriteString(m.Value)
	}
	w.WriteObjectEnd()
}

// UnmarshalMD unmarshals the metdata.
func (m *metadataPart) UnmarshalMDP(r *jsoniter.Iterator) {
	for field := r.ReadObject(); field != ""; field = r.ReadObject() {
		switch field {
		case "value":
			m.Value = r.ReadString()
		default:
			r.Skip()
		}
	}
}
