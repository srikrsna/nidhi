package nidhi

import (
	"context"
	"fmt"

	sq "github.com/elgris/sqrl"
)

// ReplaceOptions are options for `Replace` operation.
type ReplaceOptions struct {
	// Metadata is the metadata of the document.
	// This will be merged with existing metadata.
	Metadata Metadata
	// Revision if set (>0), the document will only be replaced if the revision also matches.
	Revision int64
}

// UpdateOptions are options for `Update` operation.
type UpdateOptions struct {
	// Metadata is the metadata of the document.
	// This will be merged with existing metadata.
	Metadata Metadata
}

// ReplaceResult is the result of the replace call.
// It doesn't have any fields as of now.
//
// Having an explicit result type will not break future changes.
type ReplaceResult struct{}

// UpdateResult is the result of the update call.
type UpdateResult struct {
	// UpdateCount is the number of documents that were updated.
	UpdateCount int64
}

// Replace replaces a document, matched using it's id, in the store.
//
// Retuns a NotFound error, if the document doesn't exist or the revision doesn't exist.
func (s *Store[T, Q]) Replace(ctx context.Context, doc *T, opts ReplaceOptions) (*ReplaceResult, error) {
	docJSON, err := getJson(doc)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to marshal document of collection: %s, err: %w", s.table, err)
	}
	defer putJson(docJSON)
	mdJSON, err := getJson(opts.Metadata)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to marshal metadata of collection: %s, err: %w", s.table, err)
	}
	defer putJson(mdJSON)
	stmt := s.updateStatement(ctx, docJSON.Buffer(), false, mdJSON.Buffer())
	if opts.Revision > 0 {
		stmt = stmt.Where(sq.Eq{ColRev: opts.Revision})
	}
	stmt = stmt.Where(sq.Eq{ColId: s.idFn(doc)}).Where(notDeleted)
	rc, err := sq.RowsAffected(stmt.PlaceholderFormat(sq.Dollar).RunWith(s.db).ExecContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to put document in database: %w", err)
	}
	if rc != 1 {
		return nil, NotFound
	}
	return &ReplaceResult{}, nil
}

// Update updates all documents of this store that satify the `q`
func (s *Store[T, Q]) Update(ctx context.Context, doc *T, q Q, opts UpdateOptions) (*UpdateResult, error) {
	s.setIdFn(doc, "")
	docJSON, err := getJson(doc)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to marshal document of collection: %s, err: %w", s.table, err)
	}
	defer putJson(docJSON)
	mdJSON, err := getJson(opts.Metadata)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to marshal metadata of collection: %s, err: %w", s.table, err)
	}
	defer putJson(mdJSON)
	st := s.updateStatement(ctx, docJSON.Buffer(), true, mdJSON.Buffer())
	if any(q) != nil {
		st = st.Where(q)
	}
	st = st.Where(notDeleted)
	res, err := st.PlaceholderFormat(sq.Dollar).RunWith(s.db).ExecContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to update documents for collection: %s, err: %w", s.table, err)
	}
	updateCount, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to get updated document count for collection: %q, err: %w", s.table, err)
	}
	return &UpdateResult{UpdateCount: updateCount}, nil
}

func (s *Store[T, Q]) updateStatement(ctx context.Context, buf []byte, merge bool, m []byte) *sq.UpdateBuilder {
	st := sq.Update(s.table).
		Set(ColRev, sq.Expr(ColRev+" + 1 ")).
		Set(ColMeta, sq.Expr(ColMeta+" || ? ", m))
	if merge {
		st = st.Set(ColDoc, sq.Expr(ColDoc+" || ? ", buf))
	} else {
		st = st.Set(ColDoc, buf)
	}
	return st
}
