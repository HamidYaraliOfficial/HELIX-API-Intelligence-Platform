package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//go:embed all:migrations
var migrationsFS embed.FS

// Connect opens a PostgreSQL connection pool and blocks (with backoff)
// until the database is reachable. This makes the service resilient to
// the common Docker Compose race where the app container starts before
// Postgres finishes initializing.
func Connect(dsn string) (*sql.DB, error) {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(10)
	dbConn.SetConnMaxLifetime(30 * time.Minute)

	var lastErr error
	for attempt := 1; attempt <= 20; attempt++ {
		lastErr = dbConn.Ping()
		if lastErr == nil {
			return dbConn, nil
		}
		log.Printf("[db] waiting for postgres (attempt %d/20): %v", attempt, lastErr)
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %w", lastErr)
}

// Migrate applies every embedded .sql migration file exactly once, tracked
// via the schema_migrations table. Files are applied in lexical order, so
// migrations must be named with a numeric/date prefix (0001_, 0002_, ...).
func Migrate(dbConn *sql.DB) error {
	if _, err := dbConn.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename   TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}

	var filenames []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			filenames = append(filenames, e.Name())
		}
	}
	sort.Strings(filenames)

	for _, filename := range filenames {
		var alreadyApplied bool
		err := dbConn.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`,
			filename,
		).Scan(&alreadyApplied)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", filename, err)
		}
		if alreadyApplied {
			continue
		}

		content, err := migrationsFS.ReadFile(path.Join("migrations", filename))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", filename, err)
		}

		tx, err := dbConn.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", filename, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", filename, err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (filename) VALUES ($1)`, filename); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", filename, err)
		}

		log.Printf("[db] applied migration %s", filename)
	}

	return nil
}
