package nidhi

import (
	"fmt"
	"time"

	sq "github.com/elgris/sqrl"
	"github.com/elgris/sqrl/pg"
	jsoniter "github.com/json-iterator/go"
)

type Sqlizer = sq.Sqlizer

type Filter interface {
	ToSql(prefix string) (Sqlizer, error)
}

type And conj

func (a And) ToSql(prefix string) (sq.Sqlizer, error) {
	return conj(a).toSql(false, prefix)
}

type Or conj

func (a Or) ToSql(prefix string) (sq.Sqlizer, error) {
	return conj(a).toSql(true, prefix)
}

type conj []Filter

func (conj conj) toSql(or bool, prefix string) (sq.Sqlizer, error) {
	if len(conj) == 0 {
		return nil, fmt.Errorf("nidhi: or filter %q empty", prefix)
	}

	if len(conj) == 1 {
		return conj[0].ToSql(prefix)
	}

	var c []sq.Sqlizer
	for _, f := range conj {
		sql, err := f.ToSql(prefix)
		if err != nil {
			return nil, err
		}
		c = append(c, sql)
	}

	if or {
		return sq.Or(c), nil
	}

	return sq.And(c), nil
}

type Not struct {
	Expr Filter
}

func (n Not) ToSql(prefix string) (sq.Sqlizer, error) {
	f, err := n.Expr.ToSql(prefix)
	if err != nil {
		return nil, err
	}

	return not{f}, nil
}

type not struct {
	expr sq.Sqlizer
}

func (n not) ToSql() (string, []interface{}, error) {
	sql, args, err := n.expr.ToSql()
	if err != nil {
		return "", nil, err
	}

	return " NOT (" + sql + ")", args, nil
}

type StringFilter struct {
	Like *string
	Eq   *string
}

func (s *StringFilter) ToSql(name string) (sq.Sqlizer, error) {
	if s == nil {
		return nil, nil
	}

	if s.Eq != nil {
		return sq.Eq{name: *s.Eq}, nil
	}
	if s.Like != nil {
		return sq.Expr(name+" like ?", *s.Like), nil
	}

	return nil, fmt.Errorf("nidhi: string filter %q is not set", name)
}

type BoolFilter struct {
	Eq *bool
}

func (s *BoolFilter) ToSql(name string) (sq.Sqlizer, error) {
	if s.Eq != nil {
		return sq.Eq{name: *s.Eq}, nil
	}

	return nil, fmt.Errorf("nidhi: string filter %q is not set", name)
}

type FloatFilter struct {
	Eq       *float64
	Lte, Gte *float64
	Lt, Gt   *float64
}

func (i *FloatFilter) ToSql(name string) (sq.Sqlizer, error) {
	var and sq.And
	if i.Eq != nil {
		and = append(and, sq.Eq{
			name: i.Eq,
		})
	}

	if i.Lte != nil {
		and = append(and, sq.LtOrEq{name: i.Lte})
	}

	if i.Gte != nil {
		and = append(and, sq.LtOrEq{name: i.Gte})
	}

	if i.Lt != nil {
		and = append(and, sq.Lt{name: i.Lt})
	}

	if i.Gt != nil {
		and = append(and, sq.Gt{name: i.Gt})
	}

	if len(and) == 0 {
		return nil, fmt.Errorf("nidhi: float filter %q not set", name)
	}

	if len(and) == 1 {
		return and[0], nil
	}

	return and[:2], nil
}

type IntFilter struct {
	Eq       *int64
	Lte, Gte *int64
	Lt, Gt   *int64
}

func (i *IntFilter) ToSql(name string) (sq.Sqlizer, error) {
	var and sq.And
	if i.Eq != nil {
		and = append(and, sq.Eq{
			name: i.Eq,
		})
	}

	if i.Lte != nil {
		and = append(and, sq.LtOrEq{name: i.Lte})
	}

	if i.Gte != nil {
		and = append(and, sq.LtOrEq{name: i.Gte})
	}

	if i.Lt != nil {
		and = append(and, sq.Lt{name: i.Lt})
	}

	if i.Gt != nil {
		and = append(and, sq.Gt{name: i.Gt})
	}

	if len(and) == 0 {
		return nil, fmt.Errorf("nidhi: int filter %q not set", name)
	}

	if len(and) == 1 {
		return and[0], nil
	}

	return and[:2], nil
}

type TimeFilter struct {
	Eq       *int64
	Lte, Gte *int64
	Lt, Gt   *int64
}

func (t *TimeFilter) ToSql(name string) (sq.Sqlizer, error) {
	var and sq.And
	if t.Eq != nil {
		and = append(and, sq.Eq{name: time.Unix(*t.Eq, 0)})
	}

	if t.Lt != nil {
		and = append(and, sq.Lt{name: time.Unix(*t.Lt, 0)})
	}

	if t.Gt != nil {
		and = append(and, sq.Gt{name: time.Unix(*t.Gt, 0)})
	}

	if t.Lte != nil {
		and = append(and, sq.LtOrEq{name: time.Unix(*t.Lte, 0)})
	}

	if t.Gte != nil {
		and = append(and, sq.LtOrEq{name: time.Unix(*t.Gte, 0)})
	}

	if len(and) == 0 {
		return nil, fmt.Errorf("nidhi: time filter %q not set", name)
	}

	if len(and) == 1 {
		return and[0], nil
	}

	return and, nil
}

type ObjectFilter struct {
	Or bool

	Filter map[string]Filter
}

func (o *ObjectFilter) ToSql(prefix string) (sq.Sqlizer, error) {
	conj := make([]sq.Sqlizer, 0, len(o.Filter))

	for k, f := range o.Filter {
		cond, err := f.ToSql(k)
		if err != nil {
			return nil, err
		}
		if cond != nil {
			conj = append(conj, cond)
		}
	}

	if len(conj) == 0 {
		return nil, fmt.Errorf("nidhi: object filter %q empty", prefix)
	}

	if len(conj) == 1 {
		return conj[0], nil
	}

	if o.Or {
		return sq.Or(conj), nil
	}

	return sq.And(conj), nil
}

type ArrayFilter []interface{}

func (f ArrayFilter) ToSql(prefix string) (sq.Sqlizer, error) {
	return sq.Expr(prefix+" @> ?", pg.Array(f)), nil
}

type ObjectArrayFilter []Marshaler

func (f ObjectArrayFilter) ToSql(prefix string) (sq.Sqlizer, error) {
	return sq.Expr(prefix+" @> ?", JSONB(f)), nil
}

func (f ObjectArrayFilter) MarshalDocument(w *jsoniter.Stream) error {
	w.WriteArrayStart()

	for _, e := range f {
		e.MarshalDocument(w)
	}

	w.WriteArrayEnd()

	return w.Error
}

func (f ObjectArrayFilter) UnmarshalDocument(_ *jsoniter.Iterator) error {
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
