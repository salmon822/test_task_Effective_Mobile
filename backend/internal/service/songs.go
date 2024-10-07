package service

import (
	"context"
	"fmt"

	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/internal/repository"
	"github.com/salmon822/test_task/internal/service/converters"
)

type SongsService struct {
	transactionRepo repository.Transactions
	songsRepo       repository.Songs
}

func NewSongsService(
	transactionRepo repository.Transactions,
	songsRepo repository.Songs,
) Songs {
	return &SongsService{
		transactionRepo: transactionRepo,
		songsRepo:       songsRepo,
	}
}

func applyPartialUpdate(existingSong, songData *domain.Song) *domain.Song {
	if songData.GroupName != "" {
		existingSong.GroupName = songData.GroupName
	}
	if songData.SongTitle != "" {
		existingSong.SongTitle = songData.SongTitle
	}
	if songData.Link != "" {
		existingSong.Link = songData.Link
	}
	if songData.ReleaseDate != 0 {
		existingSong.ReleaseDate = songData.ReleaseDate
	}
	if songData.SongText != "" {
		existingSong.SongText = songData.SongText
	}

	return existingSong
}

func (s *SongsService) checkIfSongExists(ctx context.Context, transaction repository.Transaction, id int64) (*domain.Song, error) {
	song, err := s.songsRepo.GetByIdTX(ctx, transaction, id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if song == nil {
		err := fmt.Errorf("song with id %d does not exist", id)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return converters.SongModels2Domain(song), nil

}

func (s *SongsService) CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer tx.Rollback(ctx)

	songModel, err := s.songsRepo.CreateTX(ctx, tx, converters.SongDomain2Models(song))
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	songDomain := converters.SongModels2Domain(songModel)

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	return songDomain, nil
}

func (s *SongsService) DeleteSong(ctx context.Context, id int64) error {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("database error: %s", err)
	}
	defer tx.Rollback(ctx)

	err = s.songsRepo.DeleteTX(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("database error: %s", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("database error: %s", err)
	}

	return nil
}

func (s *SongsService) UpdateSong(ctx context.Context, id int64, songData *domain.Song) (*domain.Song, error) {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer tx.Rollback(ctx)

	beforeUpdate, err := s.checkIfSongExists(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	updatedSong := applyPartialUpdate(beforeUpdate, songData)

	updatedData, err := s.songsRepo.UpdateTX(ctx, tx, converters.SongDomain2Models(updatedSong))
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	song := converters.SongModels2Domain(updatedData)

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	return song, nil
}
