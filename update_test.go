package nidhi_test

import (
	"context"
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
