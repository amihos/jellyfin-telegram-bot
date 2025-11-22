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
	Testing  TestingConfig
}

// TestingConfig holds testing and feature flag configuration
type TestingConfig struct {
	TesterChatIDs      []int64 // Chat IDs that can access beta features
	EnableBetaFeatures bool    // Global flag to enable/disable beta features
	NotifyOnlyTesters  bool    // If true, send ALL notifications only to testers (for debugging)
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
		Testing: TestingConfig{
			TesterChatIDs:      getEnvInt64Slice("TESTER_CHAT_IDS", []int64{}),
			EnableBetaFeatures: getEnvBool("ENABLE_BETA_FEATURES", false),
			NotifyOnlyTesters:  getEnvBool("NOTIFY_ONLY_TESTERS", false),
		},
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

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvInt64Slice gets a comma-separated list of int64 values with a default
func getEnvInt64Slice(key string, defaultValue []int64) []int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result []int64
	for _, part := range splitAndTrim(value, ",") {
		if intValue, err := strconv.ParseInt(part, 10, 64); err == nil {
			result = append(result, intValue)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}

// splitAndTrim splits a string and trims whitespace from each part
func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitString splits a string by separator
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	var parts []string
	var current string
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	parts = append(parts, current)
	return parts
}

// trimSpace removes leading and trailing whitespace
func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}

// IsTester checks if a chat ID is in the tester allowlist
func (c *Config) IsTester(chatID int64) bool {
	if !c.Testing.EnableBetaFeatures {
		return false
	}

	for _, testerID := range c.Testing.TesterChatIDs {
		if testerID == chatID {
			return true
		}
	}
	return false
}
