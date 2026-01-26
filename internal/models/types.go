package models

type IdentityResponse struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type AvailabilityResponse struct {
	VideoID            int                `json:"video_id"`
	AvailabilityWindow AvailabilityWindow `json:"availability_window"`
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
