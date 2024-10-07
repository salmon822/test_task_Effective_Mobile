package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/salmon822/test_task/internal/domain"
	"github.com/salmon822/test_task/internal/handler/writes"
	"github.com/salmon822/test_task/models"
)

func (h *handler) createSong(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	var req models.SongCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to decode request body: %w", err))
		return
	}
	if err := req.Validate(h.validationFormats); err != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("validation failed: %w", err))
		return
	}

	res, err := h.songs.CreateSong(ctx, domain.SongModels2Domain(req.Song))
	if err != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to create song: %w", err))
		return
	}

	song := domain.SongDomain2Models(res)

	writes.WriteResponseWithErrorLog(w, http.StatusOK, song)
}

func (h *handler) deleteSong(w http.ResponseWriter, r *http.Request) {
	var (
		id int64
	)
	defer r.Body.Close()
	if err := h.parsePathInt64Param(r, "id", &id); err != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	e := h.songs.DeleteSong(ctx, id)
	if e != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to delete song: %w", e))
		return
	}

	writes.WriteResponseWithErrorLog(w, http.StatusOK, successResponse(true))
}

func (h *handler) updateSong(w http.ResponseWriter, r *http.Request) {
	var (
		id int64
	)
	defer r.Body.Close()
	if err := h.parsePathInt64Param(r, "id", &id); err != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}

	var req models.SongUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to decode request body: %w", err))
		return
	}
	if err := req.Validate(h.validationFormats); err != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("validation failed: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	songToUpdate, e := h.songs.UpdateSong(ctx, id, domain.SongModels2Domain(req.Song))
	if e != nil {
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to update song: %w", e))
		return
	}

	song := domain.SongDomain2Models(songToUpdate)

	writes.WriteResponseWithErrorLog(w, http.StatusOK, song)
}

// func (h *handler) getSongText(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		id       int64
// 		page     int64
// 		pageSize int64
// 	)
// 	defer r.Body.Close()
// 	if err := h.parsePathInt64Param(r, "id", &id); err != nil {
// 		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
// 		return
// 	}

// 	page, err := h.parseQueryInt64Param(r, "page", 1)
// 	if err != nil {
// 		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
// 		return
// 	}
// 	pageSize, err = h.parseQueryInt64Param(r, "pageSize", 5)
// 	if err != nil {
// 		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
// 	defer cancel()

// 	res, err := h.songs.GetSongTextByID(ctx, id, page, pageSize)
// 	if err != nil {
// 		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("fa: %w", err))
// 		return
// 	}

// 	writes.WriteResponseWithErrorLog(w, http.StatusOK, res)
// }
