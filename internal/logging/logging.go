package logging

import (
	"log/slog"
	"os"
)

func NewLogger(env string, serviceName string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler).With("service", serviceName)
}
