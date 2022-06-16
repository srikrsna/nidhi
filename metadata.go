package nidhi

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/protojsoniter"
)

var (
	_ protojsoniter.Reader = Metadata{}
	_ protojsoniter.Writer = Metadata{}
)

// Metadata is the metadata stored against a document.
//
// Metadata is designed to be extendible. These should typically be document agnostic information.
// Eg:
// 	- Activity log (create, update timestamps)
//  - Common information extracted from the document
type Metadata map[string]MetadataPart

type MetadataPart interface {
	// MarshalMD marshals the metadata.
	MarshalMDP(w *jsoniter.Stream)
	// UnmarshalMD unmarshals the metdata.
	UnmarshalMDP(r *jsoniter.Iterator)
}

func (m Metadata) WriteJSON(w *jsoniter.Stream) {
	w.WriteObjectStart()
	var count int
	for k, v := range m {
		count++
		w.WriteObjectField(k)
		v.MarshalMDP(w)
		if count != len(m) {
			w.WriteMore()
		}
	}
	w.WriteObjectEnd()
}

func (m Metadata) ReadJSON(r *jsoniter.Iterator) {
	for field := r.ReadObject(); field != ""; field = r.ReadObject() {
		part, ok := m[field]
		if ok {
			part.UnmarshalMDP(r)
		}
	}
}
