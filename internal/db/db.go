package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx stdlib driver
)

func Open() (*sql.DB, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	// pgx registers as "pgx"
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	// reasonable pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, nil
}

func Ping(db *sql.DB) error {
	ctx := context.Background()
	return db.PingContext(ctx)
}
