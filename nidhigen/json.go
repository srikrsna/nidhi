package nidhigen

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

func WriteStringSlice(w *jsoniter.Stream, field string, value []string, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteString(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteString(v)
	}
	w.WriteArrayEnd()

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

func WriteBoolSlice(w *jsoniter.Stream, field string, value []bool, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteBool(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteBool(v)
	}
	w.WriteArrayEnd()

	return false
}

func WriteUint32(w *jsoniter.Stream, field string, value uint32, first bool) bool {
	if value == 0 {
		return false
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteUint32(value)

	return false
}

func WriteUint32Slice(w *jsoniter.Stream, field string, value []uint32, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteUint32(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteUint32(v)
	}
	w.WriteArrayEnd()

	return false
}

func WriteUint64(w *jsoniter.Stream, field string, value uint64, first bool) bool {
	if value == 0 {
		return false
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteUint64(value)

	return false
}

func WriteUint64Slice(w *jsoniter.Stream, field string, value []uint64, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteUint64(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteUint64(v)
	}
	w.WriteArrayEnd()

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

func WriteIntSlice(w *jsoniter.Stream, field string, value []int, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteInt(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteInt(v)
	}
	w.WriteArrayEnd()

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

func WriteFloat32Slice(w *jsoniter.Stream, field string, value []float32, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteFloat32(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteFloat32(v)
	}
	w.WriteArrayEnd()

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

func WriteFloat64Slice(w *jsoniter.Stream, field string, value []float64, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteFloat64(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteFloat64(v)
	}
	w.WriteArrayEnd()

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

func WriteInt32Slice(w *jsoniter.Stream, field string, value []int32, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteInt32(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteInt32(v)
	}
	w.WriteArrayEnd()

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

func WriteInt64Slice(w *jsoniter.Stream, field string, value []int64, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)

	w.WriteArrayStart()
	w.WriteInt64(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteInt64(v)
	}
	w.WriteArrayEnd()

	return false
}

func WriteMarshaler(w *jsoniter.Stream, field string, value nidhi.Marshaler, first bool) bool {
	if value == nil {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.Error = value.MarshalDocument(w)

	return false
}

func WriteBytes(w *jsoniter.Stream, field string, value []byte, first bool) bool {
	if value == nil {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteVal(value)

	return false
}

func WriteOneOf(w *jsoniter.Stream, field interface{}, first bool) bool {
	if field == nil {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.Error = field.(nidhi.Marshaler).MarshalDocument(w)

	return false
}
