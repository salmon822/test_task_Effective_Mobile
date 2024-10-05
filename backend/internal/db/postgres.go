package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	DB *pgxpool.Pool
}

func NewPostgresClient(ctx context.Context, DSN string) (*PostgresClient, error) {
	cfg, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MaxConnIdleTime = 5 * time.Second

	client, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the postgresql database: %w", err)
	}
	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error while ping to postgres")
	}

	return &PostgresClient{
		DB: client,
	}, nil
}

func (c *PostgresClient) Close() {
	c.DB.Close()
}
