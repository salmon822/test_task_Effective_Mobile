package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"github.com/salmon822/test_task/internal/config"
	"github.com/salmon822/test_task/internal/pkg/logger"
	"github.com/salmon822/test_task/internal/service"
	"github.com/salmon822/test_task/models"
)

type Handler interface {
	Init() http.Handler
}

type handler struct {
	songs             service.Songs
	cfg               *config.HandlerConfig
	logger            logger.Logger
	validationFormats strfmt.Registry
}

func NewHandler(
	songs service.Songs,
	cfg *config.HandlerConfig,
	logger logger.Logger,
) Handler {
	return &handler{
		songs:             songs,
		cfg:               cfg,
		logger:            logger,
		validationFormats: strfmt.NewFormats(),
	}
}

func (h *handler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func successResponse(s bool) *models.SuccessResponse {
	return &models.SuccessResponse{
		Success: &s,
	}
}

func (h *handler) parsePathInt64Param(r *http.Request, paramName string, paramValue *int64) error {
	param := mux.Vars(r)[paramName]
	if param == "" {
		return fmt.Errorf("parsePathInt64Param: %s", "param is empty")
	}

	var err error
	*paramValue, err = strconv.ParseInt(param, 10, 64)
	if err != nil {
		return fmt.Errorf("parsePathInt64Param: %w", err)
	}

	return nil
}

func (h *handler) parseQueryInt64Param(r *http.Request, paramName string, defaultValue int64) (int64, error) {
	param := r.URL.Query().Get(paramName)
	if param == "" {
		return defaultValue, nil
	}

	paramValue, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parseQueryInt64Param: %w", err)
	}

	return paramValue, nil
}

func (h *handler) parseQueryStringParam(r *http.Request, paramName string, dest **string) error {
	param := r.URL.Query().Get(paramName)
	if param == "" {
		return nil
	}

	*dest = &param
	return nil
}

func (h *handler) Init() http.Handler {
	router := mux.NewRouter()

	songsRouter := router.PathPrefix("/songs").Subrouter()
	songsRouter.Handle("/create", http.HandlerFunc(h.createSong)).Methods(http.MethodPost)
	songsRouter.Handle("/{id}/delete", http.HandlerFunc(h.deleteSong)).Methods(http.MethodDelete)
	songsRouter.Handle("/{id}/update", http.HandlerFunc(h.updateSong)).Methods(http.MethodPatch)
	songsRouter.Handle("/{id}/song-text", http.HandlerFunc(h.getSongText)).Methods(http.MethodPost)
	songsRouter.Handle("/filter", http.HandlerFunc(h.getFilteredSongs)).Methods(http.MethodGet)

	router.Use(h.corsMiddleware)
	return router
}
