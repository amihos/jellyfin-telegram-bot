package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"jellyfin-telegram-bot/internal/config"
	"jellyfin-telegram-bot/internal/database"
	"jellyfin-telegram-bot/internal/handlers"
	"jellyfin-telegram-bot/internal/jellyfin"
	"jellyfin-telegram-bot/internal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up structured logging with file rotation
	logger := config.SetupLogger(cfg.Logger)
	slog.SetDefault(logger)

	// Create context with cancellation for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	slog.Info("Jellyfin Telegram Bot starting...",
		"version", "0.1.0",
		"port", cfg.Webhook.Port,
		"database", cfg.Database.Path,
	)

	// Initialize database connection
	db, err := database.NewDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	slog.Info("Database initialized", "path", cfg.Database.Path)

	// Initialize Jellyfin API client
	jellyfinClient := jellyfin.NewClient(cfg.Jellyfin.ServerURL, cfg.Jellyfin.APIKey)
	slog.Info("Jellyfin client initialized", "server", cfg.Jellyfin.ServerURL)

	// Create Jellyfin client adapter for the bot
	jellyfinAdapter := telegram.NewJellyfinClientAdapter(jellyfinClient)

	// Initialize Telegram bot
	bot, err := telegram.NewBot(cfg.Telegram.BotToken, db, jellyfinAdapter)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	slog.Info("Telegram bot initialized")

	// Create broadcaster adapter for webhook handler
	broadcaster := telegram.NewBroadcasterAdapter(bot)

	// Initialize webhook handler
	webhookHandler := handlers.NewWebhookHandler(db, cfg.Webhook.Secret)
	webhookHandler.SetBroadcaster(broadcaster)
	slog.Info("Webhook handler initialized")

	// Start webhook server in goroutine
	go func() {
		slog.Info("Starting webhook server", "port", cfg.Webhook.Port)
		port := fmt.Sprintf("%d", cfg.Webhook.Port)
		if err := handlers.StartWebhookServer(port, webhookHandler); err != nil {
			slog.Error("Webhook server failed", "error", err)
			cancel()
		}
	}()

	// Start bot polling
	go bot.Start(ctx)

	slog.Info("Bot is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("Received shutdown signal, shutting down gracefully...")

	// Cleanup happens automatically when context is cancelled

	slog.Info("Shutdown complete")
}
