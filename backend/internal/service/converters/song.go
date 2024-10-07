package converters

import (
	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/internal/repository/models"
)

func SongDomain2Models(s *domain.Song) *models.Song {
	if s == nil {
		return nil
	}
	song := &models.Song{
		ID:          s.ID,
		GroupName:   s.GroupName,
		SongTitle:   s.SongTitle,
		ReleaseDate: s.ReleaseDate,
		SongText:    s.SongText,
		Link:        s.Link,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	return song
}

func SongModels2Domain(s *models.Song) *domain.Song {
	if s == nil {
		return nil
	}
	song := &domain.Song{
		ID:          s.ID,
		GroupName:   s.GroupName,
		SongTitle:   s.SongTitle,
		ReleaseDate: s.ReleaseDate,
		SongText:    s.SongText,
		Link:        s.Link,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	return song
}
