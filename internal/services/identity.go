package services

import (
	"bytestream/internal/models"
	"context"
	"fmt"
	"net/http"
)

type IdentityClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewIdentityClient(baseURL string, httpClient *http.Client) *IdentityClient {
	return &IdentityClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *IdentityClient) GetUserInfo(ctx context.Context, token string) (*models.IdentityResponse, error) {
	url := c.baseURL + "/identity/userinfo"

	var identity models.IdentityResponse
	if err := doAuthenticatedRequest(ctx, c.httpClient, url, token, &identity); err != nil {
		return nil, fmt.Errorf("identity service: %w", err)
	}

	return &identity, nil
}
