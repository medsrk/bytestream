package services

import (
	"bytestream/internal/models"
	"context"
	"fmt"
	"time"
)

type VideoService struct {
	identityClient     *IdentityClient
	availabilityClient *AvailabilityClient
	s3BaseURL          string
	videoCatalog       map[int]models.VideoMetadata
}

func NewVideoService(identity *IdentityClient, availability *AvailabilityClient, s3BaseURL string) *VideoService {
	return &VideoService{
		identityClient:     identity,
		availabilityClient: availability,
		s3BaseURL:          s3BaseURL,
		videoCatalog: map[int]models.VideoMetadata{
			46325: {
				Title:    "Example Video 001",
				Filename: "example001",
			},
		},
	}
}

func (s *VideoService) ResolveVideo(ctx context.Context, token string, videoID int) (*models.VideoResponse, error) {
	metadata, exists := s.videoCatalog[videoID]
	if !exists {
		return nil, fmt.Errorf("video not found")
	}

	identity, err := s.identityClient.GetUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("getting user info: %w", err)
	}

	availability, err := s.availabilityClient.GetAvailabilityInfo(ctx, token, videoID)
	if err != nil {
		return nil, fmt.Errorf("getting availability: %w", err)
	}

	if !availability.IsAvailable(time.Now()) {
		return nil, fmt.Errorf("video not available")
	}

	filename := metadata.Filename
	if identity.IsPremium() {
		filename += "-premium"
	}

	return &models.VideoResponse{
		VideoID:           videoID,
		Title:             metadata.Title,
		PlaybackBaseURL:   s.s3BaseURL,
		PlaybackFilename:  filename,
		PlaybackExtension: ".mp4",
	}, nil
}
