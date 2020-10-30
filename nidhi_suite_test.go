package nidhi_test

import (
	"errors"
	"testing"

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

func (doc *testDoc) SetDocumentId(id string) {
	doc.Id = id
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

type testQuery nidhi.Query

func newTestQuery() *testQuery {
	return (*testQuery)(nidhi.GetQuery())
}

func (q *testQuery) query() *nidhi.Query {
	return (*nidhi.Query)(q)
}

func (q *testQuery) Id(f *nidhi.StringQuery) testConjuction {
	q.query().Id(f)
	return q
}

func (q *testQuery) Number(f *nidhi.IntQuery) testConjuction {
	q.query().Field("->'Number'", f)
	return q
}

func (q *testQuery) Where(query string, args ...interface{}) testConjuction {
	q.query().Where(query, args...)
	return q
}

func (q *testQuery) Not() *testQuery {
	q.query().Not()
	return q
}

func (q *testQuery) Paren(iq *testQuery) testConjuction {
	q.query().Paren(iq)
	return q
}

func (q *testQuery) And() *testQuery {
	q.query().And()
	return q
}

func (q *testQuery) Or() *testQuery {
	q.query().Or()
	return q
}

func (q *testQuery) ToSql() (string, []interface{}, error) {
	return q.query().ToSql()
}

type testConjuction interface {
	And() *testQuery
	Or() *testQuery

	nidhi.Sqlizer
}

func (q *testQuery) Reset() {
	q.query().Reset()
}
