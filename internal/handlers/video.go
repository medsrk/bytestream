package handlers

import (
	"bytestream/internal/models"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type VideoResolver interface {
	ResolveVideo(ctx context.Context, token string, videoID int) (*models.VideoResponse, error)
}

type VideoHandler struct {
	videoService VideoResolver
}

func NewVideoHandler(videoService VideoResolver) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	videoIDStr := r.PathValue("video_id")
	videoID, err := strconv.Atoi(videoIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid video ID")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, "missing auth header")
		return
	}

	token := strings.TrimPrefix(authHeader, "bearer ")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == authHeader {
		writeError(w, http.StatusUnauthorized, "invalid auth header format")
		return
	}

	video, err := h.videoService.ResolveVideo(r.Context(), token, videoID)
	if err != nil {
		log.Printf("Error resolving video: %v", err)

		if errors.Is(err, models.ErrVideoNotFound) {
			writeError(w, http.StatusNotFound, "video not found")
		} else if errors.Is(err, models.ErrVideoUnavailable) {
			writeError(w, http.StatusNotFound, "video not available")
		} else {
			writeError(w, http.StatusServiceUnavailable, "service unavailable")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

func writeError(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: err})
}
