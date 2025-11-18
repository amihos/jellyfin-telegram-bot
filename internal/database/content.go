package database

import (
	"fmt"
	"jellyfin-telegram-bot/pkg/models"
)

// IsContentNotified checks if content has already been notified
func (db *DB) IsContentNotified(jellyfinID string) (bool, error) {
	var count int64
	result := db.Model(&models.ContentCache{}).
		Where("jellyfin_id = ?", jellyfinID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check content notification status: %w", result.Error)
	}

	return count > 0, nil
}

// MarkContentNotified marks content as notified by storing it in the cache
func (db *DB) MarkContentNotified(jellyfinID, title, contentType string) error {
	content := models.ContentCache{
		JellyfinID: jellyfinID,
		Title:      title,
		Type:       contentType,
	}

	result := db.Create(&content)
	if result.Error != nil {
		return fmt.Errorf("failed to mark content as notified: %w", result.Error)
	}

	return nil
}
