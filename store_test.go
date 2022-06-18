package nidhi_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/akshayjshah/attest"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/srikrsna/nidhi"
)

const (
	schema = "resource"
	table  = "resources"
)

var (
	fields = []string{
		"id",
		"title",
		"dateOfBirth",
		"age",
		"canDrive",
	}
)

type resource struct {
	Id          string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	DateOfBirth time.Time `json:"dateOfBirth,omitempty"`
	Age         int       `json:"age,omitempty"`
	CanDrive    bool      `json:"canDrive,omitempty"`
}

func TestNewStore(t *testing.T) {
	t.Parallel()
	db := newDB(t)
	store := newStore(t, db, nidhi.StoreOptions{})
	attest.NotZero(t, store)
	// Check if schema and table were created.
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2)`, schema, table).Scan(&exists)
	attest.Ok(t, err)
	attest.True(t, exists)
	// Check columns of created table.
	rows, err := db.Query(`SELECT column_name, column_default, is_nullable, data_type FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2`, schema, table)
	attest.Ok(t, err)
	t.Cleanup(func() { attest.Ok(t, rows.Close()) })
	type column struct {
		name         string
		defaultValue sql.NullString
		nullable     string
		datatype     string
	}
	var cc []*column
	for rows.Next() {
		var c column
		attest.Ok(t, rows.Scan(&c.name, &c.defaultValue, &c.nullable, &c.datatype))
		c.datatype = strings.ToLower(c.datatype)
		c.nullable = strings.ToLower(c.nullable)
		cc = append(cc, &c)
	}
	expectedColumns := []*column{
		{"id", sql.NullString{}, "no", "text"},
		{"document", sql.NullString{}, "no", "jsonb"},
		{"metadata", sql.NullString{String: "'{}'::jsonb", Valid: true}, "no", "jsonb"},
		{"revision", sql.NullString{}, "no", "bigint"},
		{"deleted", sql.NullString{String: "false", Valid: true}, "no", "boolean"},
	}
	for _, c := range cc {
		attest.Contains(t, expectedColumns, c, attest.Allow(*c), attest.Continue())
	}
	// Check if store doesn't error on an existing table.
	store = newStore(t, db, nidhi.StoreOptions{})
	attest.NotZero(t, store)
}

func newStore(tb testing.TB, db *sql.DB, opts nidhi.StoreOptions) *nidhi.Store[resource, nidhi.Sqlizer] {
	store, err := nidhi.NewStore[resource, nidhi.Sqlizer](
		context.Background(),
		db,
		schema,
		table,
		fields,
		func(r *resource) string { return r.Id },
		func(r *resource, id string) { r.Id = id },
		opts,
	)
	attest.Ok(tb, err)
	return store
}

func newDB(tb testing.TB) *sql.DB {
	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("myuser", "mypass"),
		Path:   "mydatabase",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()
	pool, err := dockertest.NewPool("")
	attest.Ok(tb, err)
	pw, _ := pgURL.User.Password()
	runOptions := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15beta1-alpine",
		Env: []string{
			"POSTGRES_USER=" + pgURL.User.Username(),
			"POSTGRES_PASSWORD=" + pw,
			"POSTGRES_DB=" + pgURL.Path,
		},
		Labels: map[string]string{"postgrestesting": "1"},
	}
	resource, err := pool.RunWithOptions(
		runOptions,
		func(config *docker.HostConfig) {
			// Set AutoRemove to true so that stopped container goes away by itself.
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	attest.Ok(tb, err)
	tb.Cleanup(func() {
		err = pool.Purge(resource)
		attest.Ok(tb, err)
	})
	pgURL.Host = resource.Container.NetworkSettings.IPAddress
	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}
	pool.MaxWait = 5 * time.Minute
	err = pool.Retry(func() (err error) {
		db, err := sql.Open("pgx", pgURL.String())
		if err != nil {
			return err
		}
		defer func() {
			closeErr := db.Close()
			if err == nil {
				err = closeErr
			}
		}()
		return db.Ping()
	})
	attest.Ok(tb, err)
	db, err := sql.Open("pgx", pgURL.String())
	attest.Ok(tb, err)
	tb.Cleanup(func() {
		attest.Ok(tb, db.Close())
	})
	attest.Ok(tb, db.Ping())
	return db
}

func loadDoc(t testing.TB, db *sql.DB, id string, md nidhi.Metadata) *nidhi.Document[resource] {
	var (
		docJsonData, mdJsonData []byte
		revision                int64
		deleted                 bool
	)
	attest.Ok(t,
		db.QueryRow(
			fmt.Sprintf(
				`SELECT document, revision, metadata, deleted FROM %s.%s WHERE id = $1`,
				schema,
				table,
			),
			id,
		).Scan(&docJsonData, &revision, &mdJsonData, &deleted),
	)
	var er resource
	attest.Ok(t, json.Unmarshal(docJsonData, &er))
	if md != nil {
		attest.Ok(t, nidhi.UnmarshalJson(mdJsonData, md))
	}
	return &nidhi.Document[resource]{
		&er,
		revision,
		md,
		deleted,
	}
}