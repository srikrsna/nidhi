package nidhi_test

import (
	"strings"
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/elgris/sqrl/pg"
	"github.com/srikrsna/nidhi"
)

func TestPtr(t *testing.T) {
	t.Parallel()
	i := "12312"
	attest.Equal(t, nidhi.Ptr(i), &i)
}

func TestStringCond(t *testing.T) {
	t.Parallel()
	testCond(t, "eq", "field", &nidhi.StringCond{Eq: nidhi.Ptr("id")}, `field = ?`, []any{"id"}, nil)
	testCond(t, "like", "field", &nidhi.StringCond{Like: nidhi.Ptr("id")}, `field like ?`, []any{"id"}, nil)
	testCond(t, "in", "field", &nidhi.StringCond{In: []string{"id"}}, `field = Any(?::text[])`, []any{pg.Array([]string{"id"})}, nil)
	testCond(t, "empty", "field", &nidhi.StringCond{}, ``, nil, nidhi.ErrInvalidCond)
}

func TestOrderedCond(t *testing.T) {
	t.Parallel()
	testCond(t, "eq", "field", &nidhi.OrderedCond[int]{Eq: nidhi.Ptr(123)}, `field = ?`, []any{123}, nil)
	testCond(t, "gte", "field", &nidhi.OrderedCond[int]{Gte: nidhi.Ptr(123)}, `field >= ?`, []any{123}, nil)
	testCond(t, "lte", "field", &nidhi.OrderedCond[int]{Lte: nidhi.Ptr(123)}, `field <= ?`, []any{123}, nil)
	testCond(t, "gt", "field", &nidhi.OrderedCond[int]{Gt: nidhi.Ptr(123)}, `field > ?`, []any{123}, nil)
	testCond(t, "lt", "field", &nidhi.OrderedCond[int]{Lt: nidhi.Ptr(123)}, `field < ?`, []any{123}, nil)
	testCond(t, "multi", "field", &nidhi.OrderedCond[int]{Lt: nidhi.Ptr(123), Gt: nidhi.Ptr(234)}, `field < ? AND field > ?`, []any{123, 234}, nil)
	testCond(t, "empty", "field", &nidhi.OrderedCond[int]{}, ``, nil, nidhi.ErrInvalidCond)
}

func TestBoolCond(t *testing.T) {
	t.Parallel()
	testCond(t, "eq", "field", &nidhi.BoolCond{Eq: nidhi.Ptr(true)}, `field = ?`, []any{true}, nil)
	testCond(t, "empty", "field", &nidhi.BoolCond{}, ``, nil, nidhi.ErrInvalidCond)
}

func TestSliceCond(t *testing.T) {
	t.Parallel()
	arg := []int{1}
	jsonArg := pg.JSONB(arg)
	testCond(t, "eq", "field", &nidhi.SliceCond[int, []int]{Eq: arg}, `field = ?`, []any{jsonArg}, nil)
	testCond(t, "gte", "field", &nidhi.SliceCond[int, []int]{Gte: arg}, `field >= ?`, []any{jsonArg}, nil)
	testCond(t, "lte", "field", &nidhi.SliceCond[int, []int]{Lte: arg}, `field <= ?`, []any{jsonArg}, nil)
	testCond(t, "gt", "field", &nidhi.SliceCond[int, []int]{Gt: arg}, `field > ?`, []any{jsonArg}, nil)
	testCond(t, "lt", "field", &nidhi.SliceCond[int, []int]{Lt: arg}, `field < ?`, []any{jsonArg}, nil)
	testCond(t, "ct", "field", &nidhi.SliceCond[int, []int]{Ct: arg}, `field @> ?`, []any{jsonArg}, nil)
	testCond(t, "ctb", "field", &nidhi.SliceCond[int, []int]{Ctb: arg}, `field <@ ?`, []any{jsonArg}, nil)
	testCond(t, "ovl", "field", &nidhi.SliceCond[int, []int]{Ovl: arg}, `field && ?`, []any{jsonArg}, nil)
	testCond(t, "multi", "field", &nidhi.SliceCond[int, []int]{Lt: []int{1}, Gt: []int{2}}, `field < ? AND field > ?`, []any{pg.JSONB([]int{1}), pg.JSONB([]int{2})}, nil)
	testCond(t, "empty", "field", &nidhi.SliceCond[int, []int]{}, ``, nil, nidhi.ErrInvalidCond)
}

func testCond(t *testing.T, name, field string, cond nidhi.Cond, wantQuery string, wantArgs []any, wantErr error) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		var sb strings.Builder
		var gotArgs []any
		gotErr := cond.AppendCond(field, &sb, &gotArgs)
		gotQuery := sb.String()
		if wantErr == nil {
			attest.Ok(t, gotErr)
		} else {
			attest.ErrorIs(t, gotErr, wantErr)
		}
		attest.Equal(t, gotQuery, wantQuery)
		attest.Equal(t, gotArgs, wantArgs, attest.Allow(pg.Array(nil), pg.JSONB(nil)))
	})
}
