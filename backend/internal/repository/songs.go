package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/salmon822/test_task/internal/pkg/logger"
	"github.com/salmon822/test_task/internal/repository/models"
)

type SongsRepository struct {
	db     sqlx.ExtContext
	logger logger.Logger
}

func NewSongsRepository(
	db *sqlx.DB,
	logger logger.Logger,
) Songs {
	return &SongsRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SongsRepository) WithTX(tx *sqlx.Tx) Songs {
	return &SongsRepository{
		db:     tx,
		logger: r.logger,
	}
}

func (r *SongsRepository) Create(ctx context.Context, song *models.Song) (*models.Song, error) {
	if r.logger == nil {
		return nil, fmt.Errorf("SongsRepo/Create: logger is nil")
	}
	query := `
		INSERT INTO songs (id, group_name, song_title, release_date, song_text, link, created_at, updated_at)
		VALUES (default, $1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	row := r.db.QueryRowxContext(ctx, query, song.GroupName, song.SongTitle, song.ReleaseDate,
		song.SongText, song.Link, song.UpdatedAt, song.CreatedAt)

	r.logger.Debugf("SQL Query: %s", query)

	err := row.Scan(&song.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warnf("No rows returned for song creation")
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

	r.logger.Debugf("SQL Query: %s", query)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("SongsRepo/Delete: error: %w", err)
	}

	return nil
}

func (r *SongsRepository) GetById(ctx context.Context, id int64) (*models.Song, error) {
	query := `
		SELECT id, group_name, song_title, release_date, song_text, link, created_at, updated_at 
		FROM songs
		WHERE id = $1
	`

	r.logger.Debugf("SQL Query: %s", query)

	var song = &models.Song{
		ID: id,
	}

	row := r.db.QueryRowxContext(ctx, query, id)
	err := row.Scan(&song.ID, &song.GroupName, &song.SongTitle,
		&song.ReleaseDate, &song.SongText, &song.Link,
		&song.CreatedAt, &song.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("SongsRepo/GetById: error: %w", err)
	}

	return song, nil
}

func (r *SongsRepository) Update(ctx context.Context, data *models.Song) (*models.Song, error) {
	query := `
		UPDATE songs 
		SET group_name = $2, link = $3, release_date = $4, song_text = $5, song_title = $6
		WHERE id = $1
	`

	r.logger.Debugf("SQL Query: %s", query)

	_, err := r.db.ExecContext(ctx, query, data.ID, data.GroupName, data.Link, data.ReleaseDate, data.SongText, data.SongTitle)
	if err != nil {
		return nil, fmt.Errorf("SongsRepo/Update: error: %w", err)
	}

	return data, nil
}

func (r *SongsRepository) GetFilteredSongs(ctx context.Context, filters *models.SongFilters, page int64, pageSize int64) ([]*models.Song, error) {
	query := `
		SELECT id, group_name, song_title, release_date, song_text, link, created_at, updated_at
		FROM songs
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filters.GroupName != nil {
		query += fmt.Sprintf(" AND group_name ILIKE $%d", argIndex)
		args = append(args, "%"+*filters.GroupName+"%")
		argIndex++
	}

	if filters.SongTitle != nil {
		query += fmt.Sprintf(" AND song_title ILIKE $%d", argIndex)
		args = append(args, "%"+*filters.SongTitle+"%")
		argIndex++
	}

	if filters.ReleaseDate != nil {
		query += fmt.Sprintf(" AND release_date = $%d", argIndex)
		args = append(args, *filters.ReleaseDate)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, (page-1)*pageSize)

	r.logger.Debugf("SQL Query: %s", query)
	r.logger.Debugf("Query Arguments: %+v", args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SongsRepo/GetFilteredSongs: error executing query: %w", err)
	}
	defer rows.Close()

	var songs []*models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.GroupName, &song.SongTitle, &song.ReleaseDate, &song.SongText, &song.Link, &song.CreatedAt, &song.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("SongsRepo/GetFilteredSongs: error scanning row: %w", err)
		}
		songs = append(songs, &song)
	}

	return songs, nil
}
