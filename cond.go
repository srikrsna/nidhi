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
	// `field` is the selector on which the condition is applied on.
	AppendCond(field string, w io.StringWriter, args *[]interface{}) error
}

// Ptr is a helper function to use along with various conds.
// It returns a pointer to a value.
//
// Eg:
//	&StringCond{Eq: Ptr("nidhi")}
func Ptr[T any](v T) *T {
	return &v
}

type StringCond struct {
	Like *string
	Eq   *string
	In   []string
}

func (s *StringCond) AppendCond(name string, sb io.StringWriter, args *[]interface{}) error {
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
	return fmt.Errorf("nidhi: condition for %q is not set: %w", name, ErrInvalidCond)
}

type BoolCond struct {
	Eq *bool
}

func (s *BoolCond) AppendCond(name string, sb io.StringWriter, args *[]interface{}) error {
	if s.Eq != nil {
		sb.WriteString(name)
		sb.WriteString(" = ?")
		*args = append(*args, *s.Eq)
		return nil
	}
	return fmt.Errorf("nidhi: condition for %q is not set: %w", name, ErrInvalidCond)
}

type OrderedCond[T constraints.Ordered | time.Time] struct {
	Eq       *T
	Lte, Gte *T
	Lt, Gt   *T
}

func (c *OrderedCond[T]) AppendCond(name string, sb io.StringWriter, args *[]interface{}) error {
	first := true
	wc := func(op string, value interface{}) {
		if !first {
			sb.WriteString(" AND ")
		}
		first = false
		sb.WriteString(name)
		sb.WriteString(" " + op + " ?")
		*args = append(*args, value)
	}
	if c.Eq != nil {
		wc("=", *c.Eq)
	}
	if c.Lte != nil {
		wc("<=", *c.Lte)
	}
	if c.Gte != nil {
		wc(">=", *c.Gte)
	}
	if c.Lt != nil {
		wc("<", *c.Lt)
	}
	if c.Gt != nil {
		wc(">", *c.Gt)
	}
	if first {
		return fmt.Errorf("nidhi: condition for %q not set: %w", name, ErrInvalidCond)
	}
	return nil
}

type FloatCond = OrderedCond[float64]

type IntCond = OrderedCond[int64]

type TimeCond = OrderedCond[time.Time]

// type SliceQuery[S ~[]S, E any] struct {
// 	Slice   S
// 	Options SliceOptions
// }

// func (s *SliceQuery[S, E]) AppendCond(name string, sb io.StringWriter, args *S) error {
// 	first := true
// 	wc := func(sym string) {
// 		if !first {
// 			sb.WriteString(" AND")
// 		}
// 		first = false
// 		sb.WriteString(" ")
// 		sb.WriteString(name)
// 		sb.WriteString(" " + sym + " ?")
// 		*args = append(*args, Jsonb{s.Slice})
// 	}
// 	if s.Options.Eq != nil {
// 		wc("=")
// 	}

// 	if s.Options.Neq != nil {
// 		wc("<>")
// 	}

// 	if s.Options.Lte != nil {
// 		wc("<=")
// 	}

// 	if s.Options.Gte != nil {
// 		wc(">=")
// 	}

// 	if s.Options.Lt != nil {
// 		wc("<")
// 	}

// 	if s.Options.Gt != nil {
// 		wc(">")
// 	}

// 	if s.Options.Ct != nil {
// 		wc("@>")
// 	}

// 	if s.Options.Ctb != nil {
// 		wc("<@")
// 	}

// 	if s.Options.Ovl != nil {
// 		wc("&&")
// 	}

// 	if first {
// 		return fmt.Errorf("nidhi: int filter %q not set", name)
// 	}

// 	return nil
// }

// type MarshalerQuery struct {
// 	Any any
// }

// func (f MarshalerQuery) AppendCond(name string, w io.StringWriter, args *[]any) error {
// 	w.WriteString(name)
// 	w.WriteString(" @> ?")
// 	*args = append(*args, JSONB(f.Any))
// 	return nil
// }

// type SliceOptions struct {
// 	// Equal (=), Not Equal (<>), Less than (<), Greater Than (>), Less than Equal (<=), Greater Than Equal (>=), Contains (@>), Contained By (<@), Overlap (&&)
// 	Eq, Neq, Lt, Gt, Lte, Gte, Ct, Ctb, Ovl *struct{}
// }
