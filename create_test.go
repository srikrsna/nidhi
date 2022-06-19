package nidhi_test

import (
	"context"
	"testing"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{})

	t.Run("new", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "create-new"
		res, err := store.Create(context.Background(), r, nidhi.CreateOptions{})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		ed := getDoc(t, db, r.Id, nil, nil)
		attest.Equal(t, ed.Revision, 1)
		attest.Equal(t, ed.Value, r)
	})
	t.Run("replace", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "create-replace"
		storeDoc(t, db, r, nil)
		r.Title = "Updated Title"
		res, err := store.Create(context.Background(), r, nidhi.CreateOptions{Replace: true})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		ed := getDoc(t, db, r.Id, nil, nil)
		attest.Equal(t, ed.Revision, 2)
		attest.Equal(t, ed.Value, r)
	})
	t.Run("new-metadata", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "create-new-meta"
		md := nidhi.Metadata{"part": &metadataPart{Value: "some"}}
		res, err := store.Create(context.Background(), r, nidhi.CreateOptions{
			CreateMetadata: md,
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		emd := nidhi.Metadata{"part": &metadataPart{}}
		ed := getDoc(t, db, r.Id, emd, nil)
		attest.Equal(t, ed.Revision, 1)
		attest.Equal(t, ed.Value, r)
		attest.Equal(t, emd, md)
	})
	t.Run("replace-meta", func(t *testing.T) {
		t.Parallel()
		r := defaultResource()
		r.Id = "create-replace-meta"
		md := nidhi.Metadata{"part": &metadataPart{Value: "some"}}
		storeDoc(t, db, r, md)
		r.Title = "Updated Title"
		md = nidhi.Metadata{"part": &metadataPart{Value: "updated"}}
		res, err := store.Create(context.Background(), r, nidhi.CreateOptions{Replace: true, ReplaceMetadata: md})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		emd := nidhi.Metadata{"part": &metadataPart{}}
		ed := getDoc(t, db, r.Id, emd, nil)
		attest.Equal(t, ed.Revision, 2)
		attest.Equal(t, ed.Value, r)
		attest.Equal(t, emd, md)
	})
}
