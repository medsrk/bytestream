package models

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrVideoNotFound    = errors.New("video not found")
	ErrVideoUnavailable = errors.New("video not available")
)
