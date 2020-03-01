package nidhi_test

import (
	"context"
	"database/sql"
	"math/rand"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/srikrsna/nidhi"
	"gocloud.dev/postgres"
)

var _ = Describe("Collection", func() {
	var (
		db  *sql.DB
		col *nidhi.Collection
		ctx = context.TODO()
	)
	BeforeSuite(func() {
		var err error
		db, err = postgres.Open(ctx, "postgres://krsna@localhost/test?sslmode=disable")
		Expect(db, err).NotTo(BeNil())
		Expect(db.Ping()).To(Succeed())
		Expect(db.Exec(`TRUNCATE collection_test.test_docs;`)).ToNot(BeNil())
		col, err = nidhi.OpenCollection(ctx, db, "collection_test", "test_docs")
		Expect(col, err).NotTo(BeNil())
	})

	AfterSuite(func() {
		Expect(col.DeleteMany(ctx, nil, []nidhi.DeleteOption{nidhi.WithDeleteOptions(nidhi.DeleteOptions{Permanent: true})})).To(Succeed())
		db.Close()
	})

	Context("single document operations", func() {
		var doc testDoc
		BeforeEach(func() {
			doc = testDoc{Id: uuid.New().String(), Number: rand.Int()}
			Expect(col.Create(ctx, &doc, nil)).To(Equal(doc.Id))
		})

		It("should get a document by its id", func() {
			var act testDoc
			Expect(col.Get(ctx, doc.Id, &act, nil)).To(Succeed())
			Expect(act).To(Equal(doc))
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
				Expect(col.Create(ctx, &exp, []nidhi.CreateOption{nidhi.WithCreateOptions(nidhi.CreateOptions{ReplaceIfExists: true})})).To(Equal(doc.Id))
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
				docs[i] = &testDoc{Id: uuid.New().String(), Number: rand.Int() % 10}
				Expect(col.Create(ctx, docs[i], nil)).To(Equal(docs[i].Id))
				if docs[i].Number > marker {
					aboveMarker = append(aboveMarker, docs[i])
				}
			}
		})

		It("returns results based on a query", func() {
			exp := aboveMarker
			act := make([]*testDoc, 0, len(docs))
			Expect(col.Query(
				ctx,
				&testFilter{
					Number: &nidhi.IntFilter{Gt: nidhi.Int64(int64(marker))},
				},
				func() nidhi.Document {
					var doc testDoc
					act = append(act, &doc)
					return &doc
				},
				nil,
			)).To(Succeed())
			Expect(act).To(Equal(exp))
		})

		It("count documents based on a query", func() {
			Expect(col.Count(ctx, &testFilter{
				Number: &nidhi.IntFilter{Gt: nidhi.Int64(int64(marker))},
			}, nil)).To(Equal(int64(len(aboveMarker))))
		})
	})

})
