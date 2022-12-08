package nidhi

import (
	"strings"

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

// GetMetadataField returns the field accessor query for the given field.
//
// It is intented to be used by metadata packages.
func GetMetadataField(part string, typ string, def string) string {
	op := `->`
	if strings.ToUpper(typ) == "TEXT" {
		op = `->>`
	}
	// Eg: metadata->'activity'::jsonb
	return `(` + ColMeta + op + `'` + part + `')::` + typ
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
