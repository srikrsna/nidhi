package nidhi_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{})
	baseResource := &resource{
		Title:       "Resource",
		DateOfBirth: time.Now().UTC(),
		Age:         12,
		CanDrive:    true,
	}
	t.Run("single-soft", func(t *testing.T) {
		t.Parallel()
		r := nidhi.Ptr(*baseResource)
		r.Id = "del-single-soft"
		storeDoc(t, db, r, nil)
		res, err := store.Delete(context.Background(), r.Id, nidhi.DeleteOptions{})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		doc := getDoc(t, db, r.Id, nil, nil)
		attest.NotZero(t, doc)
		attest.True(t, doc.Deleted)
	})
	t.Run("single-hard", func(t *testing.T) {
		t.Parallel()
		r := nidhi.Ptr(*baseResource)
		r.Id = "del-single-hard"
		storeDoc(t, db, r, nil)
		res, err := store.Delete(context.Background(), r.Id, nidhi.DeleteOptions{Permanent: true})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		doc := getDoc(t, db, r.Id, nil, sql.ErrNoRows)
		attest.Zero(t, doc)
	})
	t.Run("single-soft-metadata", func(t *testing.T) {
		t.Parallel()
		r := nidhi.Ptr(*baseResource)
		r.Id = "del-single-soft-md"
		md := nidhi.Metadata{"part": &metadataPart{Value: "some"}}
		storeDoc(t, db, r, md)
		md = nidhi.Metadata{"part": &metadataPart{Value: "del"}}
		res, err := store.Delete(context.Background(), r.Id, nidhi.DeleteOptions{
			Metadata: md,
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		emd := nidhi.Metadata{"part": &metadataPart{}}
		doc := getDoc(t, db, r.Id, emd, nil)
		attest.NotZero(t, doc)
		attest.True(t, doc.Deleted)
		attest.Equal(t, emd, md)
	})
}
