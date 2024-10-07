package models

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
