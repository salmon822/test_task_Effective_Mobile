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
		h.logger.Errorf("Failed to decode request body: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to decode request body: %w", err))
		return
	}
	h.logger.Infof("SongCreateRequest decoded successfully")

	if err := req.Validate(h.validationFormats); err != nil {
		h.logger.Errorf("Validation failed: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("validation failed: %w", err))
		return
	}

	res, err := h.songs.CreateSong(ctx, domain.SongModels2Domain(req.Song))
	if err != nil {
		h.logger.Errorf("Failed to create song: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to create song: %w", err))
		return
	}

	h.logger.Infof("Song created successfully with ID: %d", res.ID)
	song := domain.SongDomain2Models(res)

	writes.WriteResponseWithErrorLog(w, http.StatusOK, song)
}

func (h *handler) deleteSong(w http.ResponseWriter, r *http.Request) {
	var id int64
	defer r.Body.Close()

	if err := h.parsePathInt64Param(r, "id", &id); err != nil {
		h.logger.Errorf("Failed to parse ID from path: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	if err := h.songs.DeleteSong(ctx, id); err != nil {
		h.logger.Errorf("Failed to delete song with ID %d: %v", id, err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to delete song: %w", err))
		return
	}

	h.logger.Infof("Song deleted successfully with ID: %d", id)
	writes.WriteResponseWithErrorLog(w, http.StatusOK, successResponse(true))
}

func (h *handler) updateSong(w http.ResponseWriter, r *http.Request) {
	var id int64
	defer r.Body.Close()

	if err := h.parsePathInt64Param(r, "id", &id); err != nil {
		h.logger.Errorf("Failed to parse ID from path: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}

	var req models.SongUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Failed to decode request body: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to decode request body: %w", err))
		return
	}

	if err := req.Validate(h.validationFormats); err != nil {
		h.logger.Errorf("Validation failed: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("validation failed: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	songToUpdate, err := h.songs.UpdateSong(ctx, id, domain.SongModels2Domain(req.Song))
	if err != nil {
		h.logger.Errorf("Failed to update song with ID %d: %v", id, err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to update song: %w", err))
		return
	}

	h.logger.Infof("Song updated successfully with ID: %d", id)
	song := domain.SongDomain2Models(songToUpdate)

	writes.WriteResponseWithErrorLog(w, http.StatusOK, song)
}

func (h *handler) getSongText(w http.ResponseWriter, r *http.Request) {
	var id, page, pageSize int64
	defer r.Body.Close()

	if err := h.parsePathInt64Param(r, "id", &id); err != nil {
		h.logger.Errorf("Failed to parse ID from path: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}

	page, err := h.parseQueryInt64Param(r, "page", 1)
	if err != nil {
		h.logger.Errorf("Failed to parse page: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}
	pageSize, err = h.parseQueryInt64Param(r, "pageSize", 5)
	if err != nil {
		h.logger.Errorf("Failed to parse pageSize: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	res, err := h.songs.GetSongTextByID(ctx, id, page, pageSize)
	if err != nil {
		h.logger.Errorf("Failed to get song text for ID %d: %v", id, err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to get song text: %w", err))
		return
	}

	h.logger.Infof("Retrieved song text successfully for ID: %d", id)
	writes.WriteResponseWithErrorLog(w, http.StatusOK, res)
}

func (h *handler) getFilteredSongs(w http.ResponseWriter, r *http.Request) {
	var filters domain.SongFilters
	defer r.Body.Close()

	if err := h.parseQueryStringParam(r, "groupName", &filters.GroupName); err != nil {
		h.logger.Errorf("Failed to parse groupName: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}
	if err := h.parseQueryStringParam(r, "songTitle", &filters.SongTitle); err != nil {
		h.logger.Errorf("Failed to parse songTitle: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}
	releaseDateParam := r.URL.Query().Get("releaseDate")
	if releaseDateParam != "" {
		releaseDateValue, err := h.parseQueryInt64Param(r, "releaseDate", 0)
		if err != nil {
			h.logger.Errorf("Failed to parse releaseDate: %v", err)
			writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
			return
		}
		filters.ReleaseDate = &releaseDateValue
	}

	page, err := h.parseQueryInt64Param(r, "page", 1)
	if err != nil {
		h.logger.Errorf("Failed to parse page: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}
	pageSize, err := h.parseQueryInt64Param(r, "pageSize", 5)
	if err != nil {
		h.logger.Errorf("Failed to parse pageSize: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("parse failed: %w", err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	res, err := h.songs.GetFilteredSongs(ctx, &filters, page, pageSize)
	if err != nil {
		h.logger.Errorf("Failed to get filtered songs: %v", err)
		writes.WriteErrorResponseWithErrorLog(w, fmt.Errorf("failed to get filtered songs: %w", err))
		return
	}

	h.logger.Infof("Retrieved filtered songs successfully")
	writes.WriteResponseWithErrorLog(w, http.StatusOK, res)
}
