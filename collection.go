package nidhi

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/elgris/sqrl"
	"github.com/elgris/sqrl/pg"
	jsoniter "github.com/json-iterator/go"
)

const (
	ColId   = "id"
	ColDoc  = "document"
	ColRev  = `"revision"`
	ColMeta = "metadata"
	ColDel  = "deleted"

	notDeleted = ColDel + " = false "
	seperator  = ":"
)

type collection struct {
	table string
	tx    sq.BaseRunner

	fields []string
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

	args := []interface{}{
		id,
		stream.Buffer(),
		JSONB(mdMarshaler(cop.CreateMetadata)),
	}

	query := `INSERT INTO ` + c.table + ` (id, revision, document, metadata) VALUES ($1, 1, $2, $3)`
	if cop.Replace {
		query += ` ON CONFLICT(id) DO UPDATE SET revision = ` + c.table + `.revision + 1, document = $2, metadata = ` + c.table + `.metadata || $4`
		args = append(args, JSONB(mdMarshaler(cop.ReplaceMetadata)))
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

	stmt := c.updateStatement(ctx, stream.Buffer(), false, rop.Metadata)
	if rop.Revision > 0 {
		stmt = stmt.Where(sq.Eq{ColRev: rop.Revision})
	}

	stmt = stmt.Where(sq.Eq{ColId: doc.DocumentId()}).Where(notDeleted)

	rc, err := sq.RowsAffected(stmt.PlaceholderFormat(sq.Dollar).RunWith(c.tx).ExecContext(ctx))
	if err != nil {
		return fmt.Errorf("nidhi: unable to put document in database: %w", err)
	}

	if rc != 1 {
		return fmt.Errorf("nidhi: no document matched with the given id and revision")
	}

	return nil
}

func (c *collection) Update(ctx context.Context, doc Document, f Sqlizer, ops []UpdateOption) error {
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

	st := c.updateStatement(ctx, stream.Buffer(), true, uop.Metadata)

	if f != nil {
		st = st.Where(f)
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
		sqlStr string
		args   []interface{}
		err    error
	)
	if dop.Permanent {
		sqlStr, args, err = sq.Delete(c.table).Where(sq.Eq{ColId: id}).PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	} else {
		st := sq.Update(c.table).
			Where(sq.Eq{ColId: id}).
			Where(notDeleted).
			Set(ColDel, true).
			Set(ColMeta, merge(ColMeta, mdMarshaler(dop.Metadata)))

		sqlStr, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	}

	if _, err := c.tx.ExecContext(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("nidhi: unable to delete a document of collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) DeleteMany(ctx context.Context, f Sqlizer, ops []DeleteOption) error {
	var dop DeleteOptions
	for _, op := range ops {
		op(&dop)
	}

	var (
		sqlStr string
		args   []interface{}
		err    error
	)
	if dop.Permanent {
		st := sq.Delete(c.table)
		if f != nil {
			st = st.Where(f)
		}
		sqlStr, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	} else {
		st := sq.Update(c.table).Set(ColDel, true).Set(ColMeta, merge(ColMeta, mdMarshaler(dop.Metadata)))
		if f != nil {
			st = st.Where(f)
		}
		st = st.Where(notDeleted)

		sqlStr, args, err = st.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("nidhi: there seems to be a bug: unable to build delete statement: %w", err)
		}
	}

	if _, err := c.tx.ExecContext(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("nidhi: unable to delete a document of collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) Query(ctx context.Context, f Sqlizer, ctr func() Document, ops []QueryOption) error {
	var qop QueryOptions
	for _, op := range ops {
		op(&qop)
	}

	var selection interface{} = ColDoc
	if len(c.fields) > 0 && len(qop.ViewMask) > 0 {
		selection = sq.Expr(ColDoc+" - ?::text[]", pg.Array(difference(c.fields, qop.ViewMask)))
	}

	st := sq.Select().Column(selection).From(c.table)
	if f != nil {
		st = st.Where(f)
	}

	st = st.Where(notDeleted)

	st, scans, err := addPagination(st, qop.PaginationOptions)
	if err != nil {
		return fmt.Errorf("nidhi: unable to paginate: %s, err: %w", c.table, err)
	}
	var md sql.RawBytes
	if len(qop.CreateMetadata) > 0 {
		st = st.Column(ColMeta)
		scans = append(scans, &md)
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
	scans = append([]interface{}{&entity}, scans...)
	for rows.Next() {
		if qop.PaginationOptions != nil && qop.PaginationOptions.Limit <= count {
			qop.PaginationOptions.HasMore = true
			break
		}

		count++
		doc := ctr()
		if err := rows.Scan(scans...); err != nil {
			return fmt.Errorf("nidhi: unexpected error while querying collection: %s, err: %w", c.table, err)
		}

		iter.ResetBytes(entity)
		if err := doc.UnmarshalDocument(iter); err != nil {
			return fmt.Errorf("nidhi: unable to unmarshal document of type %s, err: %w", c.table, err)
		}

		if len(qop.CreateMetadata) > 0 {
			iter.ResetBytes(md)
			mdu := make(mdUnmarshaler, 0, len(qop.CreateMetadata))
			for _, cmd := range qop.CreateMetadata {
				mdu = append(mdu, cmd())
			}
			if err := mdu.UnmarshalDocument(iter); err != nil {
				return fmt.Errorf("nidhi: unable to unmarshal metadata of type %s, err: %w", c.table, err)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("nidhi: unexpected error while querying collection: %s, err: %w", c.table, err)
	}

	if qop.PaginationOptions != nil {
		offset := 0
		if len(qop.CreateMetadata) > 0 {
			offset = 1
		}
		if len(scans) == 2+offset {
			qop.PaginationOptions.NextCursor = *(scans[1].(*string))
		} else {
			qop.PaginationOptions.NextCursor = qop.PaginationOptions.OrderBy[0].Field.Encode(scans[1], *(scans[2].(*string)))
		}
	}

	return nil
}

func (c *collection) Get(ctx context.Context, id string, doc Unmarshaler, ops []GetOption) error {
	var gop GetOptions
	for _, op := range ops {
		op(&gop)
	}

	var selection interface{} = ColDoc
	if len(c.fields) > 0 && len(gop.ViewMask) > 0 {
		selection = sq.Expr(ColDoc+" - ?::text[]", pg.Array(difference(c.fields, gop.ViewMask)))
	}

	st := sq.Select().Column(selection)

	scans := make([]interface{}, 0, 2)
	scans = append(scans, JSONB(NoopMarshaler{Unmarshaler: doc}))

	if len(gop.Metadata) > 0 {
		st = st.Column(ColMeta)
		scans = append(scans, JSONB(mdUnmarshaler(gop.Metadata)))
	}

	st = st.From(c.table).Where(sq.Eq{ColId: id}).Where(notDeleted)

	if err := st.PlaceholderFormat(sq.Dollar).RunWith(c.tx).QueryRowContext(ctx).Scan(scans...); err != nil {
		return fmt.Errorf("nidhi: unable to get a document from collection %q, err: %w", c.table, err)
	}

	return nil
}

func (c *collection) updateStatement(ctx context.Context, buf []byte, merge bool, m mdMarshaler) *sq.UpdateBuilder {
	st := sq.Update(c.table).
		Set(ColRev, sq.Expr(ColRev+" + 1 ")).
		Set(ColMeta, sq.Expr(ColMeta+" || ? ", JSONB(m)))

	if merge {
		st = st.Set(ColDoc, sq.Expr(ColDoc+" || ? ", buf))
	} else {
		st = st.Set(ColDoc, buf)
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
