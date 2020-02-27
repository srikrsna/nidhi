package nidhi

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/appointy/idgen"
	jsoniter "github.com/json-iterator/go"
)

const (
	idCol   = "id"
	delCol  = "deleted"
	docCol  = "document"
	revCol  = "revision"
	metaCol = "metadata"
)

type Collection struct {
	table  string
	db     *sql.DB
	prefix string
}

func OpenCollection(ctx context.Context, db *sql.DB, schema, name, prefix string) (*Collection, error) {
	const query = `CREATE TABLE IF NOT EXISTS %s.%s (id TEXT NOT NULL PRIMARY KEY, revision bigint NOT NULL, deleted BOOLEAN NOT NULL DEFAULT false, document JSONB NOT NULL, metadata JSONB NOT NULL DEFAULT '{}')`
	if _, err := db.ExecContext(ctx, fmt.Sprintf(query, schema, name)); err != nil {
		return nil, fmt.Errorf("nidhi: unable to open collection: %s, err: %w", name, err)
	}

	return &Collection{table: schema + "." + name, db: db}, nil
}

func (c *Collection) Create(ctx context.Context, doc Document, ops []CreateOption) (string, error) {
	var cop CreateOptions
	for _, op := range ops {
		op(&cop)
	}
	id := idgen.New(c.prefix)
	if doc.DocumentId() == "" {
		doc.SetDocumentId(id)
	}

	stream := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(stream)

	if err := doc.MarshalDocument(stream); err != nil {
		return "", fmt.Errorf("nidhi: unable to marshal document of collection: %s, err: %w", c.table, err)
	}

	if _, err := c.db.ExecContext(ctx,
		`INSERT INTO `+c.table+` (id, revision, document) VALUES ($1, 1, $2) on conflict(id) set revision = revision + 1, document = $2, deleted = false`,
		id,
		stream.Buffer(),
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
		Set(docCol, stream.Buffer())

	if rop.Revision > 0 {
		stmt = stmt.Set(revCol, rop.Revision+1).
			Where(sq.Eq{revCol: rop.Revision})
	}

	stmt = stmt.Where(sq.Eq{idCol: doc.DocumentId(), "deleted": false})

	res, err := stmt.PlaceholderFormat(sq.Dollar).RunWith(c.db).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("nidhi: unable to put document in database: %w", err)
	}

	rc, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("nidhi: unable to determine if put was executed: %w", err)
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
		sql, args, err = sq.Update(c.table).Where(sq.Eq{idCol: id}).Set(delCol, true).PlaceholderFormat(sq.Dollar).ToSql()
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

	st := sq.Select(docCol).From(c.table)

	if f != nil {
		cond, err := f.ToSql(docCol)
		if err != nil {
			return fmt.Errorf("nidhi: invalid filter received for collection: %s, err: %w", c.table, err)
		}

		st = st.Where(cond)
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
