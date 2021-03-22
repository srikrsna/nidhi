package nidhi_test

import (
	"context"
	"database/sql"
	"math/rand"
	"sort"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gocloud.dev/postgres"

	"github.com/srikrsna/nidhi"
)

var _ = Describe("Collection", func() {
	var (
		db    *sql.DB
		col   *nidhi.Collection
		ctx   = context.TODO()
		count = 0
	)
	BeforeSuite(func() {
		var err error
		db, err = postgres.Open(ctx, "postgres://krsna@localhost/postgres?sslmode=disable")
		Expect(db, err).NotTo(BeNil())
		Expect(db.Ping()).To(Succeed())
		Expect(db.Exec(`TRUNCATE collection_test.test_docs;`)).ToNot(BeNil())
		col, err = nidhi.OpenCollection(ctx, db, "collection_test", "test_doc", nidhi.CollectionOptions{
			Fields: []string{"Id", "Number"},
		})
		Expect(col, err).NotTo(BeNil())
	})

	AfterSuite(func() {
		Expect(col.DeleteMany(ctx, nil, []nidhi.DeleteOption{nidhi.WithDeleteOptions(nidhi.DeleteOptions{Permanent: true})})).To(Succeed())
		db.Close()
	})

	Context("single document operations", func() {
		var doc testDoc
		BeforeEach(func() {
			doc = testDoc{Id: strconv.Itoa(count), Number: rand.Int()}
			count++
			Expect(col.Create(ctx, &doc, nil)).To(Equal(doc.Id))
		})

		It("should get a document by its id", func() {
			var act testDoc
			Expect(col.Get(ctx, doc.Id, &act, nil)).To(Succeed())
			Expect(act).To(Equal(doc))
		})

		It("should get a partial document by it's id", func() {
			var act testDoc
			Expect(col.Get(ctx, doc.Id, &act, []nidhi.GetOption{nidhi.WithGetOptions(nidhi.GetOptions{ViewMask: []string{"Number"}})})).To(Succeed())
			exp := doc
			exp.Id = ""
			Expect(act).To(Equal(exp))
		})

		It("should delete a document by its id", func() {
			Expect(col.Delete(ctx, doc.Id, nil)).To(Succeed())
			var act testDoc
			Expect(col.Get(ctx, doc.Id, &act, nil)).ToNot(Succeed())
		})

		Context("Upsert", func() {
			It("should not be allowed when not passing replace option", func() {
				_, err := col.Create(ctx, &doc, nil)
				Expect(err).ToNot(Succeed())
			})

			It("should replace if replace is passed", func() {
				exp := doc
				exp.Number = doc.Number + 1
				Expect(col.Create(ctx, &exp, []nidhi.CreateOption{nidhi.WithCreateOptions(nidhi.CreateOptions{Replace: true})})).To(Equal(doc.Id))
				var act testDoc
				Expect(col.Get(ctx, exp.Id, &act, nil)).To(Succeed())
				Expect(act).To(Equal(exp))
			})
		})

		It("should be able to replace a document", func() {
			exp := doc
			exp.Number = doc.Number + 2
			Expect(col.Replace(ctx, &exp, nil)).To(Succeed())
			var act testDoc
			Expect(col.Get(ctx, exp.Id, &act, nil)).To(Succeed())
			Expect(act).To(Equal(exp))
		})
	})

	Context("multi document operations", func() {
		var (
			docs        []*testDoc
			aboveMarker []*testDoc
		)
		const marker = 5

		BeforeEach(func() {
			Expect(col.DeleteMany(ctx, nil, []nidhi.DeleteOption{nidhi.WithDeleteOptions(nidhi.DeleteOptions{Permanent: true})})).To(Succeed())
			docs = make([]*testDoc, 1+(rand.Int()%20))
			aboveMarker = make([]*testDoc, 0, len(docs))

			for i := range docs {
				docs[i] = &testDoc{Id: strconv.Itoa(count), Number: rand.Int() % 10}
				count++
				Expect(col.Create(ctx, docs[i], nil)).To(Equal(docs[i].Id))
				if docs[i].Number > marker {
					aboveMarker = append(aboveMarker, docs[i])
				}
			}
		})

		qf := func(opts ...nidhi.QueryOption) ([]*testDoc, error) {
			var act []*testDoc
			f := &testQuery{}
			f.Number(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))})
			return act, col.Query(
				ctx,
				newTestQuery().Number(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))}),
				func() nidhi.Document {
					var doc testDoc
					act = append(act, &doc)
					return &doc
				},
				opts,
			)
		}

		It("returns results based on a query", func() {
			exp := aboveMarker
			act, err := qf()
			Expect(err).To(BeNil())
			Expect(act).To(Equal(exp))
		})

		It("returns results with a partial view based on a query", func() {
			exp := aboveMarker
			for i := range exp {
				exp[i].Id = ""
			}
			act, err := qf(nidhi.WithQueryOptions(nidhi.QueryOptions{ViewMask: []string{"Number"}}))
			Expect(err).To(BeNil())
			Expect(act).To(Equal(exp))
		})

		Context("Pagination", func() {

			It("has more", func() {
				sort.Sort(byNumber(aboveMarker))
				exp := aboveMarker[:len(aboveMarker)-2]
				po := &nidhi.PaginationOptions{}
				po.Limit = uint64(len(exp))
				po.Cursor = ""
				po.OrderBy = []nidhi.OrderBy{{Field: nidhi.OrderByInt("(" + nidhi.ColDoc + "->'Number')::bigint")}}
				act, err := qf(nidhi.WithQueryOptions(nidhi.QueryOptions{PaginationOptions: po}))
				Expect(err).To(BeNil())
				Expect(act).To(Equal(exp))
				Expect(po.HasMore).To(BeTrue())
			})

			pf := func(backward bool) {
				cursor := ""
				var act []*testDoc
				for {
					po := &nidhi.PaginationOptions{
						Backward: backward,
						Limit:    1,
						Cursor:   cursor,
						OrderBy:  []nidhi.OrderBy{{Field: nidhi.OrderByInt("(" + nidhi.ColDoc + "->'Number')::bigint")}},
					}
					pr, err := qf(nidhi.WithPaginationOptions(po))
					Expect(err).To(BeNil())
					Expect(len(pr)).To(Equal(1))
					act = append(act, pr...)

					cursor = po.NextCursor

					if !po.HasMore {
						break
					}
				}

				if backward {
					sort.Sort(sort.Reverse(byId(aboveMarker)))
				} else {
					sort.Sort(byId(aboveMarker))
				}

				Expect(act).To(Equal(aboveMarker))
			}

			It("Should Paginate forward", func() {
				pf(false)
			})

			It("Should Paginate backward", func() {
				pf(true)
			})

			It("Does not have more", func() {
				po := &nidhi.PaginationOptions{}
				exp := aboveMarker
				po.Limit = uint64(len(exp))
				po.Cursor = ""
				po.OrderBy = []nidhi.OrderBy{{Field: nidhi.OrderByInt("(" + nidhi.ColDoc + "->'Number')::bigint")}}
				_, err := qf(nidhi.WithQueryOptions(nidhi.QueryOptions{PaginationOptions: po}))
				Expect(err).To(BeNil())
				Expect(po.HasMore).To(BeFalse())
			})
		})

		It("update based on query", func() {
			Expect(col.Update(
				ctx,
				&testDoc{
					Number: -1,
				},
				newTestQuery().Number(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))}),
				nil,
			)).To(Succeed())

			var act []*testDoc
			Expect(col.Query(
				ctx,
				newTestQuery().Number(&nidhi.IntQuery{Eq: nidhi.Int64(-1)}),
				func() nidhi.Document {
					var doc testDoc
					act = append(act, &doc)
					return &doc
				},
				nil,
			)).To(Succeed())

			exp := aboveMarker
			for _, e := range exp {
				e.Number = -1
			}

			Expect(act).To(Equal(exp))
		})

	})

})
