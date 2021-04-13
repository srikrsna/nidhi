package nidhigen

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/srikrsna/nidhi"
)

const (
	ErrNilUnmarshal nidhi.Error = "nil passed to unmarshal"
)

func WriteString(w *jsoniter.Stream, field, value string, first bool) bool {
	if value == "" {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteStringOneOf(w, field, value)

	return false
}

func WriteStringSlice(w *jsoniter.Stream, field string, value []string, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteStringSliceOneOf(w, field, value)

	return false
}

func WriteBool(w *jsoniter.Stream, field string, value, first bool) bool {
	if !value {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteBoolOneOf(w, field, value)

	return false
}

func WriteBoolSlice(w *jsoniter.Stream, field string, value []bool, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteBoolSliceOneOf(w, field, value)
	return false
}

func WriteUint32(w *jsoniter.Stream, field string, value uint32, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteUint32OneOf(w, field, value)

	return false
}

func WriteUint32Slice(w *jsoniter.Stream, field string, value []uint32, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteUint32SliceOneOf(w, field, value)

	return false
}

func WriteUint64(w *jsoniter.Stream, field string, value uint64, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteUint64OneOf(w, field, value)

	return false
}

func WriteUint64Slice(w *jsoniter.Stream, field string, value []uint64, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteUint64SliceOneOf(w, field, value)

	return false
}

func WriteInt(w *jsoniter.Stream, field string, value int, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteIntOneOf(w, field, value)

	return false
}

func WriteIntSlice(w *jsoniter.Stream, field string, value []int, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteIntSliceOneOf(w, field, value)

	return false
}

func WriteFloat32(w *jsoniter.Stream, field string, value float32, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteFloat32OneOf(w, field, value)

	return false
}

func WriteFloat32Slice(w *jsoniter.Stream, field string, value []float32, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteFloat32SliceOneOf(w, field, value)

	return false
}

func WriteFloat64(w *jsoniter.Stream, field string, value float64, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteFloat64OneOf(w, field, value)

	return false
}

func WriteFloat64Slice(w *jsoniter.Stream, field string, value []float64, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteFloat64SliceOneOf(w, field, value)

	return false
}

func WriteInt32(w *jsoniter.Stream, field string, value int32, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}
	WriteInt32OneOf(w, field, value)

	return false
}

func WriteInt32Slice(w *jsoniter.Stream, field string, value []int32, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteInt32SliceOneOf(w, field, value)

	return false
}

func WriteInt64(w *jsoniter.Stream, field string, value int64, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteInt64OneOf(w, field, value)

	return false
}

func WriteInt64Slice(w *jsoniter.Stream, field string, value []int64, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteInt64SliceOneOf(w, field, value)

	return false
}

func WriteFieldMask(w *jsoniter.Stream, field string, value *fieldmaskpb.FieldMask, first bool) bool {
	if len(value.GetPaths()) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteFieldMaskOneOf(w, field, value)

	return false
}

func WriteTimestamp(w *jsoniter.Stream, field string, value *timestamppb.Timestamp, first bool) bool {
	if value == nil {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteTimestampOneOf(w, field, value)

	return false
}

func WriteTimestampSlice(w *jsoniter.Stream, field string, value []*timestamppb.Timestamp, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteTimestampSliceOneOf(w, field, value)

	return false
}

func WriteDuration(w *jsoniter.Stream, field string, value *durationpb.Duration, first bool) bool {
	if value == nil {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteDurationOneOf(w, field, value)

	return false
}

func WriteDurationSlice(w *jsoniter.Stream, field string, value []*durationpb.Duration, first bool) bool {
	if value == nil {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteDurationSliceOneOf(w, field, value)

	return false
}

func WriteAny(w *jsoniter.Stream, field string, value *anypb.Any, first bool) bool {
	if value == nil {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteAnyOneOf(w, field, value)

	return false
}

func WriteAnySlice(w *jsoniter.Stream, field string, value []*anypb.Any, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteAnySliceOneOf(w, field, value)

	return false
}

func WriteMarshaler(w *jsoniter.Stream, field string, value nidhi.Marshaler, first bool) bool {
	if value == nil || reflect.ValueOf(value).IsNil() {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteMarshalerOneOf(w, field, value)

	return false
}

func WriteBytes(w *jsoniter.Stream, field string, value []byte, first bool) bool {
	if len(value) == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	WriteBytesOneOf(w, field, value)

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

func WriteStringOneOf(w *jsoniter.Stream, field, value string) {
	w.WriteObjectField(field)
	w.WriteString(value)
}

func WriteStringSliceOneOf(w *jsoniter.Stream, field string, value []string) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteString(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteString(v)
	}
	w.WriteArrayEnd()
}

func WriteBoolOneOf(w *jsoniter.Stream, field string, value bool) {
	w.WriteObjectField(field)
	w.WriteBool(value)
}

func WriteBoolSliceOneOf(w *jsoniter.Stream, field string, value []bool) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteBool(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteBool(v)
	}
	w.WriteArrayEnd()

}

func WriteUint32OneOf(w *jsoniter.Stream, field string, value uint32) {
	w.WriteObjectField(field)
	w.WriteUint32(value)
}

func WriteUint32SliceOneOf(w *jsoniter.Stream, field string, value []uint32) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteUint32(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteUint32(v)
	}
	w.WriteArrayEnd()
}

func WriteUint64OneOf(w *jsoniter.Stream, field string, value uint64) {
	w.WriteObjectField(field)
	w.WriteUint64(value)
}

func WriteUint64SliceOneOf(w *jsoniter.Stream, field string, value []uint64) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteUint64(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteUint64(v)
	}
	w.WriteArrayEnd()
}

func WriteIntOneOf(w *jsoniter.Stream, field string, value int) {
	w.WriteObjectField(field)
	w.WriteInt(value)
}

func WriteIntSliceOneOf(w *jsoniter.Stream, field string, value []int) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteInt(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteInt(v)
	}
	w.WriteArrayEnd()
}

func WriteFloat32OneOf(w *jsoniter.Stream, field string, value float32) {
	w.WriteObjectField(field)
	w.WriteFloat32(value)
}

func WriteFloat32SliceOneOf(w *jsoniter.Stream, field string, value []float32) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteFloat32(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteFloat32(v)
	}
	w.WriteArrayEnd()
}

func WriteFloat64OneOf(w *jsoniter.Stream, field string, value float64) {
	w.WriteObjectField(field)
	w.WriteFloat64(value)
}

func WriteFloat64SliceOneOf(w *jsoniter.Stream, field string, value []float64) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteFloat64(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteFloat64(v)
	}
	w.WriteArrayEnd()
}

func WriteInt32OneOf(w *jsoniter.Stream, field string, value int32) {
	w.WriteObjectField(field)
	w.WriteInt32(value)
}

func WriteInt32SliceOneOf(w *jsoniter.Stream, field string, value []int32) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteInt32(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteInt32(v)
	}
	w.WriteArrayEnd()
}

func WriteInt64OneOf(w *jsoniter.Stream, field string, value int64) {
	w.WriteObjectField(field)
	w.WriteInt64(value)
}

func WriteInt64SliceOneOf(w *jsoniter.Stream, field string, value []int64) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteInt64(value[0])
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteInt64(v)
	}
	w.WriteArrayEnd()
}

func WriteFieldMaskOneOf(w *jsoniter.Stream, field string, value *fieldmaskpb.FieldMask) {
	w.WriteObjectField(field)
	w.WriteString(strings.Join(value.GetPaths(), ","))
}

func WriteTimestampOneOf(w *jsoniter.Stream, field string, value *timestamppb.Timestamp) {
	w.WriteObjectField(field)
	w.WriteString(value.AsTime().Format(time.RFC3339))
}

func WriteTimestampSliceOneOf(w *jsoniter.Stream, field string, value []*timestamppb.Timestamp) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteString(value[0].AsTime().Format(time.RFC3339))
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteString(v.AsTime().Format(time.RFC3339))
	}
	w.WriteArrayEnd()
}

func WriteDurationOneOf(w *jsoniter.Stream, field string, value *durationpb.Duration) {
	w.WriteObjectField(field)
	w.WriteInt64(value.GetSeconds())
}

func WriteDurationSliceOneOf(w *jsoniter.Stream, field string, value []*durationpb.Duration) {
	w.WriteObjectField(field)

	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	w.WriteInt64(value[0].GetSeconds())
	for _, v := range value[1:] {
		w.WriteMore()
		w.WriteInt64(v.GetSeconds())
	}
	w.WriteArrayEnd()
}

func WriteAnyOneOf(w *jsoniter.Stream, field string, value *anypb.Any) {
	w.WriteObjectField(field)
	buf, err := protojson.Marshal(value)
	if err != nil {
		w.Error = err
		return
	}

	w.WriteVal(json.RawMessage(buf))
}

func WriteAnySliceOneOf(w *jsoniter.Stream, field string, value []*anypb.Any) {
	w.WriteObjectField(field)
	if len(value) == 0 {
		w.WriteEmptyArray()
		return
	}

	w.WriteArrayStart()
	buf, err := protojson.Marshal(value[0])
	if err != nil {
		w.Error = err
		return
	}
	w.WriteVal(json.RawMessage(buf))
	for _, v := range value[1:] {
		w.WriteMore()
		buf, err := protojson.Marshal(v)
		if err != nil {
			w.Error = err
			return
		}
		w.WriteVal(json.RawMessage(buf))
	}
	w.WriteArrayEnd()
}

func WriteMarshalerOneOf(w *jsoniter.Stream, field string, value nidhi.Marshaler) {
	w.WriteObjectField(field)
	if value == nil || reflect.ValueOf(value).IsNil() {
		w.WriteNil()
		return
	}

	w.Error = value.MarshalDocument(w)
}

func WriteBytesOneOf(w *jsoniter.Stream, field string, value []byte) {
	w.WriteObjectField(field)
	w.WriteString(base64.StdEncoding.EncodeToString(value))
}

func ReadByteSlice(r *jsoniter.Iterator) []byte {
	v, err := base64.StdEncoding.DecodeString(r.ReadString())
	if err != nil {
		r.ReportError("decoding byte slice", err.Error())
	}

	return v
}

func ReadTimestamp(r *jsoniter.Iterator) *timestamppb.Timestamp {
	t, err := time.Parse(time.RFC3339, r.ReadString())
	if err != nil {
		r.ReportError("decoding timestamp", err.Error())
	}

	return timestamppb.New(t)
}

func ReadDuration(r *jsoniter.Iterator) *durationpb.Duration {
	return &durationpb.Duration{
		Seconds: r.ReadInt64(),
	}
}

func ReadAny(r *jsoniter.Iterator) *anypb.Any {
	var buf json.RawMessage
	r.ReadVal(&buf)

	var any anypb.Any
	if err := protojson.Unmarshal([]byte(buf), &any); err != nil {
		r.ReportError("decoding any", err.Error())
	}
	return &any
}

func ReadFieldMask(r *jsoniter.Iterator) *fieldmaskpb.FieldMask {
	return &fieldmaskpb.FieldMask{
		Paths: strings.Split(r.ReadString(), ","),
	}
}

func TimestampSliceToArgs(in []*timestamppb.Timestamp) []string {
	res := make([]string, 0, len(in))
	for _, t := range in {
		res = append(res, t.AsTime().Format(time.RFC3339))
	}
	return res
}

func DurationSliceToArgs(in []*durationpb.Duration) []int64 {
	res := make([]int64, 0, len(in))
	for _, d := range in {
		res = append(res, d.GetSeconds())
	}

	return res
}
