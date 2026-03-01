package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "steam-deck-stock-alerts",
	Short: "Steam Deck stock monitoring and notification service",
	Long:  "A service that monitors Steam Deck OLED inventory and sends notifications when stock status changes.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
