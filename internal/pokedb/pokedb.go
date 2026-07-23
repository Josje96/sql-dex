// Package pokedb wraps the Pokémon SQLite database and exposes a small,
// read-only query API. Keeping this separate from the CLI means Phase 2 (GUI)
// and Phase 3 (AI helper) can reuse the same engine.
package pokedb

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite" // pure-Go driver, no CGO required
)

// DB is a read-only handle to the pokédex database.
type DB struct {
	sql *sql.DB
}

// Result is the outcome of a query: the column names and the rows, where each
// value is either a string, int64, float64, []byte, or nil (for SQL NULL).
type Result struct {
	Columns []string
	Rows    [][]any
}

// Open opens the SQLite file at path in read-only mode. Read-only keeps the
// "game" data safe: learners can only SELECT, never mutate the pokédex.
func Open(path string) (*DB, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("cannot find database at %q: %w", path, err)
	}
	// mode=ro opens read-only; immutable=1 tells SQLite the file won't change,
	// which avoids lock files and is a bit faster for our use.
	dsn := fmt.Sprintf("file:%s?mode=ro&immutable=1", path)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return &DB{sql: sqlDB}, nil
}

// Close releases the underlying database handle.
func (db *DB) Close() error {
	return db.sql.Close()
}

// Query runs a SQL statement and returns its columns and rows. Because the
// connection is read-only, write statements fail with a clear SQLite error.
func (db *DB) Query(ctx context.Context, query string) (*Result, error) {
	rows, err := db.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	res := &Result{Columns: cols}
	for rows.Next() {
		// Scan into a slice of *any so we accept whatever type SQLite gives us.
		holders := make([]any, len(cols))
		scanTargets := make([]any, len(cols))
		for i := range holders {
			scanTargets[i] = &holders[i]
		}
		if err := rows.Scan(scanTargets...); err != nil {
			return nil, err
		}
		res.Rows = append(res.Rows, holders)
	}
	return res, rows.Err()
}

// Tables returns the names of the user tables in the database, sorted.
func (db *DB) Tables(ctx context.Context) ([]string, error) {
	const q = `SELECT name FROM sqlite_master
	           WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
	           ORDER BY name`
	res, err := db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(res.Rows))
	for _, row := range res.Rows {
		if s, ok := row[0].(string); ok {
			names = append(names, s)
		}
	}
	return names, nil
}

// Schema returns the CREATE statement for a single table.
func (db *DB) Schema(ctx context.Context, table string) (string, error) {
	res, err := db.Query(ctx,
		"SELECT sql FROM sqlite_master WHERE type='table' AND name = '"+escape(table)+"'")
	if err != nil {
		return "", err
	}
	if len(res.Rows) == 0 {
		return "", fmt.Errorf("no such table: %s", table)
	}
	if s, ok := res.Rows[0][0].(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("no schema for table: %s", table)
}

// escape doubles single quotes so a table name is safe inside a SQL literal.
func escape(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r == '\'' {
			out = append(out, '\'')
		}
		out = append(out, r)
	}
	return string(out)
}
