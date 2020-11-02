// Code generated by protoc-gen-nidhi. DO NOT EDIT.
// source: internal/protoc-gen-nidhi/test_data/simple/simple.proto

package pb

import (
	"context"
	"database/sql"
	"errors"

	jsoniter "github.com/json-iterator/go"

	"github.com/srikrsna/nidhi"
	nidhigen "github.com/srikrsna/nidhi/nidhigen"
)

var (
	_ = context.Background
	_ = (*sql.DB)(nil)
	_ = errors.New
	_ = jsoniter.Marshal
	_ = nidhi.JSONB
	_ = nidhigen.WriteString
)

func (doc *Book) DocumentId() string {
	return doc.Id
}

func (doc *Book) SetDocumentId(id string) {
	doc.Id = id
}

type BookCollection struct {
	*bookCollection

	ogCol *nidhi.Collection
}

func OpenBookCollection(ctx context.Context, db *sql.DB) (*BookCollection, error) {
	col, err := nidhi.OpenCollection(ctx, db, "pb", "books", nidhi.CollectionOptions{
		Fields: []string{"id", "title", "author", "pageCount", "pages"},
	})
	if err != nil {
		return nil, err
	}
	return &BookCollection{
		&bookCollection{col: col},
		col,
	}, nil
}

func (st *BookCollection) BeginTx(ctx context.Context, opt *sql.TxOptions) (*BookTxCollection, error) {
	txCol, err := st.ogCol.BeginTx(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &BookTxCollection{&bookCollection{txCol}, txCol}, nil
}

func (st *BookCollection) WithTransaction(tx *nidhi.TxToken) *BookTxCollection {
	txCol := st.ogCol.WithTransaction(tx)
	return &BookTxCollection{&bookCollection{txCol}, txCol}
}

type BookTxCollection struct {
	*bookCollection
	txCol *nidhi.TxCollection
}

func (tx *BookTxCollection) Rollback() error {
	return tx.txCol.Rollback()
}

func (tx *BookTxCollection) Commit() error {
	return tx.txCol.Commit()
}

func (tx *BookTxCollection) TxToken() *nidhi.TxToken {
	return nidhi.NewTxToken(tx.txCol)
}

type bookCollection struct {
	col nidhigen.Collection
}

func (st *bookCollection) CreateBook(ctx context.Context, b *Book, ops ...nidhi.CreateOption) (string, error) {
	return st.col.Create(ctx, b, ops)
}

func (st *bookCollection) QueryBooks(ctx context.Context, f isBookQuery, ops ...nidhi.QueryOption) ([]*Book, error) {
	var ee []*Book
	ctr := func() nidhi.Document {
		var e Book
		ee = append(ee, &e)
		return &e
	}

	return ee, st.col.Query(ctx, f, ctr, ops)
}

func (st *bookCollection) ReplaceBook(ctx context.Context, b *Book, ops ...nidhi.ReplaceOption) error {
	return st.col.Replace(ctx, b, ops)
}

func (st *bookCollection) DeleteBook(ctx context.Context, id string, ops ...nidhi.DeleteOption) error {
	return st.col.Delete(ctx, id, ops)
}

func (st *bookCollection) CountBooks(ctx context.Context, f isBookQuery, ops ...nidhi.CountOption) (int64, error) {
	return st.col.Count(ctx, f, ops)
}

func (st *bookCollection) GetBook(ctx context.Context, id string, ops ...nidhi.GetOption) (*Book, error) {
	var entity Book
	return &entity, st.col.Get(ctx, id, &entity, ops)
}

func (st *bookCollection) UpdateBooks(ctx context.Context, b *Book, f isBookQuery, ops ...nidhi.UpdateOption) error {
	return st.col.Update(ctx, b, f, ops)
}

func GetBookQuery() BookQuery {
	return (*bookQuery)(nidhi.GetQuery())
}

func PutBookQuery(q BookQuery) {
	nidhi.PutQuery((*nidhi.Query)(q.(*bookQuery)))
}

type BookQuery interface {
	Id(*nidhi.StringQuery) BookConj
	Title(*nidhi.StringQuery) BookConj
	Author() BookAuthorQuery
	PageCount(*nidhi.IntQuery) BookConj
	Pages(...*Page) BookConj

	// Generic With Type Safety
	Paren(iq isBookQuery) BookConj
	Where(query string, args ...interface{}) BookConj
	Not() BookQuery
	ReplaceArgs(args ...interface{}) error
}

type BookConj interface {
	And() BookQuery
	Or() BookQuery
	isBookQuery
}

type isBookQuery interface {
	bookQuery()
	nidhi.Sqlizer
}

type bookQuery nidhi.Query

func (q *bookQuery) bookQuery() {}

func (q *bookQuery) Id(f *nidhi.StringQuery) BookConj {
	(*nidhi.Query)(q).Id(f)
	return q
}
func (q *bookQuery) Title(f *nidhi.StringQuery) BookConj {
	(*nidhi.Query)(q).Field(" "+nidhi.ColDoc+"->'title'", f)
	return q
}

func (q *bookQuery) Author() BookAuthorQuery {
	(*nidhi.Query)(q).Prefix(" " + nidhi.ColDoc + "->>'author'")
	return (*bookAuthorQuery)(q)
}
func (q *bookQuery) PageCount(f *nidhi.IntQuery) BookConj {
	(*nidhi.Query)(q).Field(" "+nidhi.ColDoc+"->'pageCount'", f)
	return q
}

func (q *bookQuery) Pages(arr ...*Page) BookConj {
	(*nidhi.Query)(q).Field(" "+nidhi.ColDoc+"->'pages'", nidhi.MarshalerQuery{
		Marshaler: PageSlice(arr),
	})
	return q
}

func (q *bookQuery) Paren(iq isBookQuery) BookConj {
	(*nidhi.Query)(q).Paren(iq)
	return q
}

func (q *bookQuery) Where(query string, args ...interface{}) BookConj {
	(*nidhi.Query)(q).Where(query, args...)
	return q
}

func (q *bookQuery) Not() BookQuery {
	(*nidhi.Query)(q).Not()
	return q
}

func (q *bookQuery) And() BookQuery {
	(*nidhi.Query)(q).And()
	return q
}

func (q *bookQuery) Or() BookQuery {
	(*nidhi.Query)(q).Or()
	return q
}

func (q *bookQuery) ReplaceArgs(args ...interface{}) error {
	return (*nidhi.Query)(q).ReplaceArgs()
}

func (q *bookQuery) ToSql() (string, []interface{}, error) {
	return (*nidhi.Query)(q).ToSql()
}

type BookAuthorQuery interface {
	Name(*nidhi.StringQuery) BookConj
	Bio(*nidhi.StringQuery) BookConj
}

type bookAuthorQuery nidhi.Query

func (q *bookAuthorQuery) Name(f *nidhi.StringQuery) BookConj {
	(*nidhi.Query)(q).Field(" "+nidhi.ColDoc+"->'name'", f)
	return (*bookQuery)(q)
}

func (q *bookAuthorQuery) Bio(f *nidhi.StringQuery) BookConj {
	(*nidhi.Query)(q).Field(" "+nidhi.ColDoc+"->'bio'", f)
	return (*bookQuery)(q)
}

func (doc *Book) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteString(w, "id", doc.Id, first)
	first = nidhigen.WriteString(w, "title", doc.Title, first)
	first = nidhigen.WriteMarshaler(w, "author", doc.Author, first)
	first = nidhigen.WriteInt32(w, "pageCount", doc.PageCount, first)
	first = nidhigen.WriteMarshaler(w, "pages", PageSlice(doc.Pages), first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *Book) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "id":
			doc.Id = r.ReadString()
		case "title":
			doc.Title = r.ReadString()
		case "author":
			doc.Author = &Author{}
			r.Error = doc.Author.UnmarshalDocument(r)
		case "pageCount":
			doc.PageCount = r.ReadInt32()
		case "pages":
			doc.Pages = []*Page{}
			r.Error = (*PageSlice)(&doc.Pages).UnmarshalDocument(r)
		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type BookSlice []*Book

func (s BookSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *BookSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e Book
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}

func (doc *Author) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteString(w, "name", doc.Name, first)
	first = nidhigen.WriteString(w, "bio", doc.Bio, first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *Author) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "name":
			doc.Name = r.ReadString()
		case "bio":
			doc.Bio = r.ReadString()
		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type AuthorSlice []*Author

func (s AuthorSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *AuthorSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e Author
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}

func (doc *Page) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = nidhigen.WriteInt32(w, "number", doc.Number, first)
	first = nidhigen.WriteString(w, "content", doc.Content, first)
	w.WriteObjectEnd()

	return w.Error
}

func (doc *Page) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "number":
			doc.Number = r.ReadInt32()
		case "content":
			doc.Content = r.ReadString()
		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

type PageSlice []*Page

func (s PageSlice) MarshalDocument(w *jsoniter.Stream) error {
	if len(s) == 0 {
		w.WriteArrayStart()
		w.WriteArrayEnd()
		return nil
	}

	w.WriteArrayStart()
	w.Error = s[0].MarshalDocument(w)
	for _, e := range s[1:] {
		w.WriteMore()
		w.Error = e.MarshalDocument(w)
	}
	w.WriteArrayEnd()

	return w.Error
}

func (s *PageSlice) UnmarshalDocument(r *jsoniter.Iterator) error {
	r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
		var e Page
		r.Error = e.UnmarshalDocument(r)
		*s = append(*s, &e)
		return true
	})

	return r.Error
}
