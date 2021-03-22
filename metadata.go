package nidhi

import (
	"context"

	jsoniter "github.com/json-iterator/go"
)

type MetadataCollection interface {
	Create(ctx context.Context, doc Document, ops []CreateOption) (string, error)
	Replace(ctx context.Context, doc Document, ops []ReplaceOption) error
	Update(ctx context.Context, doc Document, f Sqlizer, ops []UpdateOption) error
	Delete(ctx context.Context, id string, ops []DeleteOption) error
	DeleteMany(ctx context.Context, f Sqlizer, ops []DeleteOption) error
	Query(ctx context.Context, f Sqlizer, ctr func() Document, ops []QueryOption) error
	Get(ctx context.Context, id string, doc Unmarshaler, ops []GetOption) error
}

type MetadataProvider struct {
	Wrapper func(col MetadataCollection) MetadataCollection
	Keys    []string
}

type MetadataUnmarshaler interface {
	UnmarshalMetadata(key string, r *jsoniter.Iterator) (matched bool, err error)
}

type MetadataMarshaler interface {
	MarshalMetadata(w *jsoniter.Stream) error
}

type CreateMetadataFunc func() MetadataUnmarshaler

type mdMarshaler []MetadataMarshaler

func (md mdMarshaler) MarshalDocument(w *jsoniter.Stream) error {
	w.WriteObjectStart()
	for _, mm := range md {
		if err := mm.MarshalMetadata(w); err != nil {
			return err
		}
	}
	w.WriteObjectEnd()

	return w.Error
}

func (md mdMarshaler) UnmarshalDocument(_ *jsoniter.Iterator) error {
	panic("should not be called")
}

type mdUnmarshaler []MetadataUnmarshaler

func (md mdUnmarshaler) MarshalDocument(w *jsoniter.Stream) error {
	panic("should not be called")
}

func (md mdUnmarshaler) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadObjectCB(func(r *jsoniter.Iterator, s string) bool {
		for _, mu := range md {
			match, err := mu.UnmarshalMetadata(s, r)
			if err != nil {
				r.ReportError("metadata unmarshal for key: "+s, err.Error())
				return false
			}
			if match {
				return true
			}
		}

		r.Skip()
		return true
	})

	return r.Error
}
