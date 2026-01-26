package main

import (
	"bytestream/internal/handlers"
	"bytestream/internal/mock"
	"bytestream/internal/services"
	"log"
	"net/http"
	"time"
)

var port = ":8080"

func main() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	identityClient := services.NewIdentityClient("http://localhost"+port, httpClient)
	availabilityClient := services.NewAvailabilityClient("http://localhost"+port, httpClient)

	videoService := services.NewVideoService(
		identityClient,
		availabilityClient,
		"https://s3.eu-west-1.amazon.com/bytestreamfake",
	)

	videoHandler := handlers.NewVideoHandler(videoService)

	mux := http.NewServeMux()

	mock.RegisterRoutes(mux)

	mux.HandleFunc("GET /video/{video_id}", videoHandler.GetVideo)

	log.Printf("server listening on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
