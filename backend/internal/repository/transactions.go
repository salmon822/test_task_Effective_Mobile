package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionsRepo struct {
	pool *pgxpool.Pool
}

func NewTransactionsRepo(pool *pgxpool.Pool) *TransactionsRepo {
	return &TransactionsRepo{pool: pool}
}

func (r *TransactionsRepo) StartTransaction(ctx context.Context) (Transaction, error) {
	return r.pool.Begin(ctx)
}
