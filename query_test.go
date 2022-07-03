package nidhi_test

import (
	"io"
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
)

func TestQueryer(t *testing.T) {
	t.Parallel()
	t.Run("where", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query[testField]{}
		gotSql, gotArgs, err := q.Where(testField{}, eqCond).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("not", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query[testField]{}
		gotSql, gotArgs, err := q.Not().Where(testField{}, eqCond).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `NOT field = ?`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("paren", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query[testField]{}
		sq := &nidhi.Query[testField]{}
		gotSql, gotArgs, err := q.Paren(sq.Where(testField{}, eqCond)).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `(field = ?)`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("and", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query[testField]{}
		gotSql, gotArgs, err := q.Where(testField{}, ltCond).And().Where(testField{}, gtCond).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field < ? AND field > ?`)
		attest.Equal(t, gotArgs, []any{0, 0})
	})
	t.Run("or", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query[testField]{}
		gotSql, gotArgs, err := q.Where(testField{}, ltCond).Or().Where(testField{}, gtCond).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field < ? OR field > ?`)
		attest.Equal(t, gotArgs, []any{0, 0})
	})
	t.Run("reset", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query[testField]{}
		gotSql, gotArgs, err := q.Where(testField{}, eqCond).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
		q.Reset()
		gotSql, gotArgs, err = q.Where(testField{}, ltCond).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field < ?`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("replace", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query[testField]{}
		conj := q.Where(testField{}, eqCond)
		gotSql, gotArgs, err := conj.ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
		rConj, err := conj.ReplaceArgs([]any{1})
		attest.Ok(t, err)
		gotSql, gotArgs, err = conj.ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
		gotSql, gotArgs, err = rConj.ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{1})
	})
}

var (
	eqCond nidhi.Cond = opCond("=")
	gtCond nidhi.Cond = opCond(">")
	ltCond nidhi.Cond = opCond("<")
)

type opCond string

func (c opCond) AppendCond(field string, w io.StringWriter, args *[]any) error {
	w.WriteString(field + " " + string(c) + " ?")
	*args = append(*args, 0)
	return nil
}

type testField struct{}

func (testField) Selector() string { return "field" }
