package service

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/rossgrat/steam-deck-stock-alerts/internal/config"
	"github.com/rossgrat/steam-deck-stock-alerts/internal/repo"
	"github.com/rossgrat/steam-deck-stock-alerts/plugins/ntfy"
	"github.com/rossgrat/steam-deck-stock-alerts/plugins/steam"
)

const steamDeckStoreURL = "https://store.steampowered.com/steamdeck"

type StockService struct {
	steamClient *steam.Client
	repo        *repo.StockRepo
	ntfyClient  *ntfy.Client
	logger      *slog.Logger
	packages    []config.PackageConfig
	countryCode string
}

func NewStockService(
	steamClient *steam.Client,
	repo *repo.StockRepo,
	ntfyClient *ntfy.Client,
	logger *slog.Logger,
	packages []config.PackageConfig,
	countryCode string,
) *StockService {
	return &StockService{
		steamClient: steamClient,
		repo:        repo,
		ntfyClient:  ntfyClient,
		logger:      logger,
		packages:    packages,
		countryCode: countryCode,
	}
}

func (s *StockService) CheckAndNotify() error {
	for _, pkg := range s.packages {
		if err := s.checkPackage(pkg); err != nil {
			s.logger.Error("failed to check package",
				"package_id", pkg.ID,
				"package_name", pkg.Name,
				"error", err,
			)
		}
	}
	return nil
}

func (s *StockService) checkPackage(pkg config.PackageConfig) error {
	inventory, err := s.steamClient.CheckInventory(pkg.ID, s.countryCode)
	if err != nil {
		return fmt.Errorf("checking inventory: %w", err)
	}

	packageID := strconv.Itoa(pkg.ID)
	previousState, err := s.repo.GetState(packageID)
	if err != nil {
		return fmt.Errorf("getting previous state: %w", err)
	}

	currentlyAvailable := inventory.InventoryAvailable

	s.logger.Info("stock check completed",
		"package_id", pkg.ID,
		"package_name", pkg.Name,
		"available", currentlyAvailable,
		"high_pending_orders", inventory.HighPendingOrders,
	)

	if err := s.handleTransition(pkg, previousState, currentlyAvailable); err != nil {
		return fmt.Errorf("handling transition: %w", err)
	}

	if err := s.repo.SetState(packageID, currentlyAvailable); err != nil {
		return fmt.Errorf("setting state: %w", err)
	}

	return nil
}

func (s *StockService) handleTransition(pkg config.PackageConfig, previousState *bool, currentlyAvailable bool) error {
	if previousState == nil {
		// First run
		if currentlyAvailable {
			return s.sendInStockNotification(pkg)
		}
		return nil
	}

	wasAvailable := *previousState

	if !wasAvailable && currentlyAvailable {
		return s.sendInStockNotification(pkg)
	}

	if wasAvailable && !currentlyAvailable {
		return s.sendOutOfStockNotification(pkg)
	}

	return nil
}

func (s *StockService) sendInStockNotification(pkg config.PackageConfig) error {
	s.logger.Info("sending in-stock notification",
		"package_id", pkg.ID,
		"package_name", pkg.Name,
	)

	return s.ntfyClient.Send(ntfy.Notification{
		Title:    fmt.Sprintf("Steam Deck %s — In Stock!", pkg.Name),
		Body:     fmt.Sprintf("The Steam Deck %s is available for purchase! Go grab one before it sells out.", pkg.Name),
		Priority: 5,
		Tags:     []string{"rotating_light", "video_game"},
		ClickURL: steamDeckStoreURL,
	})
}

func (s *StockService) sendOutOfStockNotification(pkg config.PackageConfig) error {
	s.logger.Info("sending out-of-stock notification",
		"package_id", pkg.ID,
		"package_name", pkg.Name,
	)

	return s.ntfyClient.Send(ntfy.Notification{
		Title:    fmt.Sprintf("Steam Deck %s — Out of Stock", pkg.Name),
		Body:     fmt.Sprintf("The Steam Deck %s is no longer available.", pkg.Name),
		Priority: 3,
		Tags:     []string{"video_game"},
		ClickURL: steamDeckStoreURL,
	})
}
