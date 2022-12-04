package nidhi

import (
	"time"

	sq "github.com/elgris/sqrl"
	"github.com/elgris/sqrl/pg"
)

// Sqlizer is the type expected by a store to query documents.
//
// ToSql returns a SQL representation of the Sqlizer, along with a slice of args
// as passed to e.g. database/sql.Exec. It can also return an error.
type Sqlizer = sq.Sqlizer

type (
	// IntField represents a field of type [int].
	IntField = orderedField[int64]
	// FloatField represents a field of type [float64].
	FloatField = orderedField[float64]
	// TimeField represents a field of type [time.Time].
	TimeField = orderedField[time.Time]
	// BoolField represents a field of type [bool].
	BoolField string
	// StringField represents a field of type [string].
	StringField string
	// StringSliceField represents a field of type [[]string].
	StringSliceField = sliceField[string, []string]
	// IntSliceField represents a field of type [[]int].
	IntSliceField = sliceField[int64, []int64]
	// BoolSliceField is convenient alias
	BoolSliceField = sliceField[bool, []bool]
	// FloatSliceField is convenient alias
	FloatSliceField = sliceField[float64, []float64]
	// TimeSliceField is convenient alias
	TimeSliceField = sliceField[time.Time, []time.Time]
	// JsonField represents a field of type json.
	JsonField string
)

// And applies 'AND' between the conditions.
func And(first, second Sqlizer, remaining ...Sqlizer) Sqlizer {
	return sq.And(append([]sq.Sqlizer{first, second}, remaining...))
}

// Or applies 'OR' between the conditions.
func Or(first, second Sqlizer, remaining ...Sqlizer) Sqlizer {
	return sq.Or(append([]sq.Sqlizer{first, second}, remaining...))
}

// Not negates the condition.
func Not(this Sqlizer) Sqlizer {
	return sq.Expr("NOT (?)", this)
}

type (
	orderedField[T int64 | float64 | time.Time] string
	sliceField[E any, S ~[]E]                   string
)

func (f orderedField[T]) Eq(v T) Sqlizer   { return sq.Expr(string(f)+" = ?", v) }
func (f orderedField[T]) Gt(v T) Sqlizer   { return sq.Expr(string(f)+" > ?", v) }
func (f orderedField[T]) Gte(v T) Sqlizer  { return sq.Expr(string(f)+" >= ?", v) }
func (f orderedField[T]) Lt(v T) Sqlizer   { return sq.Expr(string(f)+" < ?", v) }
func (f orderedField[T]) Lte(v T) Sqlizer  { return sq.Expr(string(f)+" <= ?", v) }
func (f orderedField[T]) In(v []T) Sqlizer { return sq.Expr(string(f)+" = Any(?)", pg.Array(v)) }

func (f StringField) Like(v string) Sqlizer { return sq.Expr(string(f)+" like ?", v) }
func (f StringField) Eq(v string) Sqlizer   { return sq.Expr(string(f)+" = ?", v) }
func (f StringField) Gt(v string) Sqlizer   { return sq.Expr(string(f)+" > ?", v) }
func (f StringField) Gte(v string) Sqlizer  { return sq.Expr(string(f)+" >= ?", v) }
func (f StringField) Lt(v string) Sqlizer   { return sq.Expr(string(f)+" < ?", v) }
func (f StringField) Lte(v string) Sqlizer  { return sq.Expr(string(f)+" <= ?", v) }
func (f StringField) In(v []string) Sqlizer { return sq.Expr(string(f)+" = Any(?)", pg.Array(v)) }

func (f BoolField) Eq(v bool) Sqlizer { return sq.Expr(string(f)+" = ?", v) }

func (f sliceField[E, S]) Eq(v S) Sqlizer       { return sq.Expr(string(f)+" = ?", pg.JSONB(v)) }
func (f sliceField[E, S]) Lt(v S) Sqlizer       { return sq.Expr(string(f)+" < ?", pg.JSONB(v)) }
func (f sliceField[E, S]) Gt(v S) Sqlizer       { return sq.Expr(string(f)+" > ?", pg.JSONB(v)) }
func (f sliceField[E, S]) Lte(v S) Sqlizer      { return sq.Expr(string(f)+" <= ?", pg.JSONB(v)) }
func (f sliceField[E, S]) Gte(v S) Sqlizer      { return sq.Expr(string(f)+" >= ?", pg.JSONB(v)) }
func (f sliceField[E, S]) Contains(v S) Sqlizer { return sq.Expr(string(f)+" @> ?", pg.JSONB(v)) }
func (f sliceField[E, S]) Within(v S) Sqlizer   { return sq.Expr(string(f)+" <@ ?", pg.JSONB(v)) }
func (f sliceField[E, S]) Overlaps(v S) Sqlizer { return sq.Expr(string(f)+" && ?", pg.JSONB(v)) }

func (f JsonField) Contains(v any) Sqlizer {
	buf, err := getJson(v)
	if err != nil {
		return errSqlizer{err}
	}
	return sq.Expr(string(f)+" @> ?", buf.Buffer())
}

type errSqlizer struct{ error }

func (e errSqlizer) ToSql() (string, []interface{}, error) {
	return "", nil, error(e)
}
