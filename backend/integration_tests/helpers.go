package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/salmon822/test_task/models"
)

func makePointer[T any](value T) *T {
	return &value
}

func makeJsonRequest(handler http.Handler, method string, url string, body any, res any) (string, error) {
	var b io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		b = bytes.NewReader(data)
	}
	httpReq := httptest.NewRequest(method, url, b)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httpReq)
	resp := recorder.Result()
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non 200 response: %d %s", resp.StatusCode, string(data))
	}

	if res != nil {
		if err = json.Unmarshal(data, res); err != nil {
			return "", fmt.Errorf("json unmarshal failed: %w %s", err, string(data))
		}
	}
	return resp.Header.Get("Set-Cookie"), nil
}

func makeJsonRequestWithErrorResp(handler http.Handler, method string, url string, body any) (models.ErrorResponse, error) {
	var b io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return models.ErrorResponse{}, err
		}
		b = bytes.NewReader(data)
	}
	httpReq := httptest.NewRequest(method, url, b)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httpReq)
	resp := recorder.Result()
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ErrorResponse{}, err
	}

	if resp.StatusCode == http.StatusOK {
		return models.ErrorResponse{}, fmt.Errorf("200 response: %d %s", resp.StatusCode, string(data))
	}

	var e models.ErrorResponse
	if err = json.Unmarshal(data, &e); err != nil {
		return models.ErrorResponse{}, fmt.Errorf("json unmarshal failed: %w %s", err, string(data))
	}

	return e, nil
}
