package nidhi

import (
	jsoniter "github.com/json-iterator/go"
)

// Metadata is the metadata stored against a document.
//
// Metadata is designed to be extendible. These should typically be document agnostic information.
// Eg:
//   - Activity log (create, update timestamps)
//   - Common information extracted from the document
type Metadata map[string]MetadataPart

type MetadataPart interface {
	// MarshalMD marshals the metadata.
	MarshalMDP(w *jsoniter.Stream)
	// UnmarshalMD unmarshals the metdata.
	UnmarshalMDP(r *jsoniter.Iterator)
}

// MetadataField is a [Field] for metadata parts
//
// Metadata parts can return a MetadataField
type MetadataField struct {
	Part    string
	Type    string
	Default string
}

func (f *MetadataField) Selector() string {
	return `JSON_VALUE(` + ColMeta + `::jsonb, '$.` + f.Part + `' RETURNING ` + f.Type + ` DEFAULT ` + f.Default + ` ON EMPTY)`
}

func (m Metadata) writeJSON(w *jsoniter.Stream) {
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

func (m Metadata) readJSON(r *jsoniter.Iterator) {
	for field := r.ReadObject(); field != ""; field = r.ReadObject() {
		part, ok := m[field]
		if ok {
			part.UnmarshalMDP(r)
		} else {
			r.Skip()
		}
	}
}
