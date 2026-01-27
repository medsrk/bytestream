package main

import (
	"bytestream/internal/config"
	"bytestream/internal/handlers"
	"bytestream/internal/logging"
	"bytestream/internal/mock"
	"bytestream/internal/services"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var port = ":8080"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	httpClient := &http.Client{
		Timeout: cfg.HTTPTimeout,
	}

	identityClient := services.NewIdentityClient(cfg.IdentityURL, httpClient)
	availabilityClient := services.NewAvailabilityClient(cfg.AvailabilityURL, httpClient)

	logger := logging.NewLogger(cfg.Env, "bytestream_video")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	videoService := services.NewVideoService(
		identityClient,
		availabilityClient,
		cfg.S3BaseURL,
		logger,
	)

	videoHandler := handlers.NewVideoHandler(videoService)

	mux := http.NewServeMux()

	mock.RegisterRoutes(mux)

	mux.HandleFunc("GET /video/{video_id}", videoHandler.GetVideo)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}

	go func() {
		logger.Info("server_start", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server_error", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()

	logger.Info("server_shutdown", "status", "shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server_shutdown", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("server_shutdown", "status", "complete")
}
