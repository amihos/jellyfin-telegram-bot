package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Telegram TelegramConfig
	Jellyfin JellyfinConfig
	Webhook  WebhookConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

// TelegramConfig holds Telegram bot configuration
type TelegramConfig struct {
	BotToken string
}

// JellyfinConfig holds Jellyfin server configuration
type JellyfinConfig struct {
	ServerURL string
	APIKey    string
}

// WebhookConfig holds webhook server configuration
type WebhookConfig struct {
	Secret string
	Port   int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Telegram: TelegramConfig{
			BotToken: getEnvRequired("TELEGRAM_BOT_TOKEN"),
		},
		Jellyfin: JellyfinConfig{
			ServerURL: getEnvRequired("JELLYFIN_SERVER_URL"),
			APIKey:    getEnvRequired("JELLYFIN_API_KEY"),
		},
		Webhook: WebhookConfig{
			Secret: getEnv("WEBHOOK_SECRET", ""),
			Port:   getEnvInt("PORT", 8080),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./bot.db"),
		},
		Logger: GetLoggerFromEnv(),
	}

	// Validate required fields
	if config.Telegram.BotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if config.Jellyfin.ServerURL == "" {
		return nil, fmt.Errorf("JELLYFIN_SERVER_URL is required")
	}
	if config.Jellyfin.APIKey == "" {
		return nil, fmt.Errorf("JELLYFIN_API_KEY is required")
	}

	return config, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvRequired gets a required environment variable
func getEnvRequired(key string) string {
	return os.Getenv(key)
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
