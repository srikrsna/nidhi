package manual

import (
	"context"
	"database/sql"
	"errors"

	jsoniter "github.com/json-iterator/go"

	"github.com/srikrsna/nidhi"
)

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
		w.WriteArrayStart()
		doc.Pages[0].MarshalDocument(w)
		for _, e := range doc.Pages[1:] {
			w.WriteMore()
			w.Error = e.MarshalDocument(w)
		}
		w.WriteArrayEnd()
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
			doc.Author.UnmarshalDocument(r)
		case "pageCount":
			doc.PageCount = r.ReadInt()
		case "pages":
			r.ReadArrayCB(func(r *jsoniter.Iterator) bool {
				var e Page
				r.Error = e.UnmarshalDocument(r)
				doc.Pages = append(doc.Pages, &e)
				return true
			})
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

type BookFilter struct {
	Or bool

	Id, Title *nidhi.StringFilter
	PageCount *nidhi.IntFilter
	Author    *AuthorFilter
	Pages     []*Page
}

func (f *BookFilter) ToSql(prefix string) (nidhi.Sqlizer, error) {
	of := &nidhi.ObjectFilter{}
	of.Or = f.Or
	var fp, op string
	if prefix != "" {
		fp = prefix + "->>"
		op = prefix + "->"
		_, _ = fp, op
	}

	pages := make(nidhi.ObjectArrayFilter, 0, len(f.Pages))
	for _, e := range f.Pages {
		pages = append(pages, e)
	}

	of.Filter = map[string]nidhi.Filter{
		fp + "'id'":        f.Id,
		fp + "'title'":     f.Title,
		op + "'author'":    f.Author,
		fp + "'pageCount'": f.PageCount,
		op + "'pages'":     pages,
	}
	return of.ToSql("book")
}

type AuthorFilter struct {
	Or        bool
	Name, Bio *nidhi.StringFilter
}

func (f *AuthorFilter) ToSql(prefix string) (nidhi.Sqlizer, error) {
	of := &nidhi.ObjectFilter{}
	of.Or = f.Or
	var fp, op string
	if prefix != "" {
		fp = prefix + "->>"
		op = prefix + "->"
		_, _ = fp, op
	}
	of.Filter = map[string]nidhi.Filter{
		fp + "'name'": f.Name,
		fp + "'bio'":  f.Bio,
	}
	return of.ToSql("book")
}

type BookCollection struct {
	col *nidhi.Collection
}

func OpenBookCollection(ctx context.Context, db *sql.DB) (*BookCollection, error) {
	col, err := nidhi.OpenCollection(ctx, db, "books_v1", "books", nidhi.CollectionOptions{
		Fields: []string{"id", "title", "author", "pageCount", "pages"},
	})
	if err != nil {
		return nil, err
	}
	return &BookCollection{
		col: col,
	}, nil
}

func (bs *BookCollection) CreateBook(ctx context.Context, b *Book, ops ...nidhi.CreateOption) (string, error) {
	return bs.col.Create(ctx, b, ops)
}

func (bs *BookCollection) QueryBooks(ctx context.Context, f nidhi.Filter, ops ...nidhi.QueryOption) ([]*Book, error) {
	var ee []*Book
	ctr := func() nidhi.Document {
		var e Book
		ee = append(ee, &e)
		return &e
	}

	return ee, bs.col.Query(ctx, f, ctr, ops)
}

func (bs *BookCollection) ReplaceBook(ctx context.Context, b *Book, ops ...nidhi.ReplaceOption) error {
	return bs.col.Replace(ctx, b, ops)
}

func (bs *BookCollection) DeleteBook(ctx context.Context, id string, ops ...nidhi.DeleteOption) error {
	return bs.col.Delete(ctx, id, ops)
}

func (bs *BookCollection) CountBooks(ctx context.Context, f nidhi.Filter, ops ...nidhi.CountOption) (int64, error) {
	return bs.col.Count(ctx, f, ops)
}

func (bs *BookCollection) GetBook(ctx context.Context, id string, ops ...nidhi.GetOption) (*Book, error) {
	var entity Book
	return &entity, bs.col.Get(ctx, id, &entity, ops)
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
