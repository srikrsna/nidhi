package nidhi

import (
	"context"
	"fmt"
)

// OnCreateHook is the signature of [Hooks.OnCreate] hook.
type OnCreateHook func(*HookContext, any, *CreateOptions)

// CreateOptions are options for the `Create` operation.
type CreateOptions struct {
	CreateMetadata Metadata
	// Replace is only used if the document with the id already exists.
	// In that scenario,
	// If set to true, will replace the document
	// If set to false, will return an error (default)
	Replace bool
	// ReplaceMetadata is the metadata used for the replace operation.
	// This is only use if Replace is set to true and document already exists.
	ReplaceMetadata Metadata
}

// CreateResult is the result of the create call.
// It doesn't have any fields as of now.
//
// Having an explicit result type will not break future code.
type CreateResult struct{}

// Create creates a new document.
//
// See `CreateOptions` for replacing documents if already present.
func (s *Store[T]) Create(ctx context.Context, doc *T, opts CreateOptions) (*CreateResult, error) {
	hookCtx := NewHookContext(ctx, s)
	for _, hook := range s.hooks {
		if hook.OnCreate != nil {
			hook.OnCreate(hookCtx, doc, &opts)
		}
	}
	id := s.idFn(doc)
	if id == "" {
		return nil, fmt.Errorf("nidhi: id cannot be empty")
	}
	docJSON, err := getJson(doc)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to marshal document of collection: %s, err: %w", s.table, err)
	}
	defer putJson(docJSON)
	cmdJSON, err := getJson(opts.CreateMetadata)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to marshal metadata of collection: %s, err: %w", s.table, err)
	}
	defer putJson(cmdJSON)
	args := []any{id, docJSON.Buffer(), cmdJSON.Buffer()}
	query := `INSERT INTO ` + s.table + ` (id, revision, document, metadata) VALUES ($1, 1, $2, $3)`
	if opts.Replace {
		query += ` ON CONFLICT(id) DO UPDATE SET revision = ` + s.table + `.revision + 1, document = $2, metadata = ` + s.table + `.metadata || $4`
		rmdJSON, err := getJson(opts.ReplaceMetadata)
		if err != nil {
			return nil, fmt.Errorf("nidhi: failed to marshal replace metadata of collection: %s, err: %w", s.table, err)
		}
		defer putJson(rmdJSON)
		args = append(args, rmdJSON.Buffer())
	}
	if _, err := s.db.ExecContext(ctx,
		query,
		args...,
	); err != nil {
		return nil, fmt.Errorf("nidhi: failed to create a new document: %w", err)
	}
	return &CreateResult{}, nil
}
