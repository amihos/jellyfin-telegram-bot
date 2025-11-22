package database

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"jellyfin-telegram-bot/pkg/models"

	"gorm.io/gorm"
)

// AddMutedSeries adds a new muted series for a user
func (db *DB) AddMutedSeries(chatID int64, seriesID string, seriesName string) error {
	mutedSeries := models.MutedSeries{
		ChatID:     chatID,
		SeriesID:   seriesID,
		SeriesName: seriesName,
	}

	result := db.Create(&mutedSeries)
	if result.Error != nil {
		// Handle duplicate constraint violations gracefully
		// SQLite returns "UNIQUE constraint failed" error message
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			slog.Debug("Series already muted", "chat_id", chatID, "series_id", seriesID)
			return nil
		}
		slog.Error("Failed to add muted series", "chat_id", chatID, "series_id", seriesID, "error", result.Error)
		return fmt.Errorf("failed to add muted series: %w", result.Error)
	}

	slog.Info("Added muted series", "chat_id", chatID, "series_id", seriesID, "series_name", seriesName)
	return nil
}

// RemoveMutedSeries removes a muted series for a user
func (db *DB) RemoveMutedSeries(chatID int64, seriesID string) error {
	result := db.Where("chat_id = ? AND series_id = ?", chatID, seriesID).Delete(&models.MutedSeries{})

	if result.Error != nil {
		slog.Error("Failed to remove muted series", "chat_id", chatID, "series_id", seriesID, "error", result.Error)
		return fmt.Errorf("failed to remove muted series: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		slog.Debug("No muted series found to remove", "chat_id", chatID, "series_id", seriesID)
		return gorm.ErrRecordNotFound
	}

	slog.Info("Removed muted series", "chat_id", chatID, "series_id", seriesID)
	return nil
}

// GetMutedSeriesByUser returns all muted series for a user
func (db *DB) GetMutedSeriesByUser(chatID int64) ([]models.MutedSeries, error) {
	var mutedSeries []models.MutedSeries
	result := db.Where("chat_id = ?", chatID).Find(&mutedSeries)

	if result.Error != nil {
		slog.Error("Failed to get muted series", "chat_id", chatID, "error", result.Error)
		return nil, fmt.Errorf("failed to get muted series: %w", result.Error)
	}

	slog.Debug("Retrieved muted series", "chat_id", chatID, "count", len(mutedSeries))
	return mutedSeries, nil
}

// IsSeriesMuted checks if a series is muted for a user
func (db *DB) IsSeriesMuted(chatID int64, seriesID string) (bool, error) {
	var count int64
	result := db.Model(&models.MutedSeries{}).
		Where("chat_id = ? AND series_id = ?", chatID, seriesID).
		Count(&count)

	if result.Error != nil {
		slog.Error("Failed to check if series is muted", "chat_id", chatID, "series_id", seriesID, "error", result.Error)
		return false, fmt.Errorf("failed to check if series is muted: %w", result.Error)
	}

	return count > 0, nil
}
