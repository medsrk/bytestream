package models

import "time"

type IdentityResponse struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// IsPremium checks if the user has the premium role
func (i *IdentityResponse) IsPremium() bool {
	for _, role := range i.Roles {
		if role == "premium" {
			return true
		}
	}
	return false
}

type AvailabilityResponse struct {
	VideoID            int                `json:"video_id"`
	AvailabilityWindow AvailabilityWindow `json:"availability_window"`
}

// IsAvailable checks if the video is available in a given window
// TODO: Currently uses UTC for date comparison. This may not align with user timezones.
// Product decision needed: should availability windows be UTC-based or user-timezone-based?
func (a *AvailabilityResponse) IsAvailable(at time.Time) bool {
	from, err := time.Parse("2006-01-02", a.AvailabilityWindow.From)
	if err != nil {
		return false
	}

	to, err := time.Parse("2006-01-02", a.AvailabilityWindow.To)
	if err != nil {
		return false
	}

	// Check if 'at' is within the window using the beginning of the day
	atDay := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.UTC)

	return !atDay.Before(from) && !atDay.After(to)
}

type AvailabilityWindow struct {
	From string `json:"from"` // YYYY-MM-DD format
	To   string `json:"to"`
}

type VideoResponse struct {
	VideoID           int    `json:"video_id"`
	Title             string `json:"title"`
	PlaybackBaseURL   string `json:"playback_baseurl"`
	PlaybackFilename  string `json:"playback_filename"`
	PlaybackExtension string `json:"playback_extension"`
}

type VideoMetadata struct {
	Title    string
	Filename string
}
