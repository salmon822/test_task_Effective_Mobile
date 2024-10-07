-- +goose Up
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_title VARCHAR(255) NOT NULL,
    release_date BIGINT NOT NULL,
    song_text TEXT,
    link VARCHAR(255),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE INDEX idx_songs_group_name ON songs(group_name);
CREATE INDEX idx_songs_song_title ON songs(song_title);

-- +goose Down
DROP TABLE IF EXISTS songs;

