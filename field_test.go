package nidhi_test

import (
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/elgris/sqrl/pg"
	"github.com/srikrsna/nidhi"
)

func TestStringFilter(t *testing.T) {
	t.Parallel()
	f := nidhi.StringField("f")
	testFilter(t, "eq", f.Eq("_"), `f = ?`, []any{"_"})
	testFilter(t, "like", f.Like("_"), `f like ?`, []any{"_"})
	testFilter(t, "gte", f.Gte("_"), `f >= ?`, []any{"_"})
	testFilter(t, "lte", f.Lte("_"), `f <= ?`, []any{"_"})
	testFilter(t, "gt", f.Gt("_"), `f > ?`, []any{"_"})
	testFilter(t, "lt", f.Lt("_"), `f < ?`, []any{"_"})
	_, arrArg, _ := pg.Array([]string{"id"}).ToSql()
	testFilter(t, "in", f.In([]string{"id"}), `f = Any(?)`, arrArg)
}

func TestOrderedCond(t *testing.T) {
	t.Parallel()
	f := nidhi.IntField("f")
	const v = int64(123)
	testFilter(t, "eq", f.Eq(v), `f = ?`, []any{v})
	testFilter(t, "gte", f.Gte(v), `f >= ?`, []any{v})
	testFilter(t, "lte", f.Lte(v), `f <= ?`, []any{v})
	testFilter(t, "gt", f.Gt(v), `f > ?`, []any{v})
	testFilter(t, "lt", f.Lt(v), `f < ?`, []any{v})
	_, arrArg, _ := pg.Array([]int64{v}).ToSql()
	testFilter(t, "in", f.In([]int64{v}), `f = Any(?)`, arrArg)
}

func TestBoolFilter(t *testing.T) {
	t.Parallel()
	f := nidhi.BoolField("f")
	testFilter(t, "eq", f.Eq(true), `f = ?`, []any{true})
}

func TestSliceFilter(t *testing.T) {
	t.Parallel()
	arg := []int64{1}
	_, jsonArg, _ := pg.JSONB(arg).ToSql()
	f := nidhi.IntSliceField("f")
	testFilter(t, "eq", f.Eq(arg), `f = ?::jsonb`, jsonArg)
	testFilter(t, "gte", f.Gte(arg), `f >= ?::jsonb`, jsonArg)
	testFilter(t, "lte", f.Lte(arg), `f <= ?::jsonb`, jsonArg)
	testFilter(t, "gt", f.Gt(arg), `f > ?::jsonb`, jsonArg)
	testFilter(t, "lt", f.Lt(arg), `f < ?::jsonb`, jsonArg)
	testFilter(t, "ct", f.Contains(arg), `f @> ?::jsonb`, jsonArg)
	testFilter(t, "ctb", f.Within(arg), `f <@ ?::jsonb`, jsonArg)
	testFilter(t, "ovl", f.Overlaps(arg), `f && ?::jsonb`, jsonArg)
}

func testFilter(t *testing.T, name string, filter nidhi.Sqlizer, wantQuery string, wantArgs []any) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		gotQuery, gotArgs, gotErr := filter.ToSql()
		attest.Ok(t, gotErr)
		attest.Equal(t, gotQuery, wantQuery)
		attest.Equal(t, gotArgs, wantArgs, attest.Allow(pg.Array(nil), pg.JSONB(nil)))
	})
}
