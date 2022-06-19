package nidhi_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
)

func TestGet(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{MetadataRegistry: map[string]func() nidhi.MetadataPart{
		"part": func() nidhi.MetadataPart { return new(metadataPart) },
	}})
	t.Run("full", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "get-full"
		md := nidhi.Metadata{"part": &metadataPart{Value: "value"}}
		storeDoc(t, db, r, md)
		res, err := store.Get(context.Background(), r.Id, nidhi.GetOptions{})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, res.Metadata, md)
		attest.False(t, res.Deleted)
		attest.Equal(t, res.Value, r)
	})
	t.Run("partial", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "get-partial"
		md := nidhi.Metadata{"part": &metadataPart{Value: "value"}}
		storeDoc(t, db, r, md)
		res, err := store.Get(context.Background(), r.Id, nidhi.GetOptions{
			ViewMask: []string{"dateOfBirth"},
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, res.Metadata, md)
		attest.False(t, res.Deleted)
		attest.Equal(t, res.Value, &resource{DateOfBirth: r.DateOfBirth})
	})
	t.Run("missing", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "get-missing"
		res, err := store.Get(context.Background(), r.Id, nidhi.GetOptions{})
		attest.ErrorIs(t, err, nidhi.NotFound)
		attest.Zero(t, res)
	})
}

func TestQuery(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{MetadataRegistry: map[string]func() nidhi.MetadataPart{
		"part": func() nidhi.MetadataPart { return new(metadataPart) },
	}})
	t.Run("full", func(t *testing.T) {
		t.Parallel()
		const age = 100
		var rr []*resource
		for i := 0; i < 10; i++ {
			r := defaultResource()
			r.Id = "query-full-" + strconv.Itoa(i)
			r.Age = age
			storeDoc(t, db, r, nil)
			rr = append(rr, r)
		}
		res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, len(res.Docs), len(rr))
		for i, doc := range res.Docs {
			attest.Equal(t, doc.Value, rr[i])
		}
		attest.False(t, res.HasMore)
		attest.Zero(t, res.LastCursor)
	})
	t.Run("partial", func(t *testing.T) {
		t.Parallel()
		const age = 200
		var rr []*resource
		for i := 0; i < 10; i++ {
			r := defaultResource()
			r.Id = "query-partial-" + strconv.Itoa(i)
			r.Age = age
			storeDoc(t, db, r, nil)
			rr = append(rr, r)
		}
		res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
			ViewMask: []string{"dateOfBirth"},
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, len(res.Docs), len(rr))
		for i, doc := range res.Docs {
			attest.Equal(t, doc.Value, &resource{DateOfBirth: rr[i].DateOfBirth})
		}
		attest.False(t, res.HasMore)
		attest.Zero(t, res.LastCursor)
	})
}
