package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/salmon822/test_task/internal/config"
	"github.com/salmon822/test_task/internal/repository/models"
)

type Songs interface {
	CreateTX(ctx context.Context, transaction Transaction, song *models.Song) (*models.Song, error)
	DeleteTX(ctx context.Context, transaction Transaction, id int64) error
	GetByIdTX(ctx context.Context, transaction Transaction, id int64) (*models.Song, error)
	UpdateTX(ctx context.Context, transaction Transaction, data *models.Song) (*models.Song, error)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Transactions interface {
	StartTransaction(ctx context.Context) (Transaction, error)
}

type Repository struct {
	Transactions
	Songs
}

func NewRepository(cfg *config.Config, pool *pgxpool.Pool) (*Repository, error) {
	var (
		songs        = NewSongsRepository()
		transactions = NewTransactionsRepo(pool)
	)

	return &Repository{
		Transactions: transactions,
		Songs:        songs,
	}, nil
}
