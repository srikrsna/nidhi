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

// TODO: Add Pagination
// TODO: Parent
// TODO: Generation
// TODO: Partial Update

type Collection struct {
	table string
	db    *sql.DB

	subFunc SubjectFunc
	fields  []string
}

func OpenCollection(ctx context.Context, db *sql.DB, schema, name string, opts CollectionOptions) (*Collection, error) {
	if _, err := db.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS `+schema); err != nil {
		return nil, fmt.Errorf("nidhi: unable to open collection: %s, err: %w", name, err)
	}
	const query = `CREATE TABLE IF NOT EXISTS %s.%s (id TEXT NOT NULL PRIMARY KEY, revision bigint NOT NULL, document JSONB NOT NULL, metadata JSONB NOT NULL DEFAULT '{}')`
	if _, err := db.ExecContext(ctx, fmt.Sprintf(query, schema, name)); err != nil {
		return nil, fmt.Errorf("nidhi: unable to open collection: %s, err: %w", name, err)
	}

	return &Collection{table: schema + "." + name, db: db, subFunc: opts.SubjectFunc, fields: opts.Fields}, nil
}

func (c *Collection) Create(ctx context.Context, doc Document, ops []CreateOption) (string, error) {
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

	if _, err := c.db.ExecContext(ctx,
		query,
		args...,
	); err != nil {
		return "", fmt.Errorf("nidhi: unable to create a new document: %w", err)
	}

	return id, nil
}

func (c *Collection) Replace(ctx context.Context, doc Document, ops []ReplaceOption) error {
	var rop ReplaceOptions
	for _, op := range ops {
		op(&rop)
	}

	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	if err := doc.MarshalDocument(stream); err != nil {
		return fmt.Errorf("nidhi: unable to marshal document of collection: %s, err: %w", c.table, err)
	}

	stmt := sq.Update(c.table).
		Set(docCol, stream.Buffer()).
		Set(revCol, sq.Expr(revCol+" + 1 ")).
		Set(metaCol, sq.Expr(metaCol+" || "+" ? ", JSONB(&Metadata{Updated: c.activityLog(ctx)})))

	if rop.Revision > 0 {
		stmt = stmt.Where(sq.Eq{revCol: rop.Revision})
	}

	stmt = stmt.Where(sq.Eq{idCol: doc.DocumentId()}).Where(notDeleted)

	rc, err := sq.RowsAffected(stmt.PlaceholderFormat(sq.Dollar).RunWith(c.db).ExecContext(ctx))
	if err != nil {
		return fmt.Errorf("nidhi: unable to put document in database: %w", err)
	}

	if rc != 1 {
		return fmt.Errorf("nidhi: no document matched with the given id and revision")
	}

	return nil
}

func (c *Collection) Delete(ctx context.Context, id string, ops []DeleteOption) error {
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

	if _, err := c.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("nidhi: unable to delete a document of collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *Collection) DeleteMany(ctx context.Context, f Filter, ops []DeleteOption) error {
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

	if _, err := c.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("nidhi: unable to delete a document of collection: %s, err: %w", c.table, err)
	}

	return nil
}

func (c *Collection) Query(ctx context.Context, f Filter, ctr func() Document, ops []QueryOption) error {
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

	rows, err := st.PlaceholderFormat(sq.Dollar).RunWith(c.db).QueryContext(ctx)
	if err != nil {
		return fmt.Errorf("nidhi: unable to query collection: %s, err: %w", c.table, err)
	}
	defer rows.Close()

	iter := jsoniter.ConfigDefault.BorrowIterator(nil)
	defer jsoniter.ConfigDefault.ReturnIterator(iter)

	var entity = sql.RawBytes{}
	for rows.Next() {
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

func (c *Collection) Get(ctx context.Context, id string, doc Document, ops []GetOption) error {
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
	if err := st.PlaceholderFormat(sq.Dollar).RunWith(c.db).QueryRowContext(ctx).Scan(&entity); err != nil {
		return fmt.Errorf("nidhi: unable to get a document from collection %q, err: %w", c.table, err)
	}

	iter := jsoniter.ConfigDefault.BorrowIterator(entity)
	defer jsoniter.ConfigDefault.ReturnIterator(iter)

	if err := doc.UnmarshalDocument(iter); err != nil {
		return fmt.Errorf("nidhi: unable to unmarshal document of type %s, err: %w", c.table, err)
	}

	return nil
}

func (c *Collection) Count(ctx context.Context, f Filter, ops []CountOption) (int64, error) {
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
	if err := st.PlaceholderFormat(sq.Dollar).RunWith(c.db).QueryRowContext(ctx).Scan(&count); err != nil {
		return 0, fmt.Errorf("nidhi: unable to query collection: %s, err: %w", c.table, err)
	}

	return count, nil
}

func (c *Collection) activityLog(ctx context.Context) *ActivityLog {
	sub := ""
	if c.subFunc != nil {
		sub = c.subFunc(ctx)
	}
	return &ActivityLog{
		By: sub,
		On: time.Now(),
	}
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
