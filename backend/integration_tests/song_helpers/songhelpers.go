package song_helpers

import (
	"context"
	"fmt"

	"github.com/salmon822/test_task/internal/db"
)

type songOptions struct {
	groupName   *string
	songTitle   *string
	releaseDate *int64
	songText    *string
	link        *string
	createdAt   *int64
	updatedAt   *int64
}

type SongOption func(options *songOptions) error

func getDefaultSongOptions() songOptions {
	defaultString := ""
	defaultInt64 := int64(0)
	return songOptions{
		groupName:   &defaultString,
		songTitle:   &defaultString,
		releaseDate: &defaultInt64,
		songText:    &defaultString,
		link:        &defaultString,
		createdAt:   &defaultInt64,
		updatedAt:   &defaultInt64,
	}
}

func WithSongText(songText string) SongOption {
	return func(options *songOptions) error {
		if songText == "" {
			return fmt.Errorf("invalid song text %s", songText)
		}
		options.songText = &songText
		return nil
	}
}

func WithSongTitle(songTitle string) SongOption {
	return func(options *songOptions) error {
		if songTitle == "" {
			return fmt.Errorf("invalid song title %s", songTitle)
		}
		options.songTitle = &songTitle
		return nil
	}
}

func WithGroupName(groupName string) SongOption {
	return func(options *songOptions) error {
		if groupName == "" {
			return fmt.Errorf("invalid group name %s", groupName)
		}
		options.groupName = &groupName
		return nil
	}
}

func CreateSong(ctx context.Context, pgClient *db.PostgresClient, options ...SongOption) (int64, error) {
	songOptions := getDefaultSongOptions()
	for _, option := range options {
		err := option(&songOptions)
		if err != nil {
			return 0, err
		}
	}

	query := `
		INSERT INTO songs (id, group_name, song_title, release_date, song_text, link, created_at, updated_at)
		VALUES (default, $1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	row := pgClient.DB.QueryRowContext(ctx, query, *songOptions.groupName, *songOptions.songTitle, *songOptions.releaseDate, *songOptions.songText, *songOptions.link, *songOptions.createdAt, *songOptions.updatedAt)
	var songId int64
	err := row.Scan(&songId)
	if err != nil {
		return 0, err
	}

	return songId, nil
}
