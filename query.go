package nidhi

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	sq "github.com/elgris/sqrl"
)

type Sqlizer = sq.Sqlizer

type Queryer interface {
	ToQuery(string, io.StringWriter, *[]interface{}) error
}

var queryPool = sync.Pool{
	New: func() interface{} {
		return &Query{
			query: &strings.Builder{},
		}
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
	err   error
	query *strings.Builder
	args  []interface{}
}

func (q *Query) Reset() {
	q.query.Reset()
	q.args = q.args[:0]
	q.err = nil
}

func (q *Query) Id(f *StringQuery) {
	if err := f.ToQuery(ColId, q.query, &q.args); err != nil {
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

func (q *Query) WhereMetadata(f Queryer) {
	q.Field("("+ColMeta, f)
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

func (q *Query) Field(name string, f Queryer) {
	if err := f.ToQuery(name, q.query, &q.args); err != nil {
		q.err = err
	}
}

func (q *Query) Prefix(p string) {
	q.query.WriteString(p)
}

func (q *Query) ReplaceArgs(args []interface{}) (*Query, error) {
	if len(args) != len(q.args) {
		return nil, fmt.Errorf("nidhi: different number of args are passed")
	}

	query := strings.Builder{}
	query.WriteString(q.query.String())

	res := GetQuery()

	res.err = q.err
	res.query = &query
	res.args = args

	return res, nil
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

type DurationQuery struct {
	Eq       *time.Duration
	Lte, Gte *time.Duration
	Lt, Gt   *time.Duration
}

func (d *DurationQuery) ToQuery(name string, sb io.StringWriter, args *[]interface{}) error {
	if d == nil {
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
	if d.Eq != nil {
		wc("=", int64(d.Eq.Seconds()))
	}

	if d.Lte != nil {
		wc("<=", int64(d.Lte.Seconds()))
	}

	if d.Gte != nil {
		wc(">=", int64(d.Gte.Seconds()))
	}

	if d.Lt != nil {
		wc("<", int64(d.Lt.Seconds()))
	}

	if d.Gt != nil {
		wc(">", int64(d.Gt.Seconds()))
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
	Slice   interface{}
	Options SliceOptions
}

func (s *SliceQuery) ToQuery(name string, sb io.StringWriter, args *[]interface{}) error {
	first := true
	wc := func(sym string) {
		if !first {
			sb.WriteString(" AND")
		}
		first = false
		sb.WriteString(" ")
		sb.WriteString(name)
		sb.WriteString(" " + sym + " ?")
		*args = append(*args, Jsonb{s.Slice})
	}
	if s.Options.Eq != nil {
		wc("=")
	}

	if s.Options.Neq != nil {
		wc("<>")
	}

	if s.Options.Lte != nil {
		wc("<=")
	}

	if s.Options.Gte != nil {
		wc(">=")
	}

	if s.Options.Lt != nil {
		wc("<")
	}

	if s.Options.Gt != nil {
		wc(">")
	}

	if s.Options.Ct != nil {
		wc("@>")
	}

	if s.Options.Ctb != nil {
		wc("<@")
	}

	if s.Options.Ovl != nil {
		wc("&&")
	}

	if first {
		return fmt.Errorf("nidhi: int filter %q not set", name)
	}

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

type SliceOptions struct {
	// Equal (=), Not Equal (<>), Less than (<), Greater Than (>), Less than Equal (<=), Greater Than Equal (>=), Contains (@>), Contained By (<@), Overlap (&&)
	Eq, Neq, Lt, Gt, Lte, Gte, Ct, Ctb, Ovl *struct{}
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

var marker = &struct{}{}

func Struct() *struct{} {
	return marker
}
