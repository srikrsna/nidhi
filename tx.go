package nidhi

import (
	"context"
	"database/sql"
	"fmt"
)

type Collection struct {
	*collection

	db *sql.DB
}

func OpenCollection(ctx context.Context, db *sql.DB, schema, name string, opts CollectionOptions) (*Collection, error) {
	if _, err := db.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS `+schema); err != nil {
		return nil, fmt.Errorf("nidhi: unable to open collection: %s, err: %w", name, err)
	}
	const query = `CREATE TABLE IF NOT EXISTS %s.%s (id TEXT NOT NULL PRIMARY KEY, revision bigint NOT NULL, document JSONB NOT NULL, metadata JSONB NOT NULL DEFAULT '{}')`
	if _, err := db.ExecContext(ctx, fmt.Sprintf(query, schema, name)); err != nil {
		return nil, fmt.Errorf("nidhi: unable to open collection: %s, err: %w", name, err)
	}

	return &Collection{&collection{table: schema + "." + name, tx: db, subFunc: opts.SubjectFunc, fields: opts.Fields}, db}, nil
}

func (c *Collection) BeginTx(ctx context.Context, opt *sql.TxOptions) (*TxCollection, error) {
	tx, err := c.db.BeginTx(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("nidhi: unable to start a transaction: %w", err)
	}

	nc := *c.collection
	nc.tx = tx
	return &TxCollection{&nc, tx}, nil
}

func (c *Collection) WithTransaction(tcol *TxToken) *TxCollection {
	nc := *c.collection
	nc.tx = tcol.tx.collection.tx
	return &TxCollection{&nc, tcol.tx}
}

type TxCollection struct {
	*collection
	rollBackComitter
}

type rollBackComitter interface {
	Rollback() error
	Commit() error
}

type TxToken struct {
	tx *TxCollection
}

func NewTxToken(tx *TxCollection) *TxToken {
	return &TxToken{tx}
}
