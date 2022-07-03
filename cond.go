package nidhi

import (
	"fmt"
	"io"
	"time"

	"github.com/elgris/sqrl/pg"
	"golang.org/x/exp/constraints"
)

var (
	_ Cond = (*StringCond)(nil)
	_ Cond = (*IntCond)(nil)
	_ Cond = (*BoolCond)(nil)
	_ Cond = (*FloatCond)(nil)
	_ Cond = (*TimeCond)(nil)
)

// Cond represents a filtering condition like >, <, =.
type Cond interface {
	// AppendCond appends a condition to a query.
	// [field] is the selector on which the condition is applied on.
	AppendCond(field string, w io.StringWriter, args *[]any) error
}

// Ptr is a helper function to use along with various conds.
// It returns a pointer to a value.
//
// Eg:
//	&StringCond{Eq: Ptr("nidhi")}
func Ptr[T any](v T) *T {
	return &v
}

// StringCond is a [Cond] for [string] type.
// It supports equality, like and in operations.
type StringCond struct {
	Like *string
	Eq   *string
	In   []string
}

// AppendCond implements [Cond] for [StringCond]
func (s *StringCond) AppendCond(name string, sb io.StringWriter, args *[]any) error {
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
	if len(s.In) > 0 {
		sb.WriteString(name)
		sb.WriteString(" = Any(?::text[])")
		*args = append(*args, pg.Array(s.In))
		return nil
	}
	return fmt.Errorf("nidhi: condition %T for %q is not set: %w", s, name, ErrInvalidCond)
}

// BoolCond is a [Cond] for [bool] type.
type BoolCond struct {
	Eq *bool
}

// AppendCond implements [Cond] for [BoolCond]
func (c *BoolCond) AppendCond(name string, sb io.StringWriter, args *[]any) error {
	if c.Eq != nil {
		sb.WriteString(name)
		sb.WriteString(" = ?")
		*args = append(*args, *c.Eq)
		return nil
	}
	return fmt.Errorf("nidhi: condition %T for %q is not set: %w", c, name, ErrInvalidCond)
}

// OrderedCond is a [Cond] for ordered (=, <=, >=, <, >) types.
type OrderedCond[T constraints.Ordered | time.Time] struct {
	Eq       *T
	Lte, Gte *T
	Lt, Gt   *T
}

// AppendCond implements [Cond] for [OrderedCond]
func (c *OrderedCond[T]) AppendCond(name string, sb io.StringWriter, args *[]any) error {
	first := true
	wc := func(op string, v *T) {
		if v == nil {
			return
		}
		if !first {
			sb.WriteString(" AND ")
		}
		first = false
		sb.WriteString(name)
		sb.WriteString(" ")
		sb.WriteString(op)
		sb.WriteString(" ?")
		*args = append(*args, *v)
	}
	wc("=", c.Eq)
	wc("<=", c.Lte)
	wc(">=", c.Gte)
	wc("<", c.Lt)
	wc(">", c.Gt)
	if first {
		return fmt.Errorf("nidhi: condition %T for %q not set: %w", c, name, ErrInvalidCond)
	}
	return nil
}

// FloatCond is a [Cond] for [float64].
type FloatCond = OrderedCond[float64]

// IntCond is a [Cond] for [int64].
type IntCond = OrderedCond[int64]

// TimeCond is a [Cond] for [time.Time].
type TimeCond = OrderedCond[time.Time]

// TimeCond is a [Cond] for slice types.
type SliceCond[E any, S ~[]E] struct {
	// Equal (=), Not Equal (<>), Less than (<), Greater Than (>), Less than Equal (<=), Greater Than Equal (>=), Contains (@>), Contained By (<@), Overlap (&&)
	Eq, Neq, Lt, Gt, Lte, Gte, Ct, Ctb, Ovl S
}

func (s *SliceCond[E, S]) AppendCond(name string, sb io.StringWriter, args *[]any) error {
	first := true
	wc := func(op string, v S) {
		if v == nil {
			return
		}
		if !first {
			sb.WriteString(" AND ")
		}
		first = false
		sb.WriteString(name)
		sb.WriteString(" ")
		sb.WriteString(op)
		sb.WriteString(" ?")
		*args = append(*args, pg.JSONB(v))
	}
	wc("=", s.Eq)
	wc("<>", s.Neq)
	wc("<=", s.Lte)
	wc(">=", s.Gte)
	wc("<", s.Lt)
	wc(">", s.Gt)
	wc("@>", s.Ct)
	wc("<@", s.Ctb)
	wc("&&", s.Ovl)
	if first {
		return fmt.Errorf("nidhi: slice condition for %q not set: %w", name, ErrInvalidCond)
	}
	return nil
}

// type MarshalerQuery struct {
// 	Any any
// }

// func (f MarshalerQuery) AppendCond(name string, w io.StringWriter, args *[]any) error {
// 	w.WriteString(name)
// 	w.WriteString(" @> ?")
// 	*args = append(*args, JSONB(f.Any))
// 	return nil
// }
