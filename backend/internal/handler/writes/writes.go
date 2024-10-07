package writes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/salmon822/test_task/models"
)

func WriteResponseWithErrorLog(w http.ResponseWriter, code int64, resp any) {
	err := WriteResponse(w, code, resp)
	if err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func WriteErrorResponseWithErrorLog(w http.ResponseWriter, err error) {
	log.Printf("error occurred: %v", err.Error())

	writeErr := WriteErrorResponse(w, err)
	if writeErr != nil {
		log.Printf("write error response failed: %v", writeErr)
	}
}

func WriteResponse(w http.ResponseWriter, code int64, resp any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(code))
	if resp != nil {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(resp); err != nil {
			return fmt.Errorf("encode write resp: %w", err)
		}
	}
	return nil
}

func WriteErrorResponse(w http.ResponseWriter, err error) error {
	var code int64
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		code = http.StatusGatewayTimeout
	case errors.Is(err, context.Canceled):
		code = http.StatusRequestTimeout
	default:
		code = http.StatusInternalServerError
	}

	message := http.StatusText(int(code))
	detail := err.Error()

	return WriteResponse(w, code, models.ErrorResponse{
		Code:    &code,
		Detail:  &detail,
		Message: &message,
	})
}
