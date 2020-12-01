package pgn_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gocloud.dev/postgres"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	fuzz "github.com/google/gofuzz"

	"github.com/srikrsna/nidhi"
	pb "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi/test_data"
	pbf "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi/test_data/fuzz"
)

func TestQuery(t *testing.T) {
	type row struct {
		Name  string
		Query nidhi.Sqlizer
		Sql   string
		Args  []interface{}
	}

	table := []row{
		{
			Name:  "Simple String",
			Query: pb.GetAllQuery().StringField(&nidhi.StringQuery{Eq: nidhi.String("Yes")}),
			Sql:   nidhi.ColDoc + `->>'stringField' = ?`,
			Args:  []interface{}{"Yes"},
		},
		{
			Name:  "Simple Int",
			Query: pb.GetAllQuery().Int32Field(&nidhi.IntQuery{Eq: nidhi.Int64(1)}),
			Sql:   "(" + nidhi.ColDoc + `->'int32Field')::bigint = ?`,
			Args:  []interface{}{int64(1)},
		},
		{
			Name:  "And",
			Query: pb.GetAllQuery().StringField(&nidhi.StringQuery{Eq: nidhi.String("1")}).And().BoolField(&nidhi.BoolQuery{Eq: nidhi.Bool(true)}),
			Sql:   nidhi.ColDoc + `->>'stringField' = ? AND (` + nidhi.ColDoc + `->'boolField')::bool = ?`,
			Args:  []interface{}{"1", true},
		},
		{
			Name:  "Or",
			Query: pb.GetAllQuery().StringField(&nidhi.StringQuery{Eq: nidhi.String("1")}).Or().BoolField(&nidhi.BoolQuery{Eq: nidhi.Bool(true)}),
			Sql:   nidhi.ColDoc + `->>'stringField' = ? OR (` + nidhi.ColDoc + `->'boolField')::bool = ?`,
			Args:  []interface{}{"1", true},
		},
	}

	for _, row := range table {
		t.Run(row.Name, func(t *testing.T) {
			actSql, actArgs, err := row.Query.ToSql()
			if err != nil {
				return
			}

			if strings.TrimSpace(actSql) != row.Sql {
				t.Fatalf("sql mismatch, exp: %s, act: %s", row.Sql, actSql)
			}

			if !reflect.DeepEqual(actArgs, row.Args) {
				t.Fatalf("args mismatch, exp: %v, act: %v", row.Args, actArgs)
			}
		})
	}
}

func TestJson(t *testing.T) {
	fz := fuzz.New().Funcs(pbf.FuzzFuncs()...)
	cfg := jsoniter.Config{
		IndentionStep: 4,
	}.Froze()
	stream := cfg.BorrowStream(nil)
	iterator := cfg.BorrowIterator(nil)
	for i := 0; i < 10; i++ {
		stream.Reset(nil)

		var doc pb.All
		fz.Fuzz(&doc)

		if err := doc.MarshalDocument(stream); err != nil {
			t.Errorf("unable to marshal document, err: %v", err)
			return
		}

		if !json.Valid(stream.Buffer()) {
			t.Errorf("invalid json: \n%s", string(stream.Buffer()))
			return
		}

		var act pb.All
		iterator.ResetBytes(stream.Buffer())
		if err := act.UnmarshalDocument(iterator); err != nil {
			t.Errorf("unable to unmarshal document, err: %v", err)
			return
		}

		if !proto.Equal(&act, &doc) {
			t.Errorf("mismatch after unmarshal, act: \n%v, exp: \n%v", prototext.Format(&act), prototext.Format(&doc))
			return
		}
	}
}

var _ = Describe("Collection", func() {
	var (
		db  *sql.DB
		col *pb.AllCollection
		ctx = context.TODO()
		fz  = fuzz.New().Funcs(pbf.FuzzFuncs()...)
	)
	BeforeSuite(func() {
		var err error
		db, err = postgres.Open(ctx, "postgres://krsna@localhost/postgres?sslmode=disable")
		Expect(db, err).NotTo(BeNil())
		Expect(db.Ping()).To(Succeed())
		Expect(db.Exec(`DROP TABLE IF EXISTS pb.alls;`)).ToNot(BeNil())
		col, err = pb.OpenAllCollection(ctx, db)
		Expect(col, err).NotTo(BeNil())
	})

	AfterSuite(func() {
		db.Close()
	})

	Context("single document operations", func() {
		var doc pb.All
		BeforeEach(func() {
			fz.Fuzz(&doc)
			doc.Id = uuid.New().String()
			doc.BytesField = nil
			s := jsoniter.ConfigDefault.BorrowStream(nil)
			doc.MarshalDocument(s)
			Expect(col.CreateAll(ctx, &doc)).To(Equal(doc.Id), "%#v", string(s.Buffer()))
		})

		It("should get a document by its id", func() {
			Expect(col.GetAll(ctx, doc.Id)).To(ProtoEqual(&doc))
		})

		It("should get a partial document by it's id", func() {
			act, err := col.GetAll(ctx, doc.Id, nidhi.WithGetOptions(nidhi.GetOptions{ViewMask: []string{"stringField"}}))
			Expect(err).To(BeNil())
			exp := &pb.All{}
			exp.StringField = doc.StringField
			Expect(proto.Equal(act, exp)).To(BeTrue())
		})

		It("should delete a document by its id", func() {
			Expect(col.DeleteAll(ctx, doc.Id)).To(Succeed())
			_, err := col.GetAll(ctx, doc.Id)
			Expect(err).NotTo(BeNil())
		})

		Context("Upsert", func() {
			It("should not be allowed when not passing replace option", func() {
				_, err := col.CreateAll(ctx, &doc)
				Expect(err).ToNot(Succeed())
			})

			It("should replace if replace is passed", func() {
				exp := proto.Clone(&doc).(*pb.All)
				exp.Int32Field = doc.Int32Field + 1
				Expect(col.CreateAll(ctx, exp, nidhi.WithCreateOptions(nidhi.CreateOptions{Replace: true}))).To(Equal(doc.Id))
				act, err := col.GetAll(ctx, exp.Id)
				Expect(err).To(BeNil())
				Expect(proto.Equal(act, exp)).To(BeTrue())
			})
		})

		It("should be able to replace a document", func() {
			exp := proto.Clone(&doc).(*pb.All)
			exp.Int32Field = doc.Int32Field + 2
			Expect(col.ReplaceAll(ctx, exp)).To(Succeed())
			act, err := col.GetAll(ctx, exp.Id)
			Expect(err).To(BeNil())
			Expect(proto.Equal(act, exp)).To(BeTrue())
		})
	})

	Context("multi document operations", func() {
		var (
			docs        []*pb.All
			aboveMarker []*pb.All
		)
		const marker = 5

		BeforeEach(func() {
			Expect(col.DeleteAlls(ctx, nil, nidhi.WithDeleteOptions(nidhi.DeleteOptions{Permanent: true}))).To(Succeed())
			docs = make([]*pb.All, 1+(rand.Int()%20))
			aboveMarker = make([]*pb.All, 0, len(docs))

			for i := range docs {
				var doc pb.All
				fz.Fuzz(&doc)
				docs[i] = &doc
				Expect(col.CreateAll(ctx, docs[i], nil)).To(Equal(docs[i].Id))
				if docs[i].Int32Field > marker {
					aboveMarker = append(aboveMarker, docs[i])
				}
			}
		})

		qf := func(opts ...nidhi.QueryOption) ([]*pb.All, error) {
			return col.QueryAlls(
				ctx,
				pb.GetAllQuery().Int32Field(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))}),
				opts...,
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
				sort.Sort(byId(aboveMarker))
				exp := aboveMarker[:len(aboveMarker)-2]
				po := &nidhi.PaginationOptions{}
				po.Limit = uint64(len(exp))
				po.Cursor = ""
				act, err := qf(nidhi.WithQueryOptions(nidhi.QueryOptions{PaginationOptions: po}))
				Expect(err).To(BeNil())
				Expect(act).To(Equal(exp))
				Expect(po.HasMore).To(BeTrue())
			})

			pf := func(backward bool) {
				cursor := ""
				var act []*pb.All
				for {
					po := &nidhi.PaginationOptions{
						Backward: backward,
						Limit:    1,
						Cursor:   cursor,
					}
					pr, err := qf(nidhi.WithPaginationOptions(po))
					Expect(err).To(BeNil())
					Expect(len(pr)).To(Equal(1))
					act = append(act, pr...)

					cursor = pr[0].Id

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
				_, err := qf(nidhi.WithQueryOptions(nidhi.QueryOptions{PaginationOptions: po}))
				Expect(err).To(BeNil())
				Expect(po.HasMore).To(BeFalse())
			})
		})

		It("update based on query", func() {
			Expect(col.UpdateAlls(
				ctx,
				&pb.All{
					Int32Field: -1,
				},
				pb.GetAllQuery().Int32Field(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))}),
				nil,
			)).To(Succeed())

			act, err := col.QueryAlls(
				ctx,
				pb.GetAllQuery().Int32Field(&nidhi.IntQuery{Eq: nidhi.Int64(-1)}),
			)
			Expect(err).To(BeNil())

			exp := aboveMarker
			for _, e := range exp {
				e.Int32Field = -1
			}

			Expect(act).To(Equal(exp))
		})

		//It("count documents based on a query", func() {
		//	Expect(col.Count(ctx,
		//		newTestQuery().Number(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))}),
		//		nil,
		//	)).To(Equal(int64(len(aboveMarker))))
		//})
	})

})
