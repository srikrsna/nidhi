package nidhi

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
)

const (
	schemaTemplate = `
	CREATE SCHEMA IF NOT EXISTS %s
	`
	tableTemplate = `
	CREATE TABLE IF NOT EXISTS %s.%s (
		id TEXT NOT NULL PRIMARY KEY,
		revision BIGINT NOT NULL,
		document JSONB NOT NULL,
		deleted BOOLEAN NOT NULL DEFAULT FALSE,
		metadata JSONB NOT NULL DEFAULT '{}'
	)
	`

	ColId   = "id"
	ColDoc  = "document"
	ColRev  = `"revision"`
	ColMeta = "metadata"
	ColDel  = "deleted"

	notDeleted = ColDel + " = false "
)

type (
	// IdFn should return the unique id of a document.
	IdFn[T any] func(*T) string
	// SetIdFn should set the unique id of a document.
	SetIdFn[T any] func(*T, string)
)

// StoreOptions are options for a Store.
type StoreOptions struct {
	// MetadataRegistry is the registry of metadata parts.
	MetadataRegistry map[string]func() MetadataPart
}

// Store is the collection of documents
type Store[T any, Q Sqlizer] struct {
	db *sql.DB

	table  string
	fields []string

	idFn    IdFn[T]
	setIdFn SetIdFn[T]

	mdr map[string]func() MetadataPart
}

// NewStore returns a new store.
//
// Typically this is never direclty called. It is called via a more concrete generated function.
// See protoc-gen-nidhi.
func NewStore[T any, Q Sqlizer](
	ctx context.Context,
	db *sql.DB,
	schema, table string,
	fields []string,
	idFn IdFn[T],
	setIdFn SetIdFn[T],
	opts StoreOptions,
) (*Store[T, Q], error) {
	if _, err := db.ExecContext(ctx, fmt.Sprintf(schemaTemplate, schema)); err != nil {
		return nil, fmt.Errorf("nidhi: failed to create schema: %q, err: %w", schema, err)
	}
	if _, err := db.ExecContext(ctx, fmt.Sprintf(tableTemplate, schema, table)); err != nil {
		return nil, fmt.Errorf("nidhi: failed to create table: \"%s.%s\", err: %w", schema, table, err)
	}
	return &Store[T, Q]{
		db:      db,
		table:   schema + "." + table,
		fields:  fields,
		idFn:    idFn,
		setIdFn: setIdFn,
		mdr:     opts.MetadataRegistry,
	}, nil
}
