package nidhi_test

import (
	"context"
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/akshayjshah/attest"
	"github.com/elgris/sqrl"
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

func TestDeleteMany(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{})
	baseResource := &resource{
		Title:       "Resource",
		DateOfBirth: time.Now().UTC(),
		Age:         12,
		CanDrive:    true,
	}
	t.Run("many-soft", func(t *testing.T) {
		t.Parallel()
		const age = 100
		var rr []*resource
		for i := 0; i < 10; i++ {
			r := nidhi.Ptr(*baseResource)
			r.Id = "del-many-soft-" + strconv.Itoa(i)
			r.Age = age
			storeDoc(t, db, r, nil)
			rr = append(rr, r)
		}
		er := nidhi.Ptr(*baseResource)
		er.Id = "del-many-soft-e"
		er.Age = age + 1
		storeDoc(t, db, er, nil)
		res, err := store.DeleteMany(context.Background(), sqrl.Expr(`JSON_VALUE(`+nidhi.ColDoc+`, '$.age' RETURNING INT`+`) = ?`, age), nidhi.DeleteManyOptions{})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, res.DeleteCount, int64(len(rr)))
		for _, r := range rr {
			doc := getDoc(t, db, r.Id, nil, nil)
			attest.True(t, doc.Deleted)
		}
		doc := getDoc(t, db, er.Id, nil, nil)
		attest.False(t, doc.Deleted)
		attest.Equal(t, doc.Value, er)
	})
	t.Run("many-hard", func(t *testing.T) {
		t.Parallel()
		const age = 200
		var rr []*resource
		for i := 0; i < 10; i++ {
			r := nidhi.Ptr(*baseResource)
			r.Id = "del-many-hard-" + strconv.Itoa(i)
			r.Age = age
			storeDoc(t, db, r, nil)
			rr = append(rr, r)
		}
		er := nidhi.Ptr(*baseResource)
		er.Id = "del-many-hard-e"
		er.Age = age + 1
		storeDoc(t, db, er, nil)
		res, err := store.DeleteMany(context.Background(), sqrl.Expr(`JSON_VALUE(`+nidhi.ColDoc+`, '$.age' RETURNING INT`+`) = ?`, age), nidhi.DeleteManyOptions{
			Permanent: true,
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, res.DeleteCount, int64(len(rr)))
		for _, r := range rr {
			getDoc(t, db, r.Id, nil, sql.ErrNoRows)
		}
		doc := getDoc(t, db, er.Id, nil, nil)
		attest.False(t, doc.Deleted)
		attest.Equal(t, doc.Value, er)
	})
	t.Run("many-soft-md", func(t *testing.T) {
		t.Parallel()
		const age = 300
		var rr []*resource
		for i := 0; i < 10; i++ {
			r := nidhi.Ptr(*baseResource)
			r.Id = "del-many-soft-md-" + strconv.Itoa(i)
			r.Age = age
			md := nidhi.Metadata{"part": &metadataPart{Value: "value"}}
			storeDoc(t, db, r, md)
			rr = append(rr, r)
		}
		er := nidhi.Ptr(*baseResource)
		er.Id = "del-many-soft-md-e"
		er.Age = age + 1
		storeDoc(t, db, er, nil)
		md := nidhi.Metadata{"part": &metadataPart{Value: "deleted"}}
		res, err := store.DeleteMany(context.Background(), sqrl.Expr(`JSON_VALUE(`+nidhi.ColDoc+`, '$.age' RETURNING INT`+`) = ?`, age), nidhi.DeleteManyOptions{
			Metadata: md,
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, res.DeleteCount, int64(len(rr)))
		for _, r := range rr {
			doc := getDoc(t, db, r.Id, nidhi.Metadata{"part": &metadataPart{}}, nil)
			attest.True(t, doc.Deleted)
			attest.Equal(t, doc.Metadata, md)
		}
		doc := getDoc(t, db, er.Id, nil, nil)
		attest.False(t, doc.Deleted)
		attest.Equal(t, doc.Value, er)
	})
}
