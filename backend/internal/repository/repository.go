package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/salmon822/test_task/internal/config"
	"github.com/salmon822/test_task/internal/pkg/logger"
	"github.com/salmon822/test_task/internal/repository/models"
)

type Songs interface {
	Create(ctx context.Context, song *models.Song) (*models.Song, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*models.Song, error)
	Update(ctx context.Context, data *models.Song) (*models.Song, error)
	GetFilteredSongs(ctx context.Context, filters *models.SongFilters, page int64, pageSize int64) ([]*models.Song, error)
	WithTX(tx *sqlx.Tx) Songs
}

type Transactions interface {
	StartTransaction(ctx context.Context) (*sqlx.Tx, error)
}

type Repository struct {
	Transactions
	Songs
	logger logger.Logger
}

func NewRepository(
	cfg *config.Config,
	db *sqlx.DB,
	logger logger.Logger,
) (*Repository, error) {
	var (
		songs        = NewSongsRepository(db, logger)
		transactions = NewTransactionsRepo(db)
	)

	return &Repository{
		Transactions: transactions,
		Songs:        songs,
		logger:       logger,
	}, nil
}
