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

type SongWithVerses struct {
	Song
	TotalVerses int64
	Page        int64
	PageSize    int64
	Verses      []string
}

type SongFilters struct {
	GroupName   *string
	SongTitle   *string
	ReleaseDate *int64
}
