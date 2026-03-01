package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/rossgrat/steam-deck-stock-alerts/internal/config"
	"github.com/rossgrat/steam-deck-stock-alerts/internal/repo"
	"github.com/rossgrat/steam-deck-stock-alerts/internal/service"
	"github.com/rossgrat/steam-deck-stock-alerts/internal/worker"
	"github.com/rossgrat/steam-deck-stock-alerts/plugins/logger"
	"github.com/rossgrat/steam-deck-stock-alerts/plugins/ntfy"
	"github.com/rossgrat/steam-deck-stock-alerts/plugins/steam"
)

var configPath string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the stock monitoring service",
	RunE:  runStart,
}

func init() {
	startCmd.Flags().StringVar(&configPath, "config", "./config.yaml", "path to config file")
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log := logger.New(cfg.Log)

	stockRepo, err := repo.NewStockRepo(cfg.DB.Path)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer stockRepo.Close()

	steamClient := steam.NewClient()
	ntfyClient := ntfy.NewClient(cfg.Ntfy.URL, cfg.Ntfy.Topic, cfg.Ntfy.Token)

	svc := service.NewStockService(
		steamClient,
		stockRepo,
		ntfyClient,
		log,
		cfg.Packages,
		cfg.CountryCode,
	)

	// Send startup notification
	log.Info("service starting")
	if err := ntfyClient.Send(ntfy.Notification{
		Title:    "Steam Deck Stock Alerts",
		Body:     "Stock monitoring service started.",
		Priority: 3,
		Tags:     []string{"rocket"},
	}); err != nil {
		log.Error("failed to send startup notification", "error", err)
	}

	// Set up signal handling
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run the worker (blocks until context is cancelled)
	w := worker.New(svc, cfg.PollingInterval, log)
	w.Run(ctx)

	// Send shutdown notification
	log.Info("service stopping")
	if err := ntfyClient.Send(ntfy.Notification{
		Title:    "Steam Deck Stock Alerts",
		Body:     "Stock monitoring service stopping.",
		Priority: 3,
		Tags:     []string{"stop_sign"},
	}); err != nil {
		log.Error("failed to send shutdown notification", "error", err)
	}

	return nil
}
