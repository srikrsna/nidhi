package nidhi

import (
	"context"
	"fmt"

	sq "github.com/elgris/sqrl"
)

// DeleteOptions are options for the `Delete` operation
type DeleteOptions struct {
	// Metadata is the metadata of the document.
	// This will be merged with existing metadata.
	//
	// NOTE: This is a no op if `Permanent` is set.
	Metadata Metadata
	// Permanent if set will hard delete the document.
	Permanent bool
}

// DeleteManyOptions are options for the `DeleteMany` operation
type DeleteManyOptions struct {
	// Metadata is the metadata of the document.
	// This will be merged with existing metadata.
	//
	// NOTE: This is a no op if `Permanent` is set.
	Metadata Metadata
	// Permanent if set will hard delete the document.
	Permanent bool
}

// DeleteResult is the result of the delete call.
// It doesn't have any fields as of now.
//
// Having an explicit result type will not break future changes.
type DeleteResult struct{}

// DeleteManyResult is the result of the delete many call.
type DeleteManyResult struct {
	// DeleteCount is the number of documents that were deleted.
	DeleteCount int64
}

// Delete deletes a single record from the store using its id.
//
// By default all deletes are soft deletes. To hard delete, see `DeleteOptions`
func (s *Store[T, Q]) Delete(ctx context.Context, id string, opts DeleteOptions) (*DeleteResult, error) {
	var (
		sqlStr string
		args   []interface{}
		err    error
	)
	if opts.Permanent {
		sqlStr, args, err = sq.Delete(s.table).Where(sq.Eq{ColId: id}).PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return nil, fmt.Errorf("nidhi: this seems to be a bug: failed to build hard delete statement: %w", err)
		}
	} else {
		mdJSON, err := getJson(opts.Metadata)
		if err != nil {
			return nil, fmt.Errorf("nidhi: failed to marshal metadata of collection: %s, err: %w", s.table, err)
		}
		defer putJson(mdJSON)
		st := sq.Update(s.table).
			Where(sq.Eq{ColId: id}).
			Where(notDeleted).
			Set(ColDel, true).
			Set(ColMeta, merge(ColMeta, mdJSON.Buffer()))

		sqlStr, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return nil, fmt.Errorf("nidhi: this seems to be a bug: failed to build soft delete statement: %w", err)
		}
	}
	if _, err := s.db.ExecContext(ctx, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("nidhi: failed to delete a document of collection: %s, err: %w", s.table, err)
	}
	return &DeleteResult{}, nil
}

func (s *Store[T, Q]) DeleteMany(ctx context.Context, q Q, opts DeleteManyOptions) (*DeleteManyResult, error) {
	var (
		sqlStr string
		args   []any
		err    error
	)
	if opts.Permanent {
		st := sq.Delete(s.table)
		if any(q) != nil {
			st = st.Where(q)
		}
		sqlStr, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return nil, fmt.Errorf("nidhi: there seems to be a bug: failed to build delete statement: %w", err)
		}
	} else {
		mdJSON, err := getJson(opts.Metadata)
		if err != nil {
			return nil, fmt.Errorf("nidhi: failed to marshal metadata of collection: %s, err: %w", s.table, err)
		}
		defer putJson(mdJSON)
		st := sq.Update(s.table).Set(ColDel, true).Set(ColMeta, merge(ColMeta, mdJSON.Buffer()))
		if any(q) != nil {
			st = st.Where(q)
		}
		st = st.Where(notDeleted)
		sqlStr, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return nil, fmt.Errorf("nidhi: there seems to be a bug: failed to build delete statement: %w", err)
		}
	}
	res, err := s.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to delete a document of collection: %s, err: %w", s.table, err)
	}
	deleteCount, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("nidhi: failed to get deleted document count for collection: %q, err: %w", s.table, err)
	}
	return &DeleteManyResult{
		DeleteCount: deleteCount,
	}, nil
}

func merge(column string, value []byte) sq.Sqlizer {
	return sq.Expr(column+" || ? ", value)
}
