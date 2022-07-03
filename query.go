package nidhi

import (
	"fmt"
	"strings"

	sq "github.com/elgris/sqrl"
)

// Sqlizer is the type expected by a store to query documents.
//
// ToSql returns a SQL representation of the Sqlizer, along with a slice of args
// as passed to e.g. database/sql.Exec. It can also return an error.
type Sqlizer = sq.Sqlizer

// Field represents a field of a document.
type Field interface {
	// Selector returns the selector expression for this field.
	// Eg: For a field named `title`
	//		JSON_VALUE('$.title' RETURNING JSONB DEFAULT '{}' ON EMPTY)
	Selector() string
}

type Query[F Field] struct {
	err  error
	buf  strings.Builder
	args []any
}

func (q *Query[F]) Reset() {
	q.buf.Reset()
	q.args = q.args[:0]
	q.err = nil
}

func (q *Query[F]) Where(f F, c Cond) *Conj[F] {
	if err := c.AppendCond(f.Selector(), &q.buf, &q.args); err != nil {
		q.err = err
	}
	return (*Conj[F])(q)
}

func (q *Query[F]) Paren(iq *Conj[F]) *Conj[F] {
	query, args, err := iq.ToSql()
	if err != nil {
		q.err = err
	}
	q.buf.Grow(len(query) + 2)
	q.buf.WriteString("(")
	q.buf.WriteString(query)
	q.buf.WriteString(")")
	q.args = append(q.args, args...)
	return (*Conj[F])(q)
}

func (q *Query[F]) Not() *Query[F] {
	q.buf.WriteString("NOT ")
	return q
}

type Conj[F Field] Query[F]

func (q *Conj[F]) And() *Query[F] {
	q.buf.WriteString(" AND ")
	return (*Query[F])(q)
}

func (q *Conj[F]) Or() *Query[F] {
	q.buf.WriteString(" OR ")
	return (*Query[F])(q)
}

func (q *Conj[F]) ReplaceArgs(args []any) (*Conj[F], error) {
	if len(args) != len(q.args) {
		return nil, fmt.Errorf("nidhi: different number of args are passed")
	}
	res := &Conj[F]{
		err:  q.err,
		args: args,
	}
	res.buf.WriteString(q.buf.String())
	return res, nil
}

func (q *Conj[F]) ToSql() (string, []any, error) {
	return q.buf.String(), q.args, q.err
}
