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

type Query struct {
	err  error
	buf  strings.Builder
	args []interface{}
}

func (q *Query) Reset() {
	q.buf.Reset()
	q.args = q.args[:0]
	q.err = nil
}

func (q *Query) Field(name string, f Cond) *Conj {
	if err := f.AppendCond(name, &q.buf, &q.args); err != nil {
		q.err = err
	}
	return (*Conj)(q)
}

func (q *Query) Paren(iq Sqlizer) *Query {
	query, args, err := iq.ToSql()
	if err != nil {
		q.err = err
	}
	q.buf.Grow(len(query) + 4)
	q.buf.WriteString(" (")
	q.buf.WriteString(query)
	q.buf.WriteString(") ")
	q.args = append(q.args, args...)
	return q
}

func (q *Query) Where(query string, args ...interface{}) *Conj {
	q.buf.WriteString(query)
	q.args = append(q.args, args...)
	return (*Conj)(q)
}

func (q *Query) Not() *Query {
	q.buf.WriteString(" NOT ")
	return q
}

type Conj Query

func (q *Conj) And() *Query {
	q.buf.WriteString(" AND")
	return (*Query)(q)
}

func (q *Conj) Or() *Query {
	q.buf.WriteString(" OR")
	return (*Query)(q)
}

func (q *Conj) ReplaceArgs(args []interface{}) (*Conj, error) {
	if len(args) != len(q.args) {
		return nil, fmt.Errorf("nidhi: different number of args are passed")
	}
	res := &Conj{
		err:  q.err,
		args: args,
	}
	res.buf.WriteString(q.buf.String())
	return res, nil
}

func (q *Conj) ToSql() (string, []interface{}, error) {
	return q.buf.String(), q.args, q.err
}


