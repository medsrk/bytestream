package models_test

import (
	"bytestream/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIdentityResponse_IsPremium(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  bool
	}{
		{"has premium", []string{"basic", "premium"}, true},
		{"only premium", []string{"premium"}, true},
		{"empty roles", []string{}, false},
		{"nil roles", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &models.IdentityResponse{Roles: tt.roles}
			assert.Equal(t, tt.want, i.IsPremium())
		})
	}
}

func TestAvailabilityResponse_IsAvailable(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		at   time.Time
		want bool
	}{
		{
			name: "within window",
			from: "2026-01-01",
			to:   "2026-12-31",
			at:   time.Date(2026, 01, 27, 19, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "before window",
			from: "2026-01-01",
			to:   "2026-12-31",
			at:   time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "after window",
			from: "2025-01-01",
			to:   "2025-12-31",
			at:   time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "on start date",
			from: "2026-01-01",
			to:   "2026-12-31",
			at:   time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "on end date",
			from: "2026-01-01",
			to:   "2026-12-31",
			at:   time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC),
			want: true,
		},
		{
			name: "invalid from date",
			from: "invalid",
			to:   "2026-12-31",
			at:   time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "invalid to date",
			from: "2026-01-01",
			to:   "invalid",
			at:   time.Date(2026, 01, 01, 0, 0, 0, 0, time.UTC),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &models.AvailabilityResponse{
				AvailabilityWindow: models.AvailabilityWindow{
					From: tt.from,
					To:   tt.to,
				},
			}
			assert.Equal(t, tt.want, a.IsAvailable(tt.at))
		})
	}
}
