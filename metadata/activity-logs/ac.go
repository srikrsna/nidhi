// Package activitylogs ...
package activitylogs

import (
	"context"
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/srikrsna/nidhi"
)

const (
	CreatedKey = "created"
	UpdatedKey = "updated"
	DeletedKey = "deleted"
)

func Provider(subjectFunc SubjectFunc) *nidhi.MetadataProvider {
	return &nidhi.MetadataProvider{
		Keys: []string{CreatedKey, UpdatedKey, DeletedKey},
		Wrap: func(col nidhi.Interface) nidhi.Interface {
			return &provider{
				SubjectFunc: subjectFunc,
				Interface:   col,
			}
		},
	}
}

type SubjectFunc func(context.Context) string

// Metadata contains basic information by who and when records where added/updated/deleted.
type Metadata struct {
	Created, Updated, Deleted *ActivityLog
}

func (doc *Metadata) MarshalMetadata(w *jsoniter.Stream) error {
	if doc.Created != nil {
		w.WriteObjectField(CreatedKey)
		_ = doc.Created.MarshalDocument(w)
	}

	if doc.Updated != nil {
		if doc.Created != nil {
			w.WriteMore()
		}
		w.WriteObjectField(UpdatedKey)
		_ = doc.Updated.MarshalDocument(w)
	}

	if doc.Deleted != nil {
		if doc.Created != nil || doc.Updated != nil {
			w.WriteMore()
		}
		w.WriteObjectField(DeletedKey)
		_ = doc.Deleted.MarshalDocument(w)
	}

	return w.Error
}

func (doc *Metadata) UnmarshalMetadata(key string, r *jsoniter.Iterator) (bool, error) {
	var ac ActivityLog
	switch key {
	case CreatedKey:
		doc.Created = &ac
	case UpdatedKey:
		doc.Updated = &ac
	case DeletedKey:
		doc.Deleted = &ac
	default:
		return false, nil
	}

	if err := ac.UnmarshalDocument(r); err != nil {
		return true, err
	}

	return true, r.Error
}

type Creator struct {
	Values []*Metadata
}

func (c *Creator) Create() nidhi.MetadataUnmarshaler {
	var md Metadata
	c.Values = append(c.Values, &md)
	return &md
}

type ActivityLog struct {
	On time.Time `json:"on"`
	By string    `json:"by"`
}

func (log *ActivityLog) MarshalDocument(w *jsoniter.Stream) error {
	if log == nil {
		w.WriteNil()
		return w.Error
	}

	w.WriteObjectStart()

	w.WriteObjectField("by")
	w.WriteString(log.By)

	w.WriteMore()

	w.WriteObjectField("on")
	w.WriteString(log.On.Format(time.RFC3339Nano))

	w.WriteObjectEnd()

	return w.Error
}

func (log *ActivityLog) UnmarshalDocument(r *jsoniter.Iterator) error {
	if log == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "by":
			log.By = r.ReadString()
		case "on":
			log.On, _ = time.Parse(time.RFC3339, r.ReadString())
		default:
			r.Skip()
		}

		return true
	})

	return r.Error
}

type provider struct {
	SubjectFunc SubjectFunc
	nidhi.Interface
}

func (p *provider) Create(ctx context.Context, doc nidhi.Document, ops []nidhi.CreateOption) (string, error) {
	ops = append(ops, nidhi.WithCreateCreateMetadata(&Metadata{Created: p.activityLog(ctx)}), nidhi.WithCreateReplaceMetadata(&Metadata{Updated: p.activityLog(ctx)}))
	return p.Interface.Create(ctx, doc, ops)
}

func (p *provider) Replace(ctx context.Context, doc nidhi.Document, ops []nidhi.ReplaceOption) error {
	ops = append(ops, nidhi.WithReplaceMetadata(&Metadata{Updated: p.activityLog(ctx)}))
	return p.Interface.Replace(ctx, doc, ops)
}

func (p *provider) Update(ctx context.Context, doc nidhi.Document, f nidhi.Sqlizer, ops []nidhi.UpdateOption) error {
	ops = append(ops, nidhi.WithUpdateMetadata(&Metadata{Updated: p.activityLog(ctx)}))
	return p.Interface.Update(ctx, doc, f, ops)
}

func (p *provider) Delete(ctx context.Context, id string, ops []nidhi.DeleteOption) error {
	ops = append(ops, nidhi.WithDeleteMetadata(&Metadata{Deleted: p.activityLog(ctx)}))
	return p.Interface.Delete(ctx, id, ops)
}

func (p *provider) DeleteMany(ctx context.Context, f nidhi.Sqlizer, ops []nidhi.DeleteOption) error {
	ops = append(ops, nidhi.WithDeleteMetadata(&Metadata{Deleted: p.activityLog(ctx)}))
	return p.Interface.DeleteMany(ctx, f, ops)
}

func (p *provider) activityLog(ctx context.Context) *ActivityLog {
	return &ActivityLog{On: time.Now(), By: p.SubjectFunc(ctx)}
}
