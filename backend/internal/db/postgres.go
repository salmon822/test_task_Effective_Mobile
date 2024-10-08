package db

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresClient struct {
	DB *sqlx.DB
}

func NewPostgresClient(ctx context.Context, DSN string) (*PostgresClient, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", DSN)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the postgresql database: %w", err)
	}

	return &PostgresClient{
		DB: db,
	}, nil
}
