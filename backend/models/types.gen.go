// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package models

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Code Error code.
	Code *int64 `json:"code,omitempty"`

	// Detail Detailed information about the error.
	Detail *string `json:"detail,omitempty"`

	// Message Error message.
	Message *string `json:"message,omitempty"`
}

// Song defines model for Song.
type Song struct {
	// CreatedAt Record creation timestamp.
	CreatedAt int64 `json:"createdAt"`

	// GroupName Name of the group or artist.
	GroupName string `json:"groupName"`

	// Id Song identifier
	Id int64 `json:"id"`

	// Link Link to the song.
	Link string `json:"link"`

	// ReleaseDate Release date.
	ReleaseDate int64 `json:"releaseDate"`

	// SongText Lyrics of the song.
	SongText string `json:"songText"`

	// SongTitle Title of the song.
	SongTitle string `json:"songTitle"`

	// UpdatedAt Record update timestamp.
	UpdatedAt int64 `json:"updatedAt"`
}

// SongCreateRequest defines model for SongCreateRequest.
type SongCreateRequest struct {
	Song *Song `json:"song,omitempty"`
}

// SongTextResponse defines model for SongTextResponse.
type SongTextResponse struct {
	// Page Current page number.
	Page *int `json:"page,omitempty"`

	// PageSize Number of verses per page.
	PageSize *int `json:"pageSize,omitempty"`

	// SongId Song identifier.
	SongId *int `json:"songId,omitempty"`

	// TotalVerses Total number of verses.
	TotalVerses *int `json:"totalVerses,omitempty"`

	// Verses List of verses.
	Verses *[]string `json:"verses,omitempty"`
}

// SongUpdateRequest defines model for SongUpdateRequest.
type SongUpdateRequest struct {
	Song *Song `json:"song,omitempty"`
}

// SuccessResponse Типовой запрос для ответа на Post запросы, которые не должны возвращать никаких данных
type SuccessResponse struct {
	Success *bool `json:"success,omitempty"`
}

// GetSongsFilterParams defines parameters for GetSongsFilter.
type GetSongsFilterParams struct {
	// GroupName Filter by group name
	GroupName *string `form:"groupName,omitempty" json:"groupName,omitempty"`

	// SongTitle Filter by song title
	SongTitle *string `form:"songTitle,omitempty" json:"songTitle,omitempty"`

	// ReleaseDate Filter by release date
	ReleaseDate *int64 `form:"releaseDate,omitempty" json:"releaseDate,omitempty"`

	// Page Page number for pagination
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page
	PageSize *int `form:"pageSize,omitempty" json:"pageSize,omitempty"`
}

// GetSongsIdTextParams defines parameters for GetSongsIdText.
type GetSongsIdTextParams struct {
	// Page Page number for verse pagination.
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of verses per page.
	PageSize *int `form:"pageSize,omitempty" json:"pageSize,omitempty"`
}

// PostSongsFilterJSONRequestBody defines body for PostSongsFilter for application/json ContentType.
type PostSongsFilterJSONRequestBody = SongCreateRequest

// PatchSongsIdJSONRequestBody defines body for PatchSongsId for application/json ContentType.
type PatchSongsIdJSONRequestBody = SongUpdateRequest
