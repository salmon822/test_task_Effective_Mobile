package service

import (
	"context"

	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/internal/repository"
)

type Songs interface {
	CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error)
	DeleteSong(ctx context.Context, id int64) error
	UpdateSong(ctx context.Context, id int64, songData *domain.Song) (*domain.Song, error)
	GetSongTextByID(ctx context.Context, id, page, pageSize int64) (*domain.SongWithVerses, error)
}

type Service struct {
	Songs

	ctx context.Context
}

func NewService(
	ctx context.Context,
	repo *repository.Repository,
) (Service, error) {

	var (
		songs = NewSongsService(repo.Transactions, repo.Songs)
	)

	res := Service{
		Songs: songs,
	}

	return res, nil
}
