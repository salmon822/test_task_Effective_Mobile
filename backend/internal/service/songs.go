package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/internal/pkg/logger"
	"github.com/salmon822/test_task/internal/repository"
	"github.com/salmon822/test_task/internal/repository/models"
	"github.com/salmon822/test_task/internal/service/converters"
)

type SongsService struct {
	transactionRepo repository.Transactions
	songsRepo       repository.Songs
	logger          logger.Logger
}

func NewSongsService(
	transactionRepo repository.Transactions,
	songsRepo repository.Songs,
	logger logger.Logger,
) Songs {
	return &SongsService{
		transactionRepo: transactionRepo,
		songsRepo:       songsRepo,
		logger:          logger,
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

func (s *SongsService) checkIfSongExists(ctx context.Context, id int64) (*domain.Song, error) {
	song, err := s.songsRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if song == nil {
		err := fmt.Errorf("song with id %d does not exist", id)
		s.logger.Warnf("Validation failed: %v", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return converters.SongModels2Domain(song), nil
}

func (s *SongsService) CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer tx.Rollback()

	songModel, err := s.songsRepo.WithTX(tx).Create(ctx, converters.SongDomain2Models(song))
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	s.logger.Infof("Song created successfully with ID: %d", songModel.ID)

	songDomain := converters.SongModels2Domain(songModel)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	return songDomain, nil
}

func (s *SongsService) DeleteSong(ctx context.Context, id int64) error {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("database error: %s", err)
	}
	defer tx.Rollback()

	err = s.songsRepo.WithTX(tx).Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("database error: %s", err)
	}

	s.logger.Infof("Song with ID %d deleted successfully", id)

	err = tx.Commit()
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
	defer tx.Rollback()

	beforeUpdate, err := s.checkIfSongExists(ctx, id)
	if err != nil {
		return nil, err
	}

	updatedSong := applyPartialUpdate(beforeUpdate, songData)

	updatedData, err := s.songsRepo.WithTX(tx).Update(ctx, converters.SongDomain2Models(updatedSong))
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	song := converters.SongModels2Domain(updatedData)

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	s.logger.Infof("Song with ID %d updated successfully", id)

	return song, nil
}

func (s *SongsService) GetSongTextByID(ctx context.Context, id int64, page int64, pageSize int64) (*domain.SongWithVerses, error) {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer tx.Rollback()

	song, err := s.songsRepo.WithTX(tx).GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	versesRaw := strings.Split(song.SongText, "\n\n")
	var verses []string

	for _, verse := range versesRaw {
		lines := strings.Split(verse, "\n")
		for _, line := range lines {
			if line != "" {
				verses = append(verses, line)
			}
		}
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= int64(len(verses)) {
		s.logger.Infof("No verses found for song ID %d on page %d", id, page)
		return converters.SongModels2DomainSongDetails(&models.SongWithVerses{
			Song:        *song,
			TotalVerses: int64(len(verses)),
			Page:        page,
			PageSize:    pageSize,
			Verses:      []string{},
		}), nil
	}

	if end > int64(len(verses)) {
		end = int64(len(verses))
	}

	s.logger.Infof("Verses retrieved successfully for song ID %d", id)

	return converters.SongModels2DomainSongDetails(&models.SongWithVerses{
		Song:        *song,
		TotalVerses: int64(len(verses)),
		Page:        page,
		PageSize:    pageSize,
		Verses:      verses[start:end],
	}), nil
}

func (s *SongsService) GetFilteredSongs(ctx context.Context, filters *domain.SongFilters, page int64, pageSize int64) ([]*domain.Song, error) {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer tx.Rollback()

	songs, err := s.songsRepo.WithTX(tx).GetFilteredSongs(ctx, converters.SongFiltersDomain2Models(filters), page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	result := domain.MapSlice(songs, converters.SongModels2Domain)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}

	s.logger.Infof("Filtered songs retrieved successfully")

	return result, nil
}
