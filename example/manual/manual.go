package manual

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
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

	w.WriteObjectStart()
	nidhi.WriteString(w, "id", doc.Id)
	nidhi.WriteString(w, "title", doc.Title)
	nidhi.WriteInt(w, "pageCount", doc.PageCount)
	if len(doc.Pages) > 0 {
		w.WriteObjectField("pages")
		w.WriteArrayStart()
		for _, e := range doc.Pages {
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

	w.WriteObjectStart()
	nidhi.WriteString(w, "name", doc.Name)
	nidhi.WriteString(w, "bio", doc.Bio)
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
	w.WriteObjectStart()
	if p.Number != 0 {
		w.WriteObjectField("number")
		w.WriteInt(p.Number)
	}
	if p.Content != "" {
		w.WriteObjectField("content")
		w.WriteString(p.Content)
	}
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

func (f *BookFilter) ToSql(prefix string) (sq.Sqlizer, error) {
	of := &nidhi.ObjectFilter{}
	of.Or = f.Or
	var fp, op string
	if prefix != "" {
		fp = prefix + "->>"
		op = prefix + "->"
		_, _ = fp, op
	}

	pages := make(nidhi.ArrayFilter, 0, len(f.Pages))
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

func (f *AuthorFilter) ToSql(prefix string) (sq.Sqlizer, error) {
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
	col, err := nidhi.OpenCollection(ctx, db, "books_v1", "books", "bk")
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
