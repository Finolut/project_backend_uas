package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// ConnectPostgres opens and returns a *sql.DB using the provided DSN.
// DSN format example: "postgres://user:pass@host:port/dbname?sslmode=disable"
func ConnectPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// tuning and pool settings
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(5)

	// ping to ensure reachable
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
