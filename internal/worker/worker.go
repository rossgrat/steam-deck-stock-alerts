package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/rossgrat/steam-deck-stock-alerts/internal/service"
)

type Worker struct {
	service  *service.StockService
	interval time.Duration
	logger   *slog.Logger
}

func New(svc *service.StockService, interval time.Duration, logger *slog.Logger) *Worker {
	return &Worker{
		service:  svc,
		interval: interval,
		logger:   logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Info("worker started", "interval", w.interval.String())

	// Run immediately on startup
	w.service.CheckAndNotify()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopping")
			return
		case <-ticker.C:
			w.service.CheckAndNotify()
		}
	}
}
