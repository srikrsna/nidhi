package manual

import (
	"context"
	"database/sql"
	"errors"

	jsoniter "github.com/json-iterator/go"

	"github.com/srikrsna/nidhi"
)

type nidhiCol interface {
	Create(ctx context.Context, doc nidhi.Document, ops []nidhi.CreateOption) (string, error)
	Replace(ctx context.Context, doc nidhi.Document, ops []nidhi.ReplaceOption) error
	Update(ctx context.Context, doc nidhi.Document, f nidhi.Sqlizer, ops []nidhi.UpdateOption) error
	Delete(ctx context.Context, id string, ops []nidhi.DeleteOption) error
	DeleteMany(ctx context.Context, f nidhi.Sqlizer, ops []nidhi.DeleteOption) error
	Query(ctx context.Context, f nidhi.Sqlizer, ctr func() nidhi.Document, ops []nidhi.QueryOption) error
	Get(ctx context.Context, id string, doc nidhi.Document, ops []nidhi.GetOption) error
	Count(ctx context.Context, f nidhi.Sqlizer, ops []nidhi.CountOption) (int64, error)
}

type Book struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`

	Author    *Author `json:"author,omitempty"`
	PageCount int     `json:"pageCount,omitempty"`
	Pages     []*Page `json:"pages,omitempty"`
}

func (doc *Book) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true

	w.WriteObjectStart()
	first = WriteString(w, "id", doc.Id, first)
	first = WriteString(w, "title", doc.Title, first)
	if doc.Author != nil {
		first = WriteMarshaler(w, "author", doc.Author, first)
	}
	first = WriteInt(w, "pageCount", doc.PageCount, first)
	if len(doc.Pages) > 0 {
		if !first {
			w.WriteMore()
		}
		w.WriteObjectField("pages")
		w.Error = PageSlice(doc.Pages).MarshalDocument(w)
	}
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
			doc.PageCount = r.ReadInt()
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

func (b *Book) SetDocumentId(id string) {
	b.Id = id
}

func (b *Book) DocumentId() string {
	return b.Id
}

type Author struct {
	Name string `json:"name,omitempty"`
	Bio  string `json:"bio,omitempty"`
}

func (doc *Author) MarshalDocument(w *jsoniter.Stream) error {
	if doc == nil {
		w.WriteNil()
		return w.Error
	}

	first := true
	w.WriteObjectStart()
	first = WriteString(w, "name", doc.Name, first)
	WriteString(w, "bio", doc.Bio, first)
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

type Page struct {
	Number  int    `json:"number,omitempty"`
	Content string `json:"content,omitempty"`
}

func (doc *Page) UnmarshalDocument(r *jsoniter.Iterator) error {
	if doc == nil {
		return errors.New("empty object passed")
	}

	r.ReadObjectCB(func(r *jsoniter.Iterator, field string) bool {
		switch field {
		case "number":
			doc.Number = r.ReadInt()
		case "content":
			doc.Content = r.ReadString()
		default:
			r.Skip()
		}
		return true
	})

	return r.Error
}

func (p *Page) MarshalDocument(w *jsoniter.Stream) error {
	if p == nil {
		w.WriteNil()
		return w.Error
	}

	first := true
	w.WriteObjectStart()
	first = WriteInt(w, "number", p.Number, first)
	WriteString(w, "content", p.Content, first)
	w.WriteObjectEnd()

	return w.Error
}

type PageSlice []*Page

func (s PageSlice) MarshalDocument(w *jsoniter.Stream) error {
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

func GetBookQuery() BookQuery {
	return (*bookQuery)(nidhi.GetQuery())
}

func PutBookQuery(q BookQuery) {
	nidhi.PutQuery((*nidhi.Query)(q.(*bookQuery)))
}

type isBookQuery interface {
	bookQuery()
	nidhi.Sqlizer
}

type BookQuery interface {
	Id(*nidhi.StringQuery) BookConj
	Title(*nidhi.StringQuery) BookConj
	PageCount(*nidhi.IntQuery) BookConj
	Author() BookAuthorQuery

	// Generic With Type Safety
	Paren(iq isBookQuery) BookConj
	Where(query string, args ...interface{}) BookConj
	Not() BookQuery
	ReplaceArgs(args ...interface{}) error
}

type bookQuery nidhi.Query

var (
	_ BookQuery       = (*bookQuery)(nil)
	_ BookConj        = (*bookQuery)(nil)
	_ BookAuthorQuery = (*bookAuthorQuery)(nil)
)

type BookConj interface {
	And() BookQuery
	Or() BookQuery
	isBookQuery
}

type BookAuthorQuery interface {
	Name(*nidhi.StringQuery) BookConj
	Bio(*nidhi.StringQuery) BookConj
}

type bookAuthorQuery nidhi.Query

func (q *bookAuthorQuery) Name(f *nidhi.StringQuery) BookConj {
	(*nidhi.Query)(q).Field("->'name'", f)
	return (*bookQuery)(q)
}

func (q *bookAuthorQuery) Bio(f *nidhi.StringQuery) BookConj {
	(*nidhi.Query)(q).Field("->'bio'", f)
	return (*bookQuery)(q)
}

func (q *bookQuery) bookQuery() {}

func (q *bookQuery) Id(f *nidhi.StringQuery) BookConj {
	q.qry().Id(f)
	return q
}

func (q *bookQuery) Title(f *nidhi.StringQuery) BookConj {
	q.qry().Field(" "+nidhi.ColDoc+"->>'title'", f)
	return q
}

func (q *bookQuery) PageCount(f *nidhi.IntQuery) BookConj {
	q.qry().Field(" "+nidhi.ColDoc+"->'pageCount'", f)
	return q
}

func (q *bookQuery) Author() BookAuthorQuery {
	q.qry().Prefix(" " + nidhi.ColDoc + "->>'author'")
	return (*bookAuthorQuery)(q)
}

func (q *bookQuery) Pages(pp ...*Page) BookConj {
	q.qry().Field(" "+nidhi.ColDoc+"->'pages'", nidhi.MarshalerQuery{
		Marshaler: PageSlice(pp),
	})
	return q
}

func (q *bookQuery) Paren(iq isBookQuery) BookConj {
	q.qry().Paren(iq)
	return q
}

func (q *bookQuery) Where(query string, args ...interface{}) BookConj {
	q.qry().Where(query, args...)
	return q
}

func (q *bookQuery) Not() BookQuery {
	q.qry().Not()
	return q
}

func (q *bookQuery) And() BookQuery {
	q.qry().And()
	return q
}

func (q *bookQuery) Or() BookQuery {
	q.qry().Or()
	return q
}

func (q *bookQuery) ReplaceArgs(args ...interface{}) error {
	return q.qry().ReplaceArgs()
}

func (q *bookQuery) ToSql() (string, []interface{}, error) {
	return q.qry().ToSql()
}

func (q *bookQuery) qry() *nidhi.Query {
	return (*nidhi.Query)(q)
}

type BookCollection struct {
	*bookCollection

	ogCol *nidhi.Collection
}

func OpenBookCollection(ctx context.Context, db *sql.DB) (*BookCollection, error) {
	col, err := nidhi.OpenCollection(ctx, db, "books_v1", "books", nidhi.CollectionOptions{
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

func (bs *BookCollection) BeginTx(ctx context.Context, opt *sql.TxOptions) (*BookTxCollection, error) {
	txCol, err := bs.ogCol.BeginTx(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &BookTxCollection{&bookCollection{txCol}, txCol}, nil
}

func (bs *BookCollection) WithTransaction(tx *nidhi.TxToken) *BookTxCollection {
	txCol := bs.ogCol.WithTransaction(tx)
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
	col nidhiCol
}

func (bs *bookCollection) CreateBook(ctx context.Context, b *Book, ops ...nidhi.CreateOption) (string, error) {
	return bs.col.Create(ctx, b, ops)
}

func (bs *bookCollection) QueryBooks(ctx context.Context, f isBookQuery, ops ...nidhi.QueryOption) ([]*Book, error) {
	var ee []*Book
	ctr := func() nidhi.Document {
		var e Book
		ee = append(ee, &e)
		return &e
	}

	return ee, bs.col.Query(ctx, f, ctr, ops)
}

func (bs *bookCollection) ReplaceBook(ctx context.Context, b *Book, ops ...nidhi.ReplaceOption) error {
	return bs.col.Replace(ctx, b, ops)
}

func (bs *bookCollection) DeleteBook(ctx context.Context, id string, ops ...nidhi.DeleteOption) error {
	return bs.col.Delete(ctx, id, ops)
}

func (bs *bookCollection) CountBooks(ctx context.Context, f isBookQuery, ops ...nidhi.CountOption) (int64, error) {
	return bs.col.Count(ctx, f, ops)
}

func (bs *bookCollection) GetBook(ctx context.Context, id string, ops ...nidhi.GetOption) (*Book, error) {
	var entity Book
	return &entity, bs.col.Get(ctx, id, &entity, ops)
}

func (bs *bookCollection) UpdateBooks(ctx context.Context, b *Book, f isBookQuery, ops ...nidhi.UpdateOption) error {
	return bs.col.Update(ctx, b, f, ops)
}

func WriteString(w *jsoniter.Stream, field, value string, first bool) bool {
	if value == "" {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteString(value)

	return false
}

func WriteBool(w *jsoniter.Stream, field string, value, first bool) bool {
	if !value {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteBool(value)

	return false
}

func WriteInt(w *jsoniter.Stream, field string, value int, first bool) bool {
	if value == 0 {
		return false
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteInt(value)

	return false
}

func WriteFloat32(w *jsoniter.Stream, field string, value float32, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteFloat32(value)

	return false
}

func WriteFloat64(w *jsoniter.Stream, field string, value float64, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteFloat64(value)

	return false
}

func WriteInt32(w *jsoniter.Stream, field string, value int32, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}
	w.WriteObjectField(field)
	w.WriteInt32(value)

	return false
}

func WriteInt64(w *jsoniter.Stream, field string, value int64, first bool) bool {
	if value == 0 {
		return first
	}

	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	w.WriteInt64(value)

	return false
}

func WriteMarshaler(w *jsoniter.Stream, field string, value nidhi.Marshaler, first bool) bool {
	if !first {
		w.WriteMore()
	}

	w.WriteObjectField(field)
	value.MarshalDocument(w)

	return false
}
