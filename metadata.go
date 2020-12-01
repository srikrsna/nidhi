package nidhi

import (
	"context"
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Metadata struct {
	Created, Updated, Deleted *ActivityLog
}

type ActivityLog struct {
	By string    `json:"by"`
	On time.Time `json:"on"`
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

func (doc *Metadata) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	w.WriteObjectStart()

	if doc.Created != nil {
		w.WriteObjectField("created")
		_ = doc.Created.MarshalDocument(w)
	}

	if doc.Updated != nil {
		if doc.Created != nil {
			w.WriteMore()
		}
		w.WriteObjectField("updated")
		_ = doc.Updated.MarshalDocument(w)
	}

	if doc.Deleted != nil {
		if doc.Created != nil || doc.Updated != nil {
			w.WriteMore()
		}
		w.WriteObjectField("deleted")
		_ = doc.Deleted.MarshalDocument(w)
	}

	w.WriteObjectEnd()

	return w.Error
}

func (doc *Metadata) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "created":
			doc.Created = &ActivityLog{}
			_ = doc.Created.UnmarshalDocument(r)
		case "updated":
			doc.Updated = &ActivityLog{}
			_ = doc.Updated.UnmarshalDocument(r)
		case "deleted":
			doc.Deleted = &ActivityLog{}
			_ = doc.Deleted.UnmarshalDocument(r)
		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type SubjectFunc func(context.Context) string
