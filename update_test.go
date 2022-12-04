package nidhi_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
)

func TestReplace(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{})
	r := defaultResource()
	r.Id = "replace"
	storeDoc(t, db, r, nil)
	t.Run("id", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "replace-id"
		storeDoc(t, db, r, nil)
		r.Title = "Replaced title"
		md := nidhi.Metadata{"part": &metadataPart{Value: "value"}}
		res, err := store.Replace(context.Background(), r, nidhi.ReplaceOptions{
			Metadata: md,
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		got := getDoc(t, db, r.Id, nidhi.Metadata{"part": &metadataPart{}}, nil)
		attest.NotZero(t, got)
		attest.Equal(t, got.Revision, 2)
		attest.Equal(t, got.Value, r)
		attest.Equal(t, got.Metadata, md)
		t.Run("revision-mismatch", func(t *testing.T) {
			t.Parallel()
			res, err := store.Replace(context.Background(), r, nidhi.ReplaceOptions{Revision: 1})
			attest.ErrorIs(t, err, nidhi.ErrNotFound)
			attest.Zero(t, res)
		})
		t.Run("revision-match", func(t *testing.T) {
			t.Parallel()
			r.Title = "Replaced revision title"
			res, err := store.Replace(context.Background(), r, nidhi.ReplaceOptions{Revision: 2})
			attest.Ok(t, err)
			attest.NotZero(t, res)
			got := getDoc(t, db, r.Id, nil, nil)
			attest.NotZero(t, got)
			attest.Equal(t, got.Revision, 3)
			attest.Equal(t, got.Value, r)
		})
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{})
	r := defaultResource()
	r.Id = "update"
	storeDoc(t, db, r, nil)
	t.Run("id", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "update-id"
		storeDoc(t, db, r, nil)
		updates := &resourceUpdates{
			Title: ptr(""),
		}
		r.Title = ""
		md := nidhi.Metadata{"part": &metadataPart{Value: "value"}}
		res, err := store.Update(context.Background(), r.Id, updates, nidhi.UpdateOptions{
			Metadata: md,
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		got := getDoc(t, db, r.Id, nidhi.Metadata{"part": &metadataPart{}}, nil)
		attest.NotZero(t, got)
		attest.Equal(t, got.Revision, 2)
		attest.Equal(t, got.Value, r)
		attest.Equal(t, got.Metadata, md)
		t.Run("revision-mismatch", func(t *testing.T) {
			t.Parallel()
			res, err := store.Update(context.Background(), r.Id, r, nidhi.UpdateOptions{Revision: 1})
			attest.ErrorIs(t, err, nidhi.ErrNotFound)
			attest.Zero(t, res)
		})
		t.Run("revision-match", func(t *testing.T) {
			t.Parallel()
			updates := &resourceUpdates{
				Title: ptr("Updated"),
			}
			r.Title = "Updated"
			res, err := store.Update(context.Background(), r.Id, updates, nidhi.UpdateOptions{Revision: 2})
			attest.Ok(t, err)
			attest.NotZero(t, res)
			got := getDoc(t, db, r.Id, nil, nil)
			attest.NotZero(t, got)
			attest.Equal(t, got.Revision, 3)
			attest.Equal(t, got.Value, r)
		})
	})
}

func TestUpdateMany(t *testing.T) {
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
			r.Id = "update-many-full-" + strconv.Itoa(i)
			r.Age = age
			storeDoc(t, db, r, nil)
			r.Title = ""
			rr = append(rr, r)
		}
		updates := &resourceUpdates{
			Title: ptr(""),
		}
		md := nidhi.Metadata{"part": &metadataPart{Value: "value"}}
		res, err := store.UpdateMany(context.Background(), updates, filterByAge(age), nidhi.UpdateManyOptions{Metadata: md})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, res.UpdateCount, int64(len(rr)))
		for _, r := range rr {
			doc := getDoc(t, db, r.Id, nidhi.Metadata{"part": &metadataPart{}}, nil)
			attest.Equal(t, doc.Value, r)
			attest.Equal(t, doc.Metadata, md)
		}
	})
}
