package nidhi_test

import (
	"errors"
	"testing"

	sq "github.com/elgris/sqrl"
	jsoniter "github.com/json-iterator/go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/srikrsna/nidhi"
)

func TestNidhi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nidhi Suite")
}

type byId []*testDoc

// Len is the number of elements in the collection.
func (a byId) Len() int           { return len(a) }
func (a byId) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a byId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type testDoc struct {
	Id     string `json:"Id,omitempty"`
	Number int    `json:"Number,omitempty"`
}

func (doc *testDoc) DocumentId() string {
	return doc.Id
}

func (doc *testDoc) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}
	w.WriteObjectStart()
	first := true
	if doc.Id != "" {
		if !first {
			w.WriteMore()
		}
		w.WriteObjectField("Id")
		w.WriteString(doc.Id)
		first = false
	}
	if doc.Number != 0 {
		if !first {
			w.WriteMore()
		}
		w.WriteObjectField("Number")
		w.WriteInt(doc.Number)
		first = false
	}
	w.WriteObjectEnd()
	return w.Error
}

func (doc *testDoc) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("nil document passed for unmarshal")
	}
	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "Id":
			doc.Id = r.ReadString()
		case "Number":
			doc.Number = r.ReadInt()
		default:
			r.Skip()
		}
		return true
	})
	return r.Error
}

type testFilter struct {
	Or     bool
	Id     *nidhi.StringFilter
	Number *nidhi.IntFilter
}

func (f *testFilter) ToSql(prefix string) (sq.Sqlizer, error) {
	of := &nidhi.ObjectFilter{}
	of.Or = f.Or
	var fp, op string
	if prefix != "" {
		fp = prefix + "->>"
		op = prefix + "->"
		_, _ = fp, op
	}

	of.Filter = map[string]nidhi.Filter{
		fp + "'Id'":     f.Id,
		fp + "'Number'": f.Number,
	}

	return of.ToSql("book")
}
