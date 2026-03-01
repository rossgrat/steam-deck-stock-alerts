package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rossgrat/steam-deck-stock-alerts/internal/config"
)

func New(cfg config.LogConfig) *slog.Logger {
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
	}

	multi := io.MultiWriter(os.Stdout, fileWriter)

	return slog.New(slog.NewJSONHandler(multi, nil))
}
