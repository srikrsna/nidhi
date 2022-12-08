package nidhi

import (
	"strings"
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
	// BoolField represents a field of type [bool].
	BoolField string
	// StringField represents a field of type [string].
	StringField string
	// DocField represents a field of type document.
	DocField[T any] string
	// IntField represents a field of type [int].
	IntField = orderedField[int64]
	// FloatField represents a field of type [float64].
	FloatField = orderedField[float64]
	// TimeField represents a field of type [time.Time].
	TimeField = orderedField[time.Time]
	// ListField represents a field of type list.
	ListField[E any, S ~[]E, F any] struct {
		// The SQL accessor for the field.
		//
		// Eg: (document#>{path})::jsonb
		Accessor string
		// The path to the field.
		Path []string
		// At is the index of the element in the slice to match.
		At func(int) F
	}
)

// NewIntField returns a new IntField for the given path.
func NewIntField(path []string) IntField {
	return IntField(newPrimitiveField(path, "#>", "BIGINT"))
}

func NewFloatField(path []string) FloatField {
	return FloatField(newPrimitiveField(path, "#>", "DOUBLE PRECISION"))
}

// NewBoolField returns a new BoolField for the given path.
func NewBoolField(path []string) BoolField {
	return BoolField(newPrimitiveField(path, "#>", "BOOLEAN"))
}

// NewTimeField returns a new TimeField for the given path.
func NewTimeField(path []string) TimeField {
	return TimeField(newPrimitiveField(path, "#>>", "TIMESTAMP"))
}

// NewFloatField returns a new FloatField for the given path.
func NewStringField(path []string) StringField {
	return StringField(newPrimitiveField(path, "#>>", "TEXT"))
}

// NewDocField returns a new DocField for the given path.
func NewDocField[T any](path []string) DocField[T] {
	return DocField[T](newPrimitiveField(path, "#>", "JSONB"))
}

// NewListField returns a new ListField for the given path.
func NewListField[E any, S ~[]E, F any](path []string, at func(int) F) ListField[E, S, F] {
	return ListField[E, S, F]{
		Accessor: newPrimitiveField(path, "#>", "JSONB"),
		Path:     path,
		At:       at,
	}
}

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

func (f ListField[E, S, F]) Eq(v S) Sqlizer       { return sq.Expr(f.Accessor+" = ?", pg.JSONB(v)) }
func (f ListField[E, S, F]) Lt(v S) Sqlizer       { return sq.Expr(f.Accessor+" < ?", pg.JSONB(v)) }
func (f ListField[E, S, F]) Gt(v S) Sqlizer       { return sq.Expr(f.Accessor+" > ?", pg.JSONB(v)) }
func (f ListField[E, S, F]) Lte(v S) Sqlizer      { return sq.Expr(f.Accessor+" <= ?", pg.JSONB(v)) }
func (f ListField[E, S, F]) Gte(v S) Sqlizer      { return sq.Expr(f.Accessor+" >= ?", pg.JSONB(v)) }
func (f ListField[E, S, F]) Contains(v S) Sqlizer { return sq.Expr(f.Accessor+" @> ?", pg.JSONB(v)) }
func (f ListField[E, S, F]) Within(v S) Sqlizer   { return sq.Expr(f.Accessor+" <@ ?", pg.JSONB(v)) }
func (f ListField[E, S, F]) Overlaps(v S) Sqlizer { return sq.Expr(f.Accessor+" && ?", pg.JSONB(v)) }

func (f DocField[T]) Contains(v T) Sqlizer {
	buf, err := getJson(v)
	if err != nil {
		return errSqlizer{err}
	}
	return sq.Expr(string(f)+" @> ?", buf.Buffer())
}

func newPrimitiveField(path []string, op, typ string) string {
	return ("(" + ColDoc + op + "'{" + strings.Join(path, ",") + "}')::" + typ)
}

type errSqlizer struct{ error }

func (e errSqlizer) ToSql() (string, []interface{}, error) {
	return "", nil, error(e)
}
