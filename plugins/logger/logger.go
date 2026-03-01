package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/lumberjack.v2"

	"github.com/rossgrat/steam-deck-stock-alerts/internal/config"
)

func New(cfg config.LogConfig) zerolog.Logger {
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
	}

	multi := io.MultiWriter(os.Stdout, fileWriter)

	return zerolog.New(multi).With().Timestamp().Logger()
}
