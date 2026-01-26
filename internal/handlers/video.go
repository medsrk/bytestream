package handlers

import (
	"bytestream/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type VideoHandler struct {
	videoService *services.VideoService
}

func NewVideoHandler(videoService *services.VideoService) *VideoHandler {
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

		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "video not found")
		} else if strings.Contains(err.Error(), "not available") {
			writeError(w, http.StatusNotFound, "video not available")
		} else {
			writeError(w, http.StatusServiceUnavailable, "service unavailable")
		}
		return
	}

	json.NewEncoder(w).Encode(video)
}

func writeError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}
