package nidhi_test

import (
	"io"
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
)

func TestQueryer(t *testing.T) {
	t.Parallel()
	t.Run("field", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		gotSql, gotArgs, err := q.Field("field", eqCond{}).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("not", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		gotSql, gotArgs, err := q.Not().Field("field", eqCond{}).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `NOT field = ?`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("paren", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		sq := &nidhi.Query{}
		gotSql, gotArgs, err := q.Paren(sq.Field("field", eqCond{})).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `(field = ?)`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("where", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		gotSql, gotArgs, err := q.Where("field = ?", 0).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("and", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		gotSql, gotArgs, err := q.Where("field < ?", 1).And().Where("field > ?", 2).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field < ? AND field > ?`)
		attest.Equal(t, gotArgs, []any{1, 2})
	})
	t.Run("or", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		gotSql, gotArgs, err := q.Where("field < ?", 1).Or().Where("field > ?", 2).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field < ? OR field > ?`)
		attest.Equal(t, gotArgs, []any{1, 2})
	})
	t.Run("reset", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		gotSql, gotArgs, err := q.Where("field = ?", 0).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
		q.Reset()
		gotSql, gotArgs, err = q.Where("field = ?", 0).ToSql()
		attest.Ok(t, err)
		attest.Equal(t, gotSql, `field = ?`)
		attest.Equal(t, gotArgs, []any{0})
	})
	t.Run("replace", func(t *testing.T) {
		t.Parallel()
		q := &nidhi.Query{}
		conj := q.Where("field = ?", 0)
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
	_ nidhi.Cond = eqCond{}
)

type eqCond struct{}

func (eqCond) AppendCond(field string, w io.StringWriter, args *[]any) error {
	w.WriteString(field + " = ?")
	*args = append(*args, 0)
	return nil
}
