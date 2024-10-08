package config

import (
	"log/slog"
)


var (
	// singleton
	logger *slog.Logger = slog.Default()
)

func Logger() *slog.Logger {
	return logger
}