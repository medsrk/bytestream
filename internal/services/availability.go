package services

import (
	"bytestream/internal/models"
	"context"
	"fmt"
	"net/http"
)

type AvailabilityClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAvailabilityClient(baseURL string, httpClient *http.Client) *AvailabilityClient {
	return &AvailabilityClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *AvailabilityClient) GetAvailabilityInfo(ctx context.Context, token string, videoID int) (*models.AvailabilityResponse, error) {
	url := fmt.Sprintf("%s/availability/availabilityinfo/%d", c.baseURL, videoID)

	var availability models.AvailabilityResponse
	if err := doAuthenticatedRequest(ctx, c.httpClient, url, token, &availability); err != nil {
		return nil, fmt.Errorf("availability service: %w", err)
	}

	return &availability, nil
}
