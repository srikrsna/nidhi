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
	testCond(t, "empty", "field", &nidhi.StringCond{}, ``, []any{}, nidhi.ErrInvalidCond)
}

func TestOrderedCond(t *testing.T) {
	t.Parallel()
	testCond(t, "eq", "field", &nidhi.OrderedCond[int]{Eq: nidhi.Ptr(123)}, `field = ?`, []any{123}, nil)
	testCond(t, "gte", "field", &nidhi.OrderedCond[int]{Gte: nidhi.Ptr(123)}, `field >= ?`, []any{123}, nil)
	testCond(t, "lte", "field", &nidhi.OrderedCond[int]{Lte: nidhi.Ptr(123)}, `field <= ?`, []any{123}, nil)
	testCond(t, "gt", "field", &nidhi.OrderedCond[int]{Gt: nidhi.Ptr(123)}, `field > ?`, []any{123}, nil)
	testCond(t, "lt", "field", &nidhi.OrderedCond[int]{Lt: nidhi.Ptr(123)}, `field < ?`, []any{123}, nil)
	testCond(t, "multi", "field", &nidhi.OrderedCond[int]{Lt: nidhi.Ptr(123), Gt: nidhi.Ptr(234)}, `field < ? AND field > ?`, []any{123, 234}, nil)
	testCond(t, "empty", "field", &nidhi.OrderedCond[int]{}, ``, []any{}, nidhi.ErrInvalidCond)
}

func TestBoolCond(t *testing.T) {
	t.Parallel()
	testCond(t, "eq", "field", &nidhi.BoolCond{Eq: nidhi.Ptr(true)}, `field = ?`, []any{true}, nil)
	testCond(t, "empty", "field", &nidhi.BoolCond{}, ``, []any{}, nidhi.ErrInvalidCond)
}

func testCond(t *testing.T, name, field string, cond nidhi.Cond, wantQuery string, wantArgs []any, wantErr error) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		var sb strings.Builder
		gotArgs := []any{}
		gotErr := cond.AppendCond(field, &sb, &gotArgs)
		gotQuery := sb.String()
		if wantErr == nil {
			attest.Ok(t, gotErr)
		} else {
			attest.ErrorIs(t, gotErr, wantErr)
		}
		attest.Equal(t, gotQuery, wantQuery)
		attest.Equal(t, gotArgs, wantArgs, attest.Allow(pg.Array(nil)))
	})
}
