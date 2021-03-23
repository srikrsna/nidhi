package pgn_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gocloud.dev/postgres"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	fuzz "github.com/google/gofuzz"
	"github.com/srikrsna/protoc-gen-fuzz/wkt"

	"github.com/srikrsna/nidhi"
	pb "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi/test_data"
	pbf "github.com/srikrsna/nidhi/internal/protoc-gen-nidhi/test_data/fuzz"
	activitylogs "github.com/srikrsna/nidhi/metadata/activity-logs"
)

func TestQuery(t *testing.T) {
	type row struct {
		Name  string
		Query nidhi.Sqlizer
		SQL   string
		Args  []interface{}
	}

	table := []row{
		{
			Name:  "Simple String",
			Query: pb.GetAllQuery().StringField(&nidhi.StringQuery{Eq: nidhi.String("Yes")}),
			SQL:   nidhi.ColDoc + `->>'stringField' = ?`,
			Args:  []interface{}{"Yes"},
		},
		{
			Name:  "Simple Int",
			Query: pb.GetAllQuery().Int32Field(&nidhi.IntQuery{Eq: nidhi.Int64(1)}),
			SQL:   "(" + nidhi.ColDoc + `->'int32Field')::bigint = ?`,
			Args:  []interface{}{int64(1)},
		},
		{
			Name:  "And",
			Query: pb.GetAllQuery().StringField(&nidhi.StringQuery{Eq: nidhi.String("1")}).And().BoolField(&nidhi.BoolQuery{Eq: nidhi.Bool(true)}),
			SQL:   nidhi.ColDoc + `->>'stringField' = ? AND (` + nidhi.ColDoc + `->'boolField')::bool = ?`,
			Args:  []interface{}{"1", true},
		},
		{
			Name:  "Or",
			Query: pb.GetAllQuery().StringField(&nidhi.StringQuery{Eq: nidhi.String("1")}).Or().BoolField(&nidhi.BoolQuery{Eq: nidhi.Bool(true)}),
			SQL:   nidhi.ColDoc + `->>'stringField' = ? OR (` + nidhi.ColDoc + `->'boolField')::bool = ?`,
			Args:  []interface{}{"1", true},
		},
	}

	for _, row := range table {
		t.Run(row.Name, func(t *testing.T) {
			actSQL, actArgs, err := row.Query.ToSql()
			if err != nil {
				return
			}

			if strings.TrimSpace(actSQL) != row.SQL {
				t.Fatalf("sql mismatch, exp: %s, act: %s", row.SQL, actSQL)
			}

			if !reflect.DeepEqual(actArgs, row.Args) {
				t.Fatalf("args mismatch, exp: %v, act: %v", row.Args, actArgs)
			}
		})
	}
}

func TestJson(t *testing.T) {
	fz := fuzz.New().Funcs(pbf.FuzzFuncs()...).Funcs(wkt.FuzzWKT[:]...).NilChance(.5).NumElements(0, 70)
	cfg := jsoniter.Config{
		IndentionStep: 4,
	}.Froze()
	stream := cfg.BorrowStream(nil)
	iterator := cfg.BorrowIterator(nil)
	for i := 0; i < 10000; i++ {
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

		if !cmp.Equal(&act, &doc, protocmp.Transform()) {
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
		fz  = fuzz.New().Funcs(pbf.FuzzFuncs()...).Funcs(wkt.FuzzWKT[:]...)
	)
	BeforeSuite(func() {
		var err error
		db, err = postgres.Open(ctx, "postgres://krsna@localhost/postgres?sslmode=disable")
		Expect(db, err).NotTo(BeNil())
		Expect(db.Ping()).To(Succeed())
		Expect(db.Exec(`DROP TABLE IF EXISTS pb.alls;`)).ToNot(BeNil())
		col, err = pb.OpenAllCollection(ctx, db, activitylogs.Provider(func(ctx context.Context) string {
			return "srikrsna"
		}))
		Expect(col, err).NotTo(BeNil())
	})

	AfterSuite(func() {
		_ = db.Close()
	})

	Context("single document operations", func() {
		var doc pb.All
		BeforeEach(func() {
			fz.Fuzz(&doc)
			doc.Id = uuid.New().String()
			doc.BytesField = nil
			s := jsoniter.ConfigDefault.BorrowStream(nil)
			Expect(doc.MarshalDocument(s)).To(Succeed())
			Expect(col.CreateAll(ctx, &doc)).To(Equal(doc.Id), "%#v", string(s.Buffer()))
		})

		It("should get a document by its id", func() {
			Expect(col.GetAll(ctx, doc.Id)).To(ProtoEqual(&doc))
		})

		It("should get a partial document by it's id", func() {
			exp := &pb.All{}
			exp.StringField = doc.StringField
			Expect(col.GetAll(ctx, doc.Id, nidhi.WithGetOptions(nidhi.GetOptions{ViewMask: []string{"stringField"}}))).To(ProtoEqual(exp))
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
				Expect(col.GetAll(ctx, exp.Id)).To(ProtoEqual(exp))
			})
		})

		Context("Metadata", func() {
			It("Get the metadata", func() {
				var md activitylogs.Metadata
				_, err := col.GetAll(ctx, doc.Id, nidhi.WithGetMetadata(&md))
				Expect(err).To(BeNil())
				Expect(md.Created).ToNot(BeNil())
			})
		})

		It("should be able to replace a document", func() {
			exp := proto.Clone(&doc).(*pb.All)
			exp.Int32Field = doc.Int32Field + 2
			Expect(col.ReplaceAll(ctx, exp)).To(Succeed())
			Expect(col.GetAll(ctx, exp.Id)).To(ProtoEqual(exp))
		})
	})

	Context("multi document operations", func() {
		var (
			ct          = time.Unix(time.Now().Add(-time.Second).Unix(), 0)
			docs        []*pb.All
			aboveMarker []*pb.All
		)
		const marker = 5

		BeforeEach(func() {
		loop:
			docs = make([]*pb.All, 1+(rand.Int()%20)) //nolint:gosec
			aboveMarker = make([]*pb.All, 0, len(docs))
			Expect(db.Exec(`TRUNCATE TABLE pb.alls;`)).ToNot(BeNil())
			for i := range docs {
				var doc pb.All
				fz.Fuzz(&doc)
				doc.Id = strconv.Itoa(i)
				docs[i] = &doc
				docs[i].Timestamp = timestamppb.New(ct)
				Expect(col.CreateAll(ctx, docs[i])).To(Equal(docs[i].Id))
				if docs[i].Int32Field > marker {
					aboveMarker = append(aboveMarker, docs[i])
				}
			}
			if len(aboveMarker) <= 5 {
				goto loop
			}
		})

		qf := func(opts ...nidhi.QueryOption) ([]*pb.All, error) {
			return col.QueryAlls(
				ctx,
				pb.GetAllQuery().Int32Field(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))}),
				opts...,
			)
		}

		qfe := func(exp []*pb.All, opts ...nidhi.QueryOption) {
			Expect(qf(opts...)).To(AllSliceEqual(exp))
		}

		It("returns results based on a query", func() {
			exp := aboveMarker
			qfe(exp)
		})

		It("returns values based on a subquery of time", func() {
			exp := docs
			Expect(col.QueryAlls(ctx, pb.GetAllQuery().Timestamp(&nidhi.TimeQuery{Eq: &ct}))).To(AllSliceEqual(exp))
		})

		It("returns results with a partial view based on a query", func() {
			exp := make([]*pb.All, 0, len(aboveMarker))
			for i := range aboveMarker {
				var e pb.All
				e.Int32Field = aboveMarker[i].Int32Field
				exp = append(exp, &e)
			}
			qfe(exp, nidhi.WithQueryOptions(nidhi.QueryOptions{ViewMask: []string{"int32Field"}}))
		})

		It("orders based on the field passed", func() {
			exp := make([]*pb.All, 0, len(aboveMarker))
			for i := range aboveMarker {
				var e pb.All
				e.Int32Field = aboveMarker[i].Int32Field
				exp = append(exp, &e)
			}
			sort.Sort(byInt32Field(exp))
			qfe(exp, nidhi.WithQueryOptions(
				nidhi.QueryOptions{
					ViewMask: []string{"int32Field"},
					PaginationOptions: &nidhi.PaginationOptions{
						Limit: uint64(len(exp)),
						OrderBy: []nidhi.OrderBy{
							{
								Field: pb.AllSchema().Int32Field(),
							},
						},
					},
				},
			))
		})

		Context("Pagination", func() {

			It("has more", func() {
				sort.Sort(byId(aboveMarker))
				exp := aboveMarker[:1]
				po := &nidhi.PaginationOptions{}
				po.Limit = 1
				qfe(exp, nidhi.WithQueryOptions(nidhi.QueryOptions{PaginationOptions: po}))
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

				Expect(act).To(AllSliceEqual(aboveMarker))
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
			q := pb.GetAllQuery().Int32Field(&nidhi.IntQuery{Gt: nidhi.Int64(int64(marker))})
			Expect(col.UpdateAlls(
				ctx,
				&pb.All{
					FloatField: -1,
				},
				q,
			)).To(Succeed())

			act, err := col.QueryAlls(
				ctx,
				q,
			)
			Expect(err).To(BeNil())

			exp := make([]*pb.All, 0, len(aboveMarker))
			for _, e := range aboveMarker {
				a := proto.Clone(e).(*pb.All)
				a.FloatField = -1
				exp = append(exp, a)
			}

			Expect(act).To(AllSliceEqual(exp))
		})

		Context("Metadata", func() {
			It("should fetch activity logs on query", func() {
				exp := aboveMarker
				var mdc activitylogs.Creator
				qfe(exp, nidhi.WithQueryCreateMetadata(mdc.Create))
				Expect(len(mdc.Values)).To(Equal(len(exp)))
			})

			It("Fetches all with md user query", func() {
				exp := docs
				act, err := col.QueryAlls(ctx, pb.GetAllQuery().WhereMetadata(
					activitylogs.CreatedBy("srikrsna"),
				))
				Expect(act, err).To(AllSliceEqual(exp))
			})

			It("Fetches all with md time query", func() {
				exp := docs
				act, err := col.QueryAlls(ctx, pb.GetAllQuery().WhereMetadata(
					activitylogs.CreatedAfter(ct),
				))
				Expect(act, err).To(AllSliceEqual(exp))
			})
		})
	})

})
