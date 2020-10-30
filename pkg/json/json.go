package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi"
)

func WriteString(w *jsoniter.Stream, field, value string, first bool) bool {
	if value == "" {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteString(value)

	return false
}

func WriteBool(w *jsoniter.Stream, field string, value, first bool) bool {
	if !value {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteBool(value)

	return false
}

func WriteInt(w *jsoniter.Stream, field string, value int, first bool) bool {
	if value == 0 {
		return false
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteInt(value)

	return false
}

func WriteFloat32(w *jsoniter.Stream, field string, value float32, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteFloat32(value)

	return false
}

func WriteFloat64(w *jsoniter.Stream, field string, value float64, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteFloat64(value)

	return false
}

func WriteInt32(w *jsoniter.Stream, field string, value int32, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}
	w.WriteObjectField(field)
	w.WriteInt32(value)

	return false
}

func WriteInt64(w *jsoniter.Stream, field string, value int64, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteInt64(value)

	return false
}

func WriteMarshaler(w *jsoniter.Stream, field string, value nidhi.Marshaler, first bool) bool {
	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	value.MarshalDocument(w)

	return false
}
