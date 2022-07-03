package nidhi

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/elgris/sqrl"
	"github.com/elgris/sqrl/pg"
	"golang.org/x/exp/maps"
)

// Document is wrapper for a resource.
type Document[T any] struct {
	// Value is the resource.
	Value *T
	// Revision is the revision of this document.
	Revision int64
	// Metadata is the metadata of the document
	Metadata Metadata
	// Deleted indicates whether it a deleted document.
	Deleted bool
}

// GetOptions are options for `Get` operation.
type GetOptions struct {
	// ViewMask if set will only populate fields listed in the mask.
	// Only top level fields are supported.
	ViewMask []string
}

// QueryResult is the result of the get call.
type GetResult[T any] struct {
	Document[T]
}

// QueryOptions are options for the `Query` operation
type QueryOptions struct {
	// PaginationOptions if set will return paginated results.
	PaginationOptions *PaginationOptions
	// ViewMask if set will only populate fields listed in the mask.
	// Only top level fields are supported.
	ViewMask []string
	// OrderBy if empty and pagination options field is set
	// will default to sorting by id asc i.e. [{"id", false}]
	OrderBy []OrderBy
	// IncludeDeleted if set will include the soft deleted documents.
	IncludeDeleted bool
	// LoadMetadataParts is slice of metadata parts that need to be loaded.
	LoadMetadataParts []string
}

// PaginationOptions are options for paginating results.
type PaginationOptions struct {
	// Cursor is the pagination cursor that the result should begin from.
	// This is typically returned via the result of the operation.
	Cursor string
	// Limit is the limit of pagination result.
	Limit uint64
	// Backward indicates the direction to fetch results from the cursor.
	// The same result can be achieved by reversing the order.
	//
	// As an example, for documents ordered by their creation time,
	// With the cursor at the 50th record, one can keep fetching the next 50, and the next 50, and so on
	// until they reach the end. Let's say the end is 1000th record.
	// At this point the records can be fetched backwards with the same order.
	Backward bool
}

// OrderBy represents an order by modifer in `Query` operation
type OrderBy struct {
	// Field is the field by which the document should be ordered.
	Field Orderer
	// Desc if set, will order in descending order according to the natural of the
	// field type.
	//
	// Defaults to false. (Ascending)
	Desc bool
}

// QueryResult is the result of the query call.
type QueryResult[T any] struct {
	// Docs is the resulted docs for the query.
	Docs []*Document[T]
	// LastCursor is the token of the last element of the result.
	// It can be used to continue the search result.
	LastCursor string
	// HasMore tells if there are more fields.
	HasMore bool
}

// Get is used to get a document from the store.
func (s *Store[T]) Get(ctx context.Context, id string, opts GetOptions) (*GetResult[T], error) {
	var selection any = ColDoc
	if len(s.fields) > 0 && len(opts.ViewMask) > 0 {
		selection = sq.Expr(ColDoc+" - ?::text[]", pg.Array(difference(s.fields, opts.ViewMask)))
	}
	var (
		docBin, mdBin []byte
		revision      int64
	)
	st := sq.Select().Column(selection).Columns(ColRev, ColMeta).From(s.table).
		Where(sq.Eq{ColId: id}).
		Where(notDeleted)
	if err := st.PlaceholderFormat(sq.Dollar).RunWith(s.db).QueryRowContext(ctx).Scan(&docBin, &revision, &mdBin); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("nidhi: failed to get a document from collection %q, err: %w", s.table, err)
	}
	doc := new(T)
	if err := unmarshalJSON(docBin, doc); err != nil {
		return nil, fmt.Errorf("nidhi: failed to unmarshal document of type %s, err: %w", s.table, err)
	}
	md := make(Metadata, len(s.mdr))
	for k, v := range s.mdr {
		md[k] = v()
	}
	if err := unmarshalJSON(mdBin, md); err != nil {
		return nil, fmt.Errorf("nidhi: failed to unmarshal metadata of parts %v, err: %w", maps.Keys(s.mdr), err)
	}
	return &GetResult[T]{
		Document[T]{
			Value:    doc,
			Metadata: md,
			Revision: revision,
		},
	}, nil
}

// Query queries the store and returns all matching documents.
func (s *Store[T]) Query(ctx context.Context, q Sqlizer, opts QueryOptions) (*QueryResult[T], error) {
	selection := any(ColDoc)
	if len(s.fields) > 0 && len(opts.ViewMask) > 0 {
		selection = sq.Expr(ColDoc+" - ?::text[]", pg.Array(difference(s.fields, opts.ViewMask)))
	}
	st := sq.Select().Column(selection).Column(ColRev).From(s.table)
	var (
		docBin   sql.RawBytes
		revision int64
		scans    = make([]any, 0, 6)
	)
	scans = append(scans, &docBin, &revision)
	// Where Clause
	if any(q) != nil {
		st = st.Where(q)
	}
	var deleted bool
	if !opts.IncludeDeleted {
		st = st.Where(notDeleted)
	} else {
		st = st.Column(ColDel)
		scans = append(scans, &deleted)
	}
	// Optional Metadata
	var mdBin sql.RawBytes
	if len(opts.LoadMetadataParts) > 0 {
		st = st.Column(ColMeta)
		scans = append(scans, &mdBin)
	}
	// Pagination
	st, scans, err := addPagination(st, scans, opts.PaginationOptions, opts.OrderBy)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to paginate: %s, err: %w", s.table, err)
	}
	rows, err := st.PlaceholderFormat(sq.Dollar).RunWith(s.db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to query collection: %s, err: %w", s.table, err)
	}
	defer rows.Close()
	var (
		count   uint64
		hasMore bool
		docs    []*Document[T]
	)
	for rows.Next() {
		if opts.PaginationOptions != nil && opts.PaginationOptions.Limit == count {
			hasMore = true
			break
		}
		count++
		if err := rows.Scan(scans...); err != nil {
			return nil, fmt.Errorf("nidhi: unexpected error while querying collection: %s, err: %w", s.table, err)
		}
		doc := new(T)
		if err := unmarshalJSON(docBin, doc); err != nil {
			return nil, fmt.Errorf("nidhi: failed to unmarshal document of type %s, err: %w", s.table, err)
		}
		var md Metadata
		if len(opts.LoadMetadataParts) > 0 {
			md = make(Metadata, len(opts.LoadMetadataParts))
			for _, part := range opts.LoadMetadataParts {
				md[part] = s.mdr[part]()
			}
			if err := unmarshalJSON(mdBin, md); err != nil {
				return nil, fmt.Errorf("nidhi: failed to unmarshal metadata of parts %v, err: %w", opts.LoadMetadataParts, err)
			}
		}
		docs = append(docs, &Document[T]{
			Value:    doc,
			Revision: revision,
			Metadata: md,
			Deleted:  deleted,
		})
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("nidhi: unexpected error while querying collection: %s, err: %w", s.table, err)
	}
	var cursor string
	if opts.PaginationOptions != nil {
		// Skip doc and revision scan
		scans = scans[2:]
		if len(opts.LoadMetadataParts) > 0 {
			// Skip md scan
			scans = scans[1:]
		}
		if opts.IncludeDeleted {
			// Skip deleted scan
			scans = scans[1:]
		}
		if len(scans) == 1 {
			cursor = *(scans[0].(*string))
		} else {
			cursor = opts.OrderBy[0].Field.Encode(scans[0], *(scans[1].(*string)))
		}
	}
	return &QueryResult[T]{
		Docs:       docs,
		LastCursor: cursor,
		HasMore:    hasMore,
	}, nil
}

func difference[S ~[]T, T comparable](slice1 S, slice2 S) S {
	var diff S
	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}
