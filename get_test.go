package nidhi_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/akshayjshah/attest"
	"github.com/srikrsna/nidhi"
	"golang.org/x/exp/slices"
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
		attest.ErrorIs(t, err, nidhi.ErrNotFound)
		attest.Zero(t, res)
	})
}

func TestQuery(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{MetadataRegistry: map[string]func() nidhi.MetadataPart{
		"part": func() nidhi.MetadataPart { return new(metadataPart) },
	}})
	timeSort := func(l, r *resource) bool { return l.DateOfBirth.Before(r.DateOfBirth) }
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
	t.Run("order-by", func(t *testing.T) {
		t.Parallel()
		const age = 300
		var rr []*resource
		for i := 0; i < 10; i++ {
			r := defaultResource()
			r.Id = "query-order-by-" + strconv.Itoa(10-i)
			r.DateOfBirth = time.Now().Add(time.Duration(i) * time.Hour)
			r.Age = age
			storeDoc(t, db, r, nil)
			rr = append(rr, r)
		}
		res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
			OrderBy: []nidhi.OrderBy{orderByDateOfBirth()},
		})
		attest.Ok(t, err)
		attest.NotZero(t, res)
		attest.Equal(t, len(res.Docs), len(rr))
		slices.SortStableFunc(rr, timeSort)
		for i, doc := range res.Docs {
			attest.Equal(t, doc.Value, rr[i])
		}
		attest.False(t, res.HasMore)
		attest.Zero(t, res.LastCursor)
	})
	t.Run("pagination", func(t *testing.T) {
		t.Parallel()
		const age = 400
		var rr []*resource
		now := time.Now().UTC()
		date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		md := nidhi.Metadata{"part": &metadataPart{Value: "Paginated"}}
		for i := 0; i < 10; i++ {
			r := defaultResource()
			r.Id = "query-pagination-" + strconv.Itoa(i)
			r.DateOfBirth = date.Add(time.Duration(10-i) * time.Hour)
			r.Title = strconv.Itoa(10 - i)
			r.Age = age
			storeDoc(t, db, r, md)
			rr = append(rr, r)
		}
		dr := defaultResource()
		dr.Id = "query-pagination-del"
		dr.DateOfBirth = date
		dr.Title = "10"
		dr.Age = age
		storeDoc(t, db, dr, md)
		markDeleted(t, db, dr.Id)
		rrd := append(rr, dr)
		t.Run("forward-first-page", func(t *testing.T) {
			res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
				PaginationOptions: &nidhi.PaginationOptions{
					Limit: uint64(len(rr) / 2),
				},
			})
			attest.Ok(t, err)
			attest.True(t, res.HasMore)
			attest.NotZero(t, res.LastCursor)
			attest.Equal(t, len(res.Docs), len(rr)/2)
			for i, doc := range res.Docs {
				attest.Equal(t, doc.Value, rr[i])
			}
		})
		t.Run("forward-all-pages", func(t *testing.T) {
			var cursor string
			const limit = 2
			var docs []*nidhi.Document[resource]
			for {
				res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
					PaginationOptions: &nidhi.PaginationOptions{
						Limit:  limit,
						Cursor: cursor,
					},
				})
				attest.Ok(t, err)
				attest.Equal(t, len(res.Docs), limit)
				cursor = res.LastCursor
				docs = append(docs, res.Docs...)
				if !res.HasMore {
					break
				}
			}
			attest.Equal(t, len(docs), len(rr))
			for i, doc := range docs {
				attest.Equal(t, doc.Value, rr[i])
			}
		})
		t.Run("reverse-all-pages", func(t *testing.T) {
			var cursor string
			const limit = 2
			var docs []*nidhi.Document[resource]
			for {
				res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
					PaginationOptions: &nidhi.PaginationOptions{
						Limit:    limit,
						Cursor:   cursor,
						Backward: true,
					},
				})
				attest.Ok(t, err)
				attest.Equal(t, len(res.Docs), limit)
				cursor = res.LastCursor
				docs = append(docs, res.Docs...)
				if !res.HasMore {
					break
				}
			}
			attest.Equal(t, len(docs), len(rr))
			for i, doc := range docs {
				attest.Equal(t, doc.Value, rr[len(rr)-i-1])
			}
		})
		t.Run("order-by-all-pages", func(t *testing.T) {
			var cursor string
			const limit = 2
			var docs []*nidhi.Document[resource]
			for {
				res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
					PaginationOptions: &nidhi.PaginationOptions{
						Limit:  limit,
						Cursor: cursor,
					},
					OrderBy: []nidhi.OrderBy{orderByDateOfBirth()},
				})
				attest.Ok(t, err)
				attest.True(t, len(res.Docs) <= limit)
				cursor = res.LastCursor
				docs = append(docs, res.Docs...)
				if !res.HasMore {
					break
				}
			}
			attest.Equal(t, len(docs), len(rr))
			slices.SortStableFunc(rr, timeSort)
			for i, doc := range docs {
				attest.Equal(t, doc.Value, rr[i])
			}
		})
		t.Run("order-by-all-pages-load-metadata", func(t *testing.T) {
			var cursor string
			const limit = 2
			var docs []*nidhi.Document[resource]
			for {
				res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
					PaginationOptions: &nidhi.PaginationOptions{
						Limit:  limit,
						Cursor: cursor,
					},
					OrderBy:           []nidhi.OrderBy{orderByDateOfBirth()},
					LoadMetadataParts: []string{"part"},
				})
				attest.Ok(t, err)
				attest.True(t, len(res.Docs) <= limit)
				cursor = res.LastCursor
				docs = append(docs, res.Docs...)
				if !res.HasMore {
					break
				}
			}
			attest.Equal(t, len(docs), len(rr))
			slices.SortStableFunc(rr, timeSort)
			for i, doc := range docs {
				attest.Equal(t, doc.Value, rr[i])
				attest.Equal(t, doc.Metadata, md)
			}
		})
		t.Run("order-by-all-pages-load-metadata-include-deleted", func(t *testing.T) {
			var cursor string
			const limit = 2
			var docs []*nidhi.Document[resource]
			for {
				res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
					PaginationOptions: &nidhi.PaginationOptions{
						Limit:  limit,
						Cursor: cursor,
					},
					OrderBy:           []nidhi.OrderBy{orderByDateOfBirth()},
					LoadMetadataParts: []string{"part"},
					IncludeDeleted:    true,
				})
				attest.Ok(t, err)
				attest.True(t, len(res.Docs) <= limit)
				cursor = res.LastCursor
				docs = append(docs, res.Docs...)
				if !res.HasMore {
					break
				}
			}

			attest.Equal(t, len(docs), len(rrd))
			slices.SortStableFunc(rrd, timeSort)
			for i, doc := range docs {
				attest.Equal(t, doc.Value, rrd[i])
				attest.Equal(t, doc.Metadata, md)
			}
		})
		t.Run("order-by-all-pages-include-deleted", func(t *testing.T) {
			var cursor string
			const limit = 2
			var docs []*nidhi.Document[resource]
			for {
				res, err := store.Query(context.Background(), filterByAge(age), nidhi.QueryOptions{
					PaginationOptions: &nidhi.PaginationOptions{
						Limit:  limit,
						Cursor: cursor,
					},
					OrderBy:        []nidhi.OrderBy{orderByDateOfBirth()},
					IncludeDeleted: true,
				})
				attest.Ok(t, err)
				attest.True(t, len(res.Docs) <= limit)
				cursor = res.LastCursor
				docs = append(docs, res.Docs...)
				if !res.HasMore {
					break
				}
			}
			attest.Equal(t, len(docs), len(rrd))
			slices.SortStableFunc(rrd, timeSort)
			for i, doc := range docs {
				attest.Equal(t, doc.Value, rrd[i])
			}
		})
	})
}
