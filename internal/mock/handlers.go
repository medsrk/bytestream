package mock

import (
	"encoding/json"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /indentity/userinfo", handleUserInfo)
	mux.HandleFunc("GET /availability/availabilityinfo/46325", handleAvailability)
}

func handleUserInfo(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"id":    7654,
		"name":  "Joe Bloggs",
		"email": "joe.bloggs@noemail.notreal",
		"roles": []string{"premium"},
	}

	json.NewEncoder(w).Encode(resp)
}

func handleAvailability(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"video_id": 46325,
		"availability_window": map[string]string{
			"from": "2025-11-05",
			"to":   "2026-05-04",
		},
	}

	json.NewEncoder(w).Encode(resp)
}
