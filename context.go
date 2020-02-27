package nidhi

import "context"

type contextKey int

const (
	metadataKey contextKey = iota
)

type Metadata map[string]interface{}

func (m Metadata) Set(k string, v interface{}) {
	m[k] = v
}

func (m Metadata) Get(k string) interface{} {
	return m[k]
}

func ExtractMetadata(ctx context.Context) Metadata {
	v := ctx.Value(metadataKey)
	if v != nil {
		return v.(Metadata)
	}

	return Metadata{}
}

func WithMetadata(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataKey, md)
}
