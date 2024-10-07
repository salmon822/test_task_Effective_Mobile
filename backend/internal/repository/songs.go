package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/salmon822/test_task/internal/repository/models"
)

type SongsRepository struct {
	db sqlx.ExtContext
}

func NewSongsRepository(db *sqlx.DB) Songs {
	return &SongsRepository{
		db: db,
	}
}

func (r *SongsRepository) WithTX(tx *sqlx.Tx) Songs {
	return &SongsRepository{
		db: tx,
	}
}

func (r *SongsRepository) Create(ctx context.Context, song *models.Song) (*models.Song, error) {
	query := `
		INSERT INTO songs (id, group_name, song_title, release_date, song_text, link, created_at, updated_at)
		VALUES (default, $1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	row := r.db.QueryRowxContext(ctx, query, song.GroupName, song.SongTitle, song.ReleaseDate,
		song.SongText, song.Link, song.UpdatedAt, song.CreatedAt)
	err := row.Scan(&song.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("SongsRepo/Create: error: %w", err)
	}

	return song, nil
}

func (r *SongsRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM songs
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("SongsRepo/Delete: error: %w", err)
	}

	return nil
}

func (r *SongsRepository) GetById(ctx context.Context, id int64) (*models.Song, error) {
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

	row := r.db.QueryRowxContext(ctx, query, id)
	err := row.Scan(&song.ID, &song.GroupName, &song.SongTitle,
		&song.ReleaseDate, &song.SongText, &song.Link,
		&song.CreatedAt, &song.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Songsrepo/GetById: error: %w", err)
	}

	return song, nil
}

func (r *SongsRepository) Update(ctx context.Context, data *models.Song) (*models.Song, error) {
	query := `
		UPDATE songs 
		SET group_name = $2, link = $3, release_date = $4, song_text = $5, song_title = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, data.ID, data.GroupName, data.Link, data.ReleaseDate, data.SongText, data.SongTitle)
	if err != nil {
		return nil, fmt.Errorf("SongsRepo/Update: error: %w", err)
	}

	return data, nil
}
