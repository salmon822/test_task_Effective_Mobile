package integration_tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/salmon822/test_task/integration_tests/song_helpers"
	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/models"
	"github.com/stretchr/testify/require"
)

type SongSuite struct {
	TestSuite
}

func (s *SongSuite) SetupSuite() {
	s.RunTestServer()
}

func (s *SongSuite) TestCreateSongSuccess(t *testing.T) {
	songToCreate := &models.Song{
		GroupName:   "Test Group",
		SongTitle:   "Test Song",
		ReleaseDate: 20220101,
		SongText:    "Test song text",
		Link:        "http://testlink.com",
	}
	req := models.SongCreateRequest{
		Song: songToCreate,
	}
	var songCreateRes models.Song

	expectedCreateRes := *songToCreate
	expectedCreateRes.Id = 1

	_, err := makeJsonRequest(s.handler.Init(), http.MethodPost, "/songs/create", req, &songCreateRes)
	require.NoError(t, err)

	require.Equal(t, expectedCreateRes, songCreateRes)
}

func (s *SongSuite) TestDeleteSongSuccess(t *testing.T) {
	ctx := context.Background()

	createdSong, err := song_helpers.CreateSong(ctx, s.pgClient)
	require.NoError(t, err)

	var successRes models.SuccessResponse
	expectedDeleteRes := models.SuccessResponse{
		Success: makePointer(true),
	}

	_, err = makeJsonRequest(s.handler.Init(), http.MethodDelete, fmt.Sprintf("/songs/%d/delete", createdSong), nil, successRes)
	require.NoError(t, err)

	require.Equal(t, expectedDeleteRes, successRes)
}

func (s *SongSuite) TestUpdateSongSuccess(t *testing.T) {
	ctx := context.Background()

	createdSong, err := song_helpers.CreateSong(ctx, s.pgClient)
	require.NoError(t, err)

	songUpdate := &models.Song{
		GroupName:   "Updated Group",
		SongTitle:   "Updated Song",
		ReleaseDate: 20220102,
		SongText:    "Updated song text",
		Link:        "http://updatedlink.com",
	}
	req := models.SongUpdateRequest{
		Song: songUpdate,
	}
	var songUpdateRes models.Song

	expectedUpdateRes := *songUpdate
	expectedUpdateRes.Id = createdSong

	_, err = makeJsonRequest(s.handler.Init(), http.MethodPatch, fmt.Sprintf("/songs/%d/update", createdSong), req, &songUpdateRes)
	require.NoError(t, err)

	require.Equal(t, expectedUpdateRes, songUpdateRes)
}

func (s *SongSuite) TestGetSongTextSuccess(t *testing.T) {
	ctx := context.Background()

	createdSong, err := song_helpers.CreateSong(ctx, s.pgClient)
	require.NoError(t, err)

	url := fmt.Sprintf("/songs/%d/song-text?page=1&pageSize=2", createdSong)

	var songTextRes domain.SongWithVerses

	_, err = makeJsonRequest(s.handler.Init(), http.MethodGet, url, nil, &songTextRes)
	require.NoError(t, err)

	require.Equal(t, 2, len(songTextRes.Verses))
}

func (s *SongSuite) TestGetFilteredSongsSuccess(t *testing.T) {
	ctx := context.Background()

	_, err := song_helpers.CreateSong(ctx, s.pgClient, song_helpers.WithGroupName("TestGroup"))
	require.NoError(t, err)

	_, err = song_helpers.CreateSong(ctx, s.pgClient, song_helpers.WithGroupName("Just"))
	require.NoError(t, err)

	url := "/songs/filter?groupName=TestGroup&songTitle=TestSong"

	var filteredSongs []models.Song

	_, err = makeJsonRequest(s.handler.Init(), http.MethodGet, url, nil, &filteredSongs)
	require.NoError(t, err)

	require.NotEmpty(t, filteredSongs)
	require.Equal(t, "Test Group", filteredSongs[0].GroupName)
}
