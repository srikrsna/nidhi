package nidhi

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	sq "github.com/elgris/sqrl"
	"github.com/elgris/sqrl/pg"
)

type Sqlizer = sq.Sqlizer

var queryPool = sync.Pool{
	New: func() interface{} {
		return new(Query)
	},
}

func GetQuery() *Query {
	q := queryPool.Get().(*Query)
	q.Reset()
	return q
}

func PutQuery(q *Query) {
	queryPool.Put(q)
}

type Query struct {
	query strings.Builder
	args  []interface{}
	err   error
}

func (q *Query) Reset() {
	q.query.Reset()
	q.args = q.args[:0]
	q.err = nil
}

func (q *Query) Id(f *StringQuery) {
	if err := f.ToQuery(ColId, &q.query, &q.args); err != nil {
		q.err = err
	}
}

func (q *Query) Paren(iq Sqlizer) {
	query, args, err := iq.ToSql()
	if err != nil {
		q.err = err
	}

	q.query.WriteString(query)
	q.args = append(q.args, args...)
}

func (q *Query) Where(query string, args ...interface{}) {
	q.query.WriteString(query)
	q.args = append(q.args, args...)
}

func (q *Query) Not() {
	q.query.WriteString(" NOT")
}

func (q *Query) And() {
	q.query.WriteString(" AND")
}

func (q *Query) Or() {
	q.query.WriteString(" OR")
}

func (q *Query) Field(name string, f interface {
	ToQuery(string, io.StringWriter, *[]interface{}) error
}) {
	if err := f.ToQuery(name, &q.query, &q.args); err != nil {
		q.err = err
	}
}

func (q *Query) Prefix(p string) {
	q.query.WriteString(p)
}

func (q *Query) ReplaceArgs(args ...interface{}) error {
	if len(args) != len(q.args) {
		return fmt.Errorf("nidhi: different number of args are passed")
	}

	q.args = args

	return nil
}

func (q *Query) ToSql() (string, []interface{}, error) {
	return q.query.String(), q.args, q.err
}

type StringQuery struct {
	Like *string
	Eq   *string
}

func (s *StringQuery) ToQuery(name string, sb io.StringWriter, args *[]interface{}) error {
	if s == nil {
		return nil
	}

	if s.Eq != nil {
		sb.WriteString(name)
		sb.WriteString(" = ?")
		*args = append(*args, *s.Eq)
		return nil
	}
	if s.Like != nil {
		sb.WriteString(name)
		sb.WriteString(" like ?")
		*args = append(*args, *s.Like)
		return nil
	}

	return fmt.Errorf("nidhi: string filter %q is not set", name)
}

type BoolQuery struct {
	Eq *bool
}

func (s *BoolQuery) ToQuery(name string, sb io.StringWriter, args *[]interface{}) error {
	if s == nil {
		return nil
	}

	if s.Eq != nil {
		sb.WriteString(name)
		sb.WriteString(" = ?")
		*args = append(*args, *s.Eq)
		return nil
	}

	return fmt.Errorf("nidhi: string filter %q is not set", name)
}

type FloatQuery struct {
	Eq       *float64
	Lte, Gte *float64
	Lt, Gt   *float64
}

func (i *FloatQuery) ToQuery(name string, sb io.StringWriter, args *[]interface{}) error {
	if i == nil {
		return nil
	}

	first := true
	wc := func(sym string, value interface{}) {
		if !first {
			sb.WriteString(" AND")
		}
		first = false
		sb.WriteString(" ")
		sb.WriteString(name)
		sb.WriteString(" " + sym + " ?")
		*args = append(*args, value)
	}
	if i.Eq != nil {
		wc("=", *i.Eq)
	}

	if i.Lte != nil {
		wc("<=", *i.Lte)
	}

	if i.Gte != nil {
		wc(">=", *i.Gte)
	}

	if i.Lt != nil {
		wc("<", *i.Lt)
	}

	if i.Gt != nil {
		wc(">", *i.Gt)
	}

	if first {
		return fmt.Errorf("nidhi: int filter %q not set", name)
	}

	return nil
}

type IntQuery struct {
	Eq       *int64
	Lte, Gte *int64
	Lt, Gt   *int64
}

func (i *IntQuery) ToQuery(name string, sb io.StringWriter, args *[]interface{}) error {
	if i == nil {
		return nil
	}

	first := true
	wc := func(sym string, value interface{}) {
		if !first {
			sb.WriteString(" AND")
		}
		first = false
		sb.WriteString(" ")
		sb.WriteString(name)
		sb.WriteString(" " + sym + " ?")
		*args = append(*args, value)
	}
	if i.Eq != nil {
		wc("=", *i.Eq)
	}

	if i.Lte != nil {
		wc("<=", *i.Lte)
	}

	if i.Gte != nil {
		wc(">=", *i.Gte)
	}

	if i.Lt != nil {
		wc("<", *i.Lt)
	}

	if i.Gt != nil {
		wc(">", *i.Gt)
	}

	if first {
		return fmt.Errorf("nidhi: int filter %q not set", name)
	}

	return nil
}

type TimeQuery struct {
	Eq       *time.Time
	Lte, Gte *time.Time
	Lt, Gt   *time.Time
}

func (t *TimeQuery) ToQuery(name string, sb io.StringWriter, args *[]interface{}) error {
	if t == nil {
		return nil
	}

	first := true
	wc := func(sym string, value interface{}) {
		if !first {
			sb.WriteString(" AND")
		}
		first = false
		sb.WriteString(" ")
		sb.WriteString(name)
		sb.WriteString(" " + sym + " ?")
		*args = append(*args, value)
	}
	if t.Eq != nil {
		wc("=", *t.Eq)
	}

	if t.Lte != nil {
		wc("<=", *t.Lte)
	}

	if t.Gte != nil {
		wc(">=", *t.Gte)
	}

	if t.Lt != nil {
		wc("<", *t.Lt)
	}

	if t.Gt != nil {
		wc(">", *t.Gt)
	}

	if first {
		return fmt.Errorf("nidhi: int filter %q not set", name)
	}

	return nil
}

type SliceQuery struct {
	Slice interface{}
}

func (f *SliceQuery) ToQuery(name string, w io.StringWriter, args *[]interface{}) error {
	w.WriteString(name)
	w.WriteString(" @> ?")
	*args = append(*args, pg.Array(f.Slice))
	return nil
}

type MarshalerQuery struct {
	Marshaler
}

func (f MarshalerQuery) ToQuery(name string, w io.StringWriter, args *[]interface{}) error {
	w.WriteString(name)
	w.WriteString(" @> ?")
	*args = append(*args, JSONB(NoopUnmarshaler(f)))
	return nil
}

func String(s string) *string {
	return &s
}

func Int64(i int64) *int64 {
	return &i
}

func Int32(i int32) *int32 {
	return &i
}

func Int(i int) *int {
	return &i
}

func Float32(f float32) *float32 {
	return &f
}

func Float64(f float64) *float64 {
	return &f
}

func Bool(b bool) *bool {
	return &b
}

func Time(t time.Time) *time.Time {
	return &t
}
