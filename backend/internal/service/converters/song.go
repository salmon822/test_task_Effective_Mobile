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
func SongModels2DomainSongDetails(s *models.SongWithVerses) *domain.SongWithVerses {
	return &domain.SongWithVerses{
		Song:        *SongModels2Domain(&s.Song),
		TotalVerses: s.TotalVerses,
		Page:        s.Page,
		PageSize:    s.PageSize,
		Verses:      s.Verses,
	}
}

func SongFiltersModels2Domain(s *models.SongFilters) *domain.SongFilters {
	return &domain.SongFilters{
		GroupName:   s.GroupName,
		SongTitle:   s.SongTitle,
		ReleaseDate: s.ReleaseDate,
	}
}

func SongFiltersDomain2Models(s *domain.SongFilters) *models.SongFilters {
	return &models.SongFilters{
		GroupName:   s.GroupName,
		SongTitle:   s.SongTitle,
		ReleaseDate: s.ReleaseDate,
	}
}
