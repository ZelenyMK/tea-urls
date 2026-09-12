package storage

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

func OpenDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	conn.Exec(context.Background(), "CREATE TABLE links (id BIGSERIAL PRIMARY KEY, original_url TEXT NOT NULL, shorten_url TEXT NOT NULL, alias TEXT NOT NULL UNIQUE)")
	conn.Exec(context.Background(), "ALTER TABLE links ADD CONSTRAINT links_alias_unique UNIQUE (alias)")
	return conn, nil
}

func CloseDB(conn *pgx.Conn) {
	conn.Close(context.Background())
}
