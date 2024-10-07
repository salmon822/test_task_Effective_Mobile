package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionsRepo struct {
	db *sqlx.DB
}

func NewTransactionsRepo(db *sqlx.DB) *TransactionsRepo {
	return &TransactionsRepo{db: db}
}

func (r *TransactionsRepo) StartTransaction(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.Beginx()
}
