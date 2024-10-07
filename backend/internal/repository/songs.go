package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/salmon822/test_task/internal/repository/models"
)

type SongsRepository struct {
}

func NewSongsRepository() Songs {
	return &SongsRepository{}
}

func (r *SongsRepository) CreateTX(ctx context.Context, transaction Transaction, song *models.Song) (*models.Song, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("SongsRepo/Create: error: type assertion failed on interface Transaction")
	}
	query := `
		INSERT INTO songs (id, group_name, song_title, release_date, song_text, link, created_at, updated_at)
		VALUES (default, $1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	row := tx.QueryRow(ctx, query, song.GroupName, song.SongTitle, song.ReleaseDate, song.SongText, song.Link, song.UpdatedAt, song.CreatedAt)
	err := row.Scan(&song.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("SongsRepo/Create: error: %w", err)
	}

	return song, nil
}

func (r *SongsRepository) DeleteTX(ctx context.Context, transaction Transaction, id int64) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return fmt.Errorf("SongsRepo/Delete: error: type assertion failed on interface Transaction")
	}
	query := `
		DELETE FROM songs
		WHERE id = $1
	`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("SongsRepo/Delete: error: %w", err)
	}

	return nil
}

func (r *SongsRepository) GetByIdTX(ctx context.Context, transaction Transaction, id int64) (*models.Song, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("SongsRepo/GetById: error: type assertion failed on interface Transaction")
	}
	var (
		song = &models.Song{
			ID: id,
		}
	)

	query := `
		SELECT id, group_name, song_title, release_date, song_text, link, created_at, updated_at 
		FROM songs
		WHERE id = $1
			`

	row := tx.QueryRow(ctx, query, id)
	err := row.Scan(&song.ID, &song.GroupName, &song.SongTitle,
		&song.ReleaseDate, &song.SongText, &song.Link,
		&song.CreatedAt, &song.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Songsrepo/GetById: error: %w", err)
	}

	return song, nil
}

func (r *SongsRepository) UpdateTX(ctx context.Context, transaction Transaction, data *models.Song) (*models.Song, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("SongsRepo/Update: error: type assertion failed on interface Transaction")
	}
	query := `
		UPDATE songs 
		SET group_name = $2, link = $3, release_date = $4, song_text = $5, song_title = $6
		WHERE id = $1
	`
	_, err := tx.Exec(ctx, query, data.ID, data.GroupName, data.Link, data.ReleaseDate, data.SongText, data.SongTitle)
	if err != nil {
		return nil, fmt.Errorf("SongsRepo/Update: error: %w", err)
	}

	return data, nil
}
