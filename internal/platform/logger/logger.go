package logger

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

func New(service, env, level string) (*slog.Logger, error) {
	minLevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: minLevel})
	return slog.New(jsonHandler).With(
		"service", service,
		"env", env,
	), nil
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return (slog.Level)(-8), errors.New("wrong log level")
	}
}
