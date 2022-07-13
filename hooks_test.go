package nidhi_test

import (
	"context"
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
)

func TestHooks(t *testing.T) {
	db := newDB(t)
	const (
		create int = iota
		get
		query
		replace
		update
		updateMany
		delete
		deleteMany
	)
	called := map[int]bool{}
	store := newStore(t, db, nidhi.StoreOptions{
		Hooks: []nidhi.Hooks{
			{
				OnCreate:     func(*nidhi.HookContext, any, *nidhi.CreateOptions) { called[create] = true },
				OnGet:        func(*nidhi.HookContext, string, *nidhi.GetOptions) { called[get] = true },
				OnQuery:      func(*nidhi.HookContext, nidhi.Sqlizer, *nidhi.QueryOptions) { called[query] = true },
				OnDelete:     func(*nidhi.HookContext, string, *nidhi.DeleteOptions) { called[delete] = true },
				OnDeleteMany: func(*nidhi.HookContext, nidhi.Sqlizer, *nidhi.DeleteManyOptions) { called[deleteMany] = true },
				OnReplace:    func(*nidhi.HookContext, any, *nidhi.ReplaceOptions) { called[replace] = true },
				OnUpdate:     func(*nidhi.HookContext, string, any, *nidhi.UpdateOptions) { called[update] = true },
				OnUpdateMany: func(*nidhi.HookContext, any, nidhi.Sqlizer, *nidhi.UpdateManyOptions) { called[updateMany] = true },
			},
		},
	})
	ctx := context.TODO()
	_, _ = store.Create(ctx, defaultResource(), nidhi.CreateOptions{})
	attest.True(t, called[create])
	_, _ = store.Get(ctx, "", nidhi.GetOptions{})
	attest.True(t, called[get])
	_, _ = store.Query(ctx, nil, nidhi.QueryOptions{})
	attest.True(t, called[query])
	_, _ = store.Delete(ctx, "", nidhi.DeleteOptions{})
	attest.True(t, called[delete])
	_, _ = store.DeleteMany(ctx, nil, nidhi.DeleteManyOptions{})
	attest.True(t, called[deleteMany])
	_, _ = store.Replace(ctx, defaultResource(), nidhi.ReplaceOptions{})
	attest.True(t, called[replace])
	_, _ = store.Update(ctx, "", nil, nidhi.UpdateOptions{})
	attest.True(t, called[update])
	_, _ = store.UpdateMany(ctx, nil, nil, nidhi.UpdateManyOptions{})
	attest.True(t, called[updateMany])
}

func TestHookContext(t *testing.T) {
	var (
		idFn    = func(t *resource) string { return t.Id }
		setIdFn = func(t *resource, s string) { t.Id = s }
	)
	db := newDB(t)
	store, err := nidhi.NewStore(
		context.Background(),
		db,
		"schema",
		"table",
		[]string{"field"},
		idFn,
		setIdFn,
		nidhi.StoreOptions{},
	)
	attest.Ok(t, err)
	attest.NotZero(t, store)
	ctx := nidhi.NewHookContext(context.Background(), store)
	r := defaultResource()
	r.Id = "id"
	attest.Equal(t, ctx.Id(r), r.Id)
	ctx.SetId(r, "new")
	attest.Equal(t, r.Id, "new")
}
