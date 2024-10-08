package integration_tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/salmon822/test_task/integration_tests/song_helpers"
	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/models"
)

type SongSuite struct {
	TestSuite
}

func (s *SongSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
}

func (s *SongSuite) TestCreateSongSuccess() {
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

	_, err := makeJsonRequest(s.httpHandler, http.MethodPost, "/songs/create", req, &songCreateRes)
	s.Require().NoError(err)

	s.Require().Equal(expectedCreateRes, songCreateRes)
}

func (s *SongSuite) TestDeleteSongSuccess() {
	ctx := context.Background()

	createdSong, err := song_helpers.CreateSong(ctx, s.pgClient)
	s.Require().NoError(err)

	var successRes models.SuccessResponse
	expectedDeleteRes := models.SuccessResponse{
		Success: makePointer(true),
	}

	_, err = makeJsonRequest(s.httpHandler, http.MethodDelete, fmt.Sprintf("/songs/%d/delete", createdSong), nil, &successRes)
	s.Require().NoError(err)

	s.Require().Equal(expectedDeleteRes, successRes)
}

func (s *SongSuite) TestUpdateSongSuccess() {
	ctx := context.Background()

	createdSong, err := song_helpers.CreateSong(ctx, s.pgClient)
	s.Require().NoError(err)

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

	_, err = makeJsonRequest(s.httpHandler, http.MethodPatch, fmt.Sprintf("/songs/%d/update", createdSong), req, &songUpdateRes)
	s.Require().NoError(err)

	s.Require().Equal(expectedUpdateRes, songUpdateRes)
}

func (s *SongSuite) TestGetSongTextSuccess() {
	ctx := context.Background()

	createdSong, err := song_helpers.CreateSong(ctx, s.pgClient, song_helpers.WithSongText("Verse 1 line 1\nVerse 1 line 2\n\nVerse 2 line 1\nVerse 2 line 2"))
	s.Require().NoError(err)

	url := fmt.Sprintf("/songs/%d/song-text?page=1&pageSize=2", createdSong)

	var songTextRes domain.SongWithVerses

	_, err = makeJsonRequest(s.httpHandler, http.MethodPost, url, nil, &songTextRes)
	s.Require().NoError(err)

	s.Require().Equal(2, len(songTextRes.Verses))
}

func (s *SongSuite) TestGetFilteredSongsSuccess() {
	ctx := context.Background()

	_, err := song_helpers.CreateSong(ctx, s.pgClient,
		song_helpers.WithGroupName("TestGroup"),
		song_helpers.WithSongTitle("TestSong"))
	s.Require().NoError(err)

	_, err = song_helpers.CreateSong(ctx, s.pgClient,
		song_helpers.WithGroupName("Just"),
		song_helpers.WithSongTitle("AnotherSong"))
	s.Require().NoError(err)

	url := "/songs/filter?groupName=TestGroup&songTitle=TestSong"

	var filteredSongs []models.Song

	_, err = makeJsonRequest(s.httpHandler, http.MethodGet, url, nil, &filteredSongs)
	s.Require().NoError(err)

	s.Require().NotEmpty(filteredSongs)

	s.Require().Equal("TestGroup", filteredSongs[0].GroupName)
	s.Require().Equal("TestSong", filteredSongs[0].SongTitle)
}
