package nidhi

import (
	"context"
	"fmt"

	sq "github.com/elgris/sqrl"
)

// OnReplaceHook is the signature for the [Hooks.OnReplace] hook.
type OnReplaceHook func(*HookContext, any, *ReplaceOptions)

// OnUpdateHook is the signature for the [Hooks.OnUpdate] hook.
type OnUpdateHook func(*HookContext, string, any, *UpdateOptions)

// OnUpdateManyHook is the signature for the [Hooks.OnUpdateMany] hook.
type OnUpdateManyHook func(*HookContext, any, Sqlizer, *UpdateManyOptions)

// ReplaceOptions are options for `Replace` operation.
type ReplaceOptions struct {
	// Metadata is the metadata of the document.
	// This will be merged with existing metadata.
	Metadata Metadata
	// Revision if set (>0), the document will only be replaced if the revision also matches.
	Revision int64
}

// UpdateOptions are options for [*Store.Update] operation.
type UpdateOptions struct {
	// Metadata is the metadata of the document.
	// This will be merged with existing metadata.
	Metadata Metadata
	// Revision if set (>0), the document will only be replaced if the revision also matches.
	Revision int64
}

// UpdateManyOptions are options for [*Store.UpdateMany] operation.
type UpdateManyOptions struct {
	// Metadata is the metadata of the document.
	// This will be merged with existing metadata.
	Metadata Metadata
}

// ReplaceResult is the result of the replace call.
// It doesn't have any fields as of now.
//
// Having an explicit result type will not break future changes.
type ReplaceResult struct{}

// UpdateResult is the result of the [*Store.Update] call.
// It doesn't have any fields as of now.
//
// Having an explicit result type will not break future changes.
type UpdateResult struct{}

// UpdateManyResult is the result of the [*Store.UpdateMany] call.
type UpdateManyResult struct {
	// UpdateCount is the number of documents that were updated.
	UpdateCount int64
}

// Replace replaces a document, matched using it's id, in the store.
//
// Returns a NotFound error, if the document doesn't exist or the revision doesn't exist.
func (s *Store[T]) Replace(ctx context.Context, doc *T, opts ReplaceOptions) (*ReplaceResult, error) {
	hookCtx := NewHookContext(ctx, s)
	for _, h := range s.hooks {
		if h.OnReplace != nil {
			h.OnReplace(hookCtx, doc, &opts)
		}
	}
	rc, err := s.update(ctx, doc, sq.Eq{ColId: s.idFn(doc)}, false, opts.Metadata, opts.Revision)
	if err != nil {
		return nil, err
	}
	if rc != 1 {
		return nil, ErrNotFound
	}
	return &ReplaceResult{}, nil
}

// Update updates a document, matched using it's id, in the store.
//
// Returns a NotFound error on id and revision mismatch.
func (s *Store[T]) Update(ctx context.Context, id string, updates any, opts UpdateOptions) (*UpdateResult, error) {
	hookCtx := NewHookContext(ctx, s)
	for _, h := range s.hooks {
		if h.OnUpdate != nil {
			h.OnUpdate(hookCtx, id, updates, &opts)
		}
	}
	rc, err := s.update(ctx, updates, sq.Eq{ColId: id}, true, opts.Metadata, opts.Revision)
	if err != nil {
		return nil, err
	}
	if rc != 1 {
		return nil, ErrNotFound
	}
	return &UpdateResult{}, nil
}

// UpdateMany updates all the documents that match the given query.
func (s *Store[T]) UpdateMany(ctx context.Context, updates any, q Sqlizer, opts UpdateManyOptions) (*UpdateManyResult, error) {
	hookCtx := NewHookContext(ctx, s)
	for _, h := range s.hooks {
		if h.OnUpdateMany != nil {
			h.OnUpdateMany(hookCtx, updates, q, &opts)
		}
	}
	rc, err := s.update(ctx, updates, q, true, opts.Metadata, -1)
	if err != nil {
		return nil, err
	}
	return &UpdateManyResult{UpdateCount: rc}, nil
}

func (s *Store[T]) update(ctx context.Context, updates any, q Sqlizer, shouldMerge bool, md Metadata, revision int64) (int64, error) {
	updatesJSON, err := getJson(updates)
	if err != nil {
		return -1, fmt.Errorf("nidhi: failed to marshal document of collection: %s, err: %w", s.table, err)
	}
	defer putJson(updatesJSON)
	mdJSON, err := getJson(md)
	if err != nil {
		return -1, fmt.Errorf("nidhi: failed to marshal metadata of collection: %s, err: %w", s.table, err)
	}
	defer putJson(mdJSON)
	stmt := sq.Update(s.table).
		Set(ColRev, sq.Expr(ColRev+" + 1 ")).
		Set(ColMeta, merge(ColMeta, mdJSON.Buffer()))
	if shouldMerge {
		stmt = stmt.Set(ColDoc, merge(ColDoc, updatesJSON.Buffer()))
	} else {
		stmt = stmt.Set(ColDoc, updatesJSON.Buffer())
	}
	if revision > 0 {
		stmt = stmt.Where(sq.Eq{ColRev: revision})
	}
	if q != nil {
		stmt = stmt.Where(q)
	}
	stmt = stmt.Where(notDeleted)
	rc, err := sq.RowsAffected(stmt.PlaceholderFormat(sq.Dollar).RunWith(s.db).ExecContext(ctx))
	if err != nil {
		return -1, fmt.Errorf("nidhi: failed to put document in database: %w", err)
	}
	return rc, nil
}
