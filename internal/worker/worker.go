package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/rossgrat/steam-deck-stock-alerts/internal/service"
)

type Worker struct {
	service  *service.StockService
	interval time.Duration
	logger   zerolog.Logger
}

func New(svc *service.StockService, interval time.Duration, logger zerolog.Logger) *Worker {
	return &Worker{
		service:  svc,
		interval: interval,
		logger:   logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Info().Dur("interval", w.interval).Msg("worker started")

	// Run immediately on startup
	w.service.CheckAndNotify()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg("worker stopping")
			return
		case <-ticker.C:
			w.service.CheckAndNotify()
		}
	}
}
