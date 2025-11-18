package main

import (
	"log/slog"
	"os"

	"jellyfin-telegram-bot/internal/database"
	"jellyfin-telegram-bot/internal/handlers"
)

// Example demonstrating how to start the webhook server
func main() {
	// Initialize database
	db, err := database.NewDB("./bot.db")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Get webhook secret from environment (optional)
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		slog.Warn("No WEBHOOK_SECRET set - webhook endpoint is unprotected")
	}

	// Create webhook handler
	webhookHandler := handlers.NewWebhookHandler(db, webhookSecret)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start webhook server (blocking call)
	slog.Info("Starting webhook server", "port", port)
	if err := handlers.StartWebhookServer(port, webhookHandler); err != nil {
		slog.Error("Webhook server error", "error", err)
		os.Exit(1)
	}
}
