package config

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/mylxsw/go-utils/ternary"
)

var (
	// singleton
	logger *slog.Logger
	once   sync.Once
)

func Logger(w ...io.Writer) *slog.Logger {
	once.Do(func() {
		out := ternary.IfElse[io.Writer](len(w) != 0, io.MultiWriter(w...), os.Stderr)

		logger = slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	})

	return logger
}
