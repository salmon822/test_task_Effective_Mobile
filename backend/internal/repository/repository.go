package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/salmon822/test_task/internal/config"
	"github.com/salmon822/test_task/internal/repository/models"
)

type Songs interface {
	Create(ctx context.Context, song *models.Song) (*models.Song, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*models.Song, error)
	Update(ctx context.Context, data *models.Song) (*models.Song, error)
	WithTX(tx *sqlx.Tx) Songs
}

type Transactions interface {
	StartTransaction(ctx context.Context) (*sqlx.Tx, error)
}

type Repository struct {
	Transactions
	Songs
}

func NewRepository(cfg *config.Config, db *sqlx.DB) (*Repository, error) {
	var (
		songs        = NewSongsRepository(db)
		transactions = NewTransactionsRepo(db)
	)

	return &Repository{
		Transactions: transactions,
		Songs:        songs,
	}, nil
}
