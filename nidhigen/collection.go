package nidhigen

import (
	"context"

	"github.com/srikrsna/nidhi"
)

type Collection interface {
	Create(ctx context.Context, doc nidhi.Document, ops []nidhi.CreateOption) (string, error)
	Replace(ctx context.Context, doc nidhi.Document, ops []nidhi.ReplaceOption) error
	Update(ctx context.Context, doc nidhi.Document, f nidhi.Sqlizer, ops []nidhi.UpdateOption) error
	Delete(ctx context.Context, id string, ops []nidhi.DeleteOption) error
	DeleteMany(ctx context.Context, f nidhi.Sqlizer, ops []nidhi.DeleteOption) error
	Query(ctx context.Context, f nidhi.Sqlizer, ctr func() nidhi.Document, ops []nidhi.QueryOption) error
	Get(ctx context.Context, id string, doc nidhi.Unmarshaler, ops []nidhi.GetOption) error
	Count(ctx context.Context, f nidhi.Sqlizer, ops []nidhi.CountOption) (int64, error)
}
