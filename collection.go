package nidhi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/elgris/sqrl"
	"github.com/elgris/sqrl/pg"
	jsoniter "github.com/json-iterator/go"
)

const (
	idCol   = "id"
	docCol  = "document"
	revCol  = `"revision"`
	metaCol = "metadata"

	notDeleted = "NOT (metadata ?? 'deleted')"
)

type collection struct {
	table string
	tx    sq.BaseRunner

	subFunc SubjectFunc
	fields  []string
}

func (c *collection) Create(ctx context.Context, doc Document, ops []CreateOption) (string, error) {
	var cop CreateOptions
	for _, op := range ops {
		op(&cop)
	}

	id := doc.DocumentId()
	if id == "" {
		return "", fmt.Errorf("nidhi: id cannot be empty")
	}

	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	if err := doc.MarshalDocument(stream); err != nil {
		return "", fmt.Errorf("nidhi: unable to marshal document of collection: %s, err: %w", c.table, err)
	}
	al := c.activityLog(ctx)

	args := []interface{}{
		id,
		stream.Buffer(),
		JSONB(&Metadata{Created: al}),
	}

	query := `INSERT INTO ` + c.table + ` (id, revision, document, metadata) VALUES ($1, 1, $2, $3)`
	if cop.Replace {
		query += ` ON CONFLICT(id) DO UPDATE SET revision = ` + c.table + `.revision + 1, document = $2, metadata = ` + c.table + `.metadata || $4`
		args = append(args, JSONB(&Metadata{Updated: al}))
	}

	if _, err := c.tx.ExecContext(ctx,
		query,
		args...,
	); err != nil {
		return "", fmt.Errorf("nidhi: unable to create a new document: %w", err)
	}

	return id, nil
}

func (c *collection) Replace(ctx context.Context, doc Document, ops []ReplaceOption) error {
	var rop ReplaceOptions
	for _, op := range ops {
		op(&rop)
	}

	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	if err := doc.MarshalDocument(stream); err != nil {
		return fmt.Errorf("nidhi: unable to marshal document of collection: %s, err: %w", c.table, err)
	}

	stmt := c.updateStatement(ctx, stream.Buffer(), false)
	if rop.Revision > 0 {
		stmt = stmt.Where(sq.Eq{revCol: rop.Revision})
	}

	stmt = stmt.Where(sq.Eq{idCol: doc.DocumentId()}).Where(notDeleted)

	rc, err := sq.RowsAffected(stmt.PlaceholderFormat(sq.Dollar).RunWith(c.tx).ExecContext(ctx))
	if err != nil {
		return fmt.Errorf("nidhi: unable to put document in database: %w", err)
	}

	if rc != 1 {
		return fmt.Errorf("nidhi: no document matched with the given id and revision")
	}

	return nil
}

func (c *collection) Update(ctx context.Context, doc Document, f Filter, ops []UpdateOption) error {
	var uop UpdateOptions
	for _, op := range ops {
		op(&uop)
	}

	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	doc.SetDocumentId("")
	if err := doc.MarshalDocument(stream); err != nil {
		return fmt.Errorf("nidhi: unable to marshal document of collection: %s, err: %w", c.table, err)
	}

	st := c.updateStatement(ctx, stream.Buffer(), true)

	if f != nil {
		cond, err := f.ToSql(docCol)
		if err != nil {
			return fmt.Errorf("nidhi: invalid filter received for collection: %s, err: %w", c.table, err)
		}

		st = st.Where(cond)
	}

	st = st.Where(notDeleted)

	if _, err := st.PlaceholderFormat(sq.Dollar).RunWith(c.tx).ExecContext(ctx); err != nil {
		return fmt.Errorf("nidhi: unable to update documents for collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) Delete(ctx context.Context, id string, ops []DeleteOption) error {
	var dop DeleteOptions
	for _, op := range ops {
		op(&dop)
	}

	var (
		sql  string
		args []interface{}
		err  error
	)
	if dop.Permanent {
		sql, args, err = sq.Delete(c.table).Where(sq.Eq{idCol: id}).PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	} else {
		st := sq.Update(c.table).
			Where(sq.Eq{idCol: id}).
			Where(notDeleted).
			Set(metaCol, merge(metaCol, &Metadata{Deleted: c.activityLog(ctx)}))

		sql, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	}

	if _, err := c.tx.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("nidhi: unable to delete a document of collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) DeleteMany(ctx context.Context, f Filter, ops []DeleteOption) error {
	var dop DeleteOptions
	for _, op := range ops {
		op(&dop)
	}

	var (
		sql  string
		args []interface{}
		err  error
	)
	if dop.Permanent {
		st := sq.Delete(c.table)
		if f != nil {
			cond, err := f.ToSql(docCol)
			if err != nil {
				return fmt.Errorf("nidhi: invalid filter received for collection: %s, err: %w", c.table, err)
			}

			st = st.Where(cond)
		}
		sql, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	} else {
		st := sq.Update(c.table).Set(metaCol, merge(metaCol, &Metadata{Deleted: c.activityLog(ctx)}))
		if f != nil {
			cond, err := f.ToSql(docCol)
			if err != nil {
				return fmt.Errorf("nidhi: invalid filter received for collection: %s, err: %w", c.table, err)
			}

			st = st.Where(cond).Where(notDeleted)
		}
		sql, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	}

	if _, err := c.tx.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("nidhi: unable to delete a document of collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) Query(ctx context.Context, f Filter, ctr func() Document, ops []QueryOption) error {
	var qop QueryOptions
	for _, op := range ops {
		op(&qop)
	}

	var selection interface{} = docCol
	if len(c.fields) > 0 && len(qop.ViewMask) > 0 {
		selection = sq.Expr(docCol+" - ?::text[]", pg.Array(difference(c.fields, qop.ViewMask)))
	}
	st := sq.Select().Column(selection).From(c.table)

	if f != nil {
		cond, err := f.ToSql(docCol)
		if err != nil {
			return fmt.Errorf("nidhi: invalid filter received for collection: %s, err: %w", c.table, err)
		}

		st = st.Where(cond).Where(notDeleted)
	}

	if qop.PaginationOptions != nil {
		if qop.PaginationOptions.Backward {
			if qop.PaginationOptions.Cursor != "" {
				st = st.Where(idCol+" < ?", qop.PaginationOptions.Cursor)
			}
			st = st.OrderBy(idCol + ` DESC`)
		} else {
			if qop.PaginationOptions.Cursor != "" {
				st = st.Where(idCol+" > ?", qop.PaginationOptions.Cursor)
			}
			st = st.OrderBy(idCol + ` ASC`)
		}

		qop.PaginationOptions.HasMore = false
		st = st.Limit(qop.PaginationOptions.Limit + 1)
	}

	rows, err := st.PlaceholderFormat(sq.Dollar).RunWith(c.tx).QueryContext(ctx)
	if err != nil {
		return fmt.Errorf("nidhi: unable to query collection: %s, err: %w", c.table, err)
	}
	defer rows.Close()

	iter := jsoniter.ConfigDefault.BorrowIterator(nil)
	defer jsoniter.ConfigDefault.ReturnIterator(iter)

	var (
		count  uint64
		entity sql.RawBytes
	)
	for rows.Next() {
		if qop.PaginationOptions != nil && qop.PaginationOptions.Limit <= count {
			qop.PaginationOptions.HasMore = true
			break
		}

		count++
		doc := ctr()
		if err := rows.Scan(&entity); err != nil {
			return fmt.Errorf("nidhi: unexpected error while querying collection: %s, err: %w", c.table, err)
		}

		iter.ResetBytes(entity)
		if err := doc.UnmarshalDocument(iter); err != nil {
			return fmt.Errorf("nidhi: unable to unmarshal document of type %s, err: %w", c.table, err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("nidhi: unexpected error while querying collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) Get(ctx context.Context, id string, doc Document, ops []GetOption) error {
	var gop GetOptions
	for _, op := range ops {
		op(&gop)
	}

	var selection interface{} = docCol
	if len(c.fields) > 0 && len(gop.ViewMask) > 0 {
		selection = sq.Expr(docCol+" - ?::text[]", pg.Array(difference(c.fields, gop.ViewMask)))
	}

	st := sq.Select().Column(selection).From(c.table).Where(sq.Eq{idCol: id}).Where(notDeleted)

	var entity []byte
	if err := st.PlaceholderFormat(sq.Dollar).RunWith(c.tx).QueryRowContext(ctx).Scan(&entity); err != nil {
		return fmt.Errorf("nidhi: unable to get a document from collection %q, err: %w", c.table, err)
	}

	iter := jsoniter.ConfigDefault.BorrowIterator(entity)
	defer jsoniter.ConfigDefault.ReturnIterator(iter)

	if err := doc.UnmarshalDocument(iter); err != nil {
		return fmt.Errorf("nidhi: unable to unmarshal document of type %s, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) Count(ctx context.Context, f Filter, ops []CountOption) (int64, error) {
	var qop CountOptions
	for _, op := range ops {
		op(&qop)
	}

	st := sq.Select("count(*)").From(c.table)

	if f != nil {
		cond, err := f.ToSql(docCol)
		if err != nil {
			return 0, fmt.Errorf("nidhi: invalid filter received for collection: %s, err: %w", c.table, err)
		}

		st = st.Where(cond)
	}

	var count int64
	if err := st.PlaceholderFormat(sq.Dollar).RunWith(c.tx).QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("nidhi: unable to query collection: %s, err: %w", c.table, err)
	}

	return count, nil
}

func (c *collection) activityLog(ctx context.Context) *ActivityLog {
	sub := ""
	if c.subFunc != nil {
		sub = c.subFunc(ctx)
	}
	return &ActivityLog{
		By: sub,
		On: time.Now(),
	}
}

func (c *collection) updateStatement(ctx context.Context, buf []byte, merge bool) *sq.UpdateBuilder {
	st := sq.Update(c.table).
		Set(revCol, sq.Expr(revCol+" + 1 ")).
		Set(metaCol, sq.Expr(metaCol+" || ? ", JSONB(&Metadata{Updated: c.activityLog(ctx)})))

	if merge {
		st = st.Set(docCol, sq.Expr(docCol+" || ? ", buf))
	} else {
		st = st.Set(docCol, buf)
	}

	return st
}

func merge(column string, value interface {
	Marshaler
	Unmarshaler
}) sq.Sqlizer {
	return sq.Expr(column+" || ? ", JSONB(value))
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

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
