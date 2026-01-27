package services

import (
	"bytestream/internal/models"
	"context"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type IdentityProvider interface {
	GetUserInfo(ctx context.Context, token string) (*models.IdentityResponse, error)
}

type AvailabilityProvider interface {
	GetAvailabilityInfo(ctx context.Context, token string, videoID int) (*models.AvailabilityResponse, error)
}
type VideoService struct {
	identityClient     IdentityProvider
	availabilityClient AvailabilityProvider
	s3BaseURL          string
	videoCatalog       map[int]models.VideoMetadata
	logger             *slog.Logger
}

func NewVideoService(identity IdentityProvider, availability AvailabilityProvider, s3BaseURL string, logger *slog.Logger) *VideoService {
	return &VideoService{
		identityClient:     identity,
		availabilityClient: availability,
		s3BaseURL:          s3BaseURL,
		logger:             logger,
		videoCatalog: map[int]models.VideoMetadata{
			46325: {
				Title:    "Example Video 001",
				Filename: "example001",
			},
		},
	}
}

func (s *VideoService) ResolveVideo(ctx context.Context, token string, videoID int) (*models.VideoResponse, error) {
	log := s.logger.With("video_id", videoID)

	metadata, exists := s.videoCatalog[videoID]
	if !exists {
		return nil, models.ErrVideoNotFound
	}

	var identity *models.IdentityResponse
	var availability *models.AvailabilityResponse

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		identity, err = s.identityClient.GetUserInfo(ctx, token)
		if err != nil && ctx.Err() == nil { // only log id we get an error, not if the context is canceled
			s.logUpstreamErr(log, "identity", err)
		}
		return err
	})

	g.Go(func() error {
		var err error
		availability, err = s.availabilityClient.GetAvailabilityInfo(ctx, token, videoID)
		if err != nil && ctx.Err() == nil {
			s.logUpstreamErr(log, "availability", err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if !availability.IsAvailable(time.Now()) {
		return nil, models.ErrVideoUnavailable
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

func (s *VideoService) logUpstreamErr(log *slog.Logger, upstream string, err error) {
	log.Error("upstream_failure", "upstream", upstream, "error", err.Error())
}
