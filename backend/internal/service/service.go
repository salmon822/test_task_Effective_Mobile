package service

import (
	"context"

	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/internal/pkg/logger"
	"github.com/salmon822/test_task/internal/repository"
)

type Songs interface {
	CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error)
	DeleteSong(ctx context.Context, id int64) error
	UpdateSong(ctx context.Context, id int64, songData *domain.Song) (*domain.Song, error)
	GetSongTextByID(ctx context.Context, id, page, pageSize int64) (*domain.SongWithVerses, error)
	GetFilteredSongs(ctx context.Context, filters *domain.SongFilters, page int64, pageSize int64) ([]*domain.Song, error)
}

type Service struct {
	Songs
	logger logger.Logger
}

func NewService(
	ctx context.Context,
	repo *repository.Repository,
	logger logger.Logger,
) (Service, error) {

	var (
		songs = NewSongsService(repo.Transactions, repo.Songs, logger)
	)

	res := Service{
		Songs:  songs,
		logger: logger,
	}

	return res, nil
}
