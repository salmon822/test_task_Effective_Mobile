package domain

import (
	"github.com/salmon822/test_task/models"
)

type Song struct {
	ID          int64
	GroupName   string
	SongTitle   string
	ReleaseDate int64
	SongText    string
	Link        string
	CreatedAt   int64
	UpdatedAt   int64
}

type SongWithVerses struct {
	Song
	TotalVerses int64
	Page        int64
	PageSize    int64
	Verses      []string
}

func SongDomain2Models(s *Song) *models.Song {
	if s == nil {
		return nil
	}
	song := &models.Song{
		Id:          s.ID,
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

func SongModels2Domain(s *models.Song) *Song {
	if s == nil {
		return nil
	}
	song := &Song{
		ID:          s.Id,
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
