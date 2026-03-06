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
	apiHealthy  bool
	errorCount  int
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
		apiHealthy:  true,
		errorCount:  0,
	}
}

func (s *StockService) CheckAndNotify() error {
	var anyError bool

	for _, pkg := range s.packages {
		if err := s.checkPackage(pkg); err != nil {
			anyError = true
			s.logger.Error("failed to check package",
				"package_id", pkg.ID,
				"package_name", pkg.Name,
				"error", err,
			)
		}
	}

	// If we received an error and the API is healthy, increment the error count
	if anyError && s.apiHealthy {
		s.errorCount++
	}

	// If we reach 10 errors and the API is healthy, set to unhealthy and send an error notification
	if s.errorCount >= 10 && s.apiHealthy {
		s.sendAPIErrorNotification()
		s.apiHealthy = false
		s.errorCount = 0
	}

	// If we don't get an error, and the API is unhealthy, send a recovery notifcation
	if !anyError && !s.apiHealthy {
		s.apiHealthy = true
		s.sendAPIRecoveredNotification()
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

func (s *StockService) sendAPIErrorNotification() error {
	s.logger.Info("sending API error notification")

	return s.ntfyClient.Send(ntfy.Notification{
		Title:    "Steam Deck Stock Alerts — API Error",
		Body:     "The Steam API is returning errors. Stock checks may be unreliable until this is resolved.",
		Priority: 4,
		Tags:     []string{"warning"},
	})
}

func (s *StockService) sendAPIRecoveredNotification() error {
	s.logger.Info("sending API recovered notifcation")

	return s.ntfyClient.Send(ntfy.Notification{
		Title:    "Steam Deck Stock Alerts — API Recovered",
		Body:     "The Steam API is responding normally again.",
		Priority: 3,
		Tags:     []string{"white_check_mark"},
	})
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
