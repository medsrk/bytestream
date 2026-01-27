package services_test

import (
	"bytestream/internal/models"
	"bytestream/internal/services"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockIdentityProvider struct {
	response *models.IdentityResponse
	err      error
}

func (m *mockIdentityProvider) GetUserInfo(ctx context.Context, token string) (*models.IdentityResponse, error) {
	return m.response, m.err
}

type mockAvailabilityProvider struct {
	response *models.AvailabilityResponse
	err      error
}

func (m *mockAvailabilityProvider) GetAvailabilityInfo(ctx context.Context, token string, videoID int) (*models.AvailabilityResponse, error) {
	return m.response, m.err
}

func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func availableWindow() models.AvailabilityWindow {
	return models.AvailabilityWindow{
		From: "2020-01-01",
		To:   "2030-12-31",
	}
}

func expiredWindow() models.AvailabilityWindow {
	return models.AvailabilityWindow{
		From: "2020-01-01",
		To:   "2020-12-31",
	}
}

func TestVideoService_ResolveVideo(t *testing.T) {
	t.Run("premium user gets premium video", func(t *testing.T) {
		svc := services.NewVideoService(
			&mockIdentityProvider{
				response: &models.IdentityResponse{ID: 1, Roles: []string{"premium"}},
			},
			&mockAvailabilityProvider{
				response: &models.AvailabilityResponse{
					VideoID:            46325,
					AvailabilityWindow: availableWindow(),
				},
			},
			"https://s3.example.com",
			noopLogger(),
		)
		got, err := svc.ResolveVideo(context.Background(), "token", 46325)

		require.NoError(t, err)
		assert.Equal(t, 46325, got.VideoID)
		assert.Equal(t, "example001-premium", got.PlaybackFilename)
	})

	t.Run("non-premium user gets regular filename", func(t *testing.T) {
		svc := services.NewVideoService(
			&mockIdentityProvider{
				response: &models.IdentityResponse{ID: 1, Roles: []string{"basic"}},
			},
			&mockAvailabilityProvider{
				response: &models.AvailabilityResponse{
					VideoID:            46325,
					AvailabilityWindow: availableWindow(),
				},
			},
			"https://s3.example.com",
			noopLogger(),
		)

		got, err := svc.ResolveVideo(context.Background(), "token", 46325)

		require.NoError(t, err)
		assert.Equal(t, "example001", got.PlaybackFilename)
	})

	t.Run("video not in catalog returns ErrVideoNotFound", func(t *testing.T) {
		svc := services.NewVideoService(
			&mockIdentityProvider{},
			&mockAvailabilityProvider{},
			"https://s3.example.com",
			noopLogger(),
		)

		_, err := svc.ResolveVideo(context.Background(), "token", 99999)

		assert.ErrorIs(t, err, models.ErrVideoNotFound)
	})

	t.Run("outside availability window returns ErrVideoUnavailable", func(t *testing.T) {
		svc := services.NewVideoService(
			&mockIdentityProvider{
				response: &models.IdentityResponse{ID: 1},
			},
			&mockAvailabilityProvider{
				response: &models.AvailabilityResponse{
					VideoID:            46325,
					AvailabilityWindow: expiredWindow(),
				},
			},
			"https://s3.example.com",
			noopLogger(),
		)

		_, err := svc.ResolveVideo(context.Background(), "token", 46325)

		assert.ErrorIs(t, err, models.ErrVideoUnavailable)
	})

	t.Run("identity service failure returns error", func(t *testing.T) {
		svc := services.NewVideoService(
			&mockIdentityProvider{
				err: errors.New("connection refused"),
			},
			&mockAvailabilityProvider{
				response: &models.AvailabilityResponse{
					VideoID:            46325,
					AvailabilityWindow: availableWindow(),
				},
			},
			"https://s3.example.com",
			noopLogger(),
		)

		_, err := svc.ResolveVideo(context.Background(), "token", 46325)

		assert.Error(t, err)
		assert.NotErrorIs(t, err, models.ErrVideoNotFound)
		assert.NotErrorIs(t, err, models.ErrVideoUnavailable)
	})

	t.Run("availability service failure returns error", func(t *testing.T) {
		svc := services.NewVideoService(
			&mockIdentityProvider{
				response: &models.IdentityResponse{ID: 1},
			},
			&mockAvailabilityProvider{
				err: errors.New("timeout"),
			},
			"https://s3.example.com",
			noopLogger(),
		)

		_, err := svc.ResolveVideo(context.Background(), "token", 46325)

		assert.Error(t, err)
	})
}
