package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bytestream/internal/handlers"
	"bytestream/internal/models"
)

type mockVideoResolver struct {
	response *models.VideoResponse
	err      error
}

func (m *mockVideoResolver) ResolveVideo(ctx context.Context, token string, videoID int) (*models.VideoResponse, error) {
	return m.response, m.err
}

func TestVideoHandler_GetVideo(t *testing.T) {
	t.Run("success returns 200 with video response", func(t *testing.T) {
		handler := handlers.NewVideoHandler(&mockVideoResolver{
			response: &models.VideoResponse{
				VideoID:          46325,
				Title:            "Test Video",
				PlaybackFilename: "test",
			},
		})

		rec := doRequest(t, handler, "/video/46325", "Bearer valid-token")

		assert.Equal(t, http.StatusOK, rec.Code)

		var got models.VideoResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
		assert.Equal(t, 46325, got.VideoID)
		assert.Equal(t, "Test Video", got.Title)
	})

	t.Run("invalid video id returns 400", func(t *testing.T) {
		handler := handlers.NewVideoHandler(&mockVideoResolver{})

		rec := doRequest(t, handler, "/video/not-a-number", "Bearer token")

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assertErrorResponse(t, rec, "invalid video ID")
	})

	t.Run("missing auth header returns 401", func(t *testing.T) {
		handler := handlers.NewVideoHandler(&mockVideoResolver{})

		rec := doRequest(t, handler, "/video/46325", "")

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assertErrorResponse(t, rec, "missing auth header")
	})

	t.Run("invalid auth format returns 401", func(t *testing.T) {
		handler := handlers.NewVideoHandler(&mockVideoResolver{})

		rec := doRequest(t, handler, "/video/46325", "InvalidFormat token")

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assertErrorResponse(t, rec, "invalid auth header format")
	})

	t.Run("video not found returns 404", func(t *testing.T) {
		handler := handlers.NewVideoHandler(&mockVideoResolver{
			err: models.ErrVideoNotFound,
		})

		rec := doRequest(t, handler, "/video/99999", "Bearer token")

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assertErrorResponse(t, rec, "video not found")
	})

	t.Run("video unavailable returns 404", func(t *testing.T) {
		handler := handlers.NewVideoHandler(&mockVideoResolver{
			err: models.ErrVideoUnavailable,
		})

		rec := doRequest(t, handler, "/video/46325", "Bearer token")

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assertErrorResponse(t, rec, "video not available")
	})

	t.Run("unknown error returns 503", func(t *testing.T) {
		handler := handlers.NewVideoHandler(&mockVideoResolver{
			err: errors.New("upstream failure"),
		})

		rec := doRequest(t, handler, "/video/46325", "Bearer token")

		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assertErrorResponse(t, rec, "service unavailable")
	})
}

func doRequest(t *testing.T, handler *handlers.VideoHandler, path, authHeader string) *httptest.ResponseRecorder {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /video/{video_id}", handler.GetVideo)

	req := httptest.NewRequest(http.MethodGet, path, nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedError string) {
	t.Helper()

	var errResp models.ErrorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&errResp))
	assert.Equal(t, expectedError, errResp.Error)
}
