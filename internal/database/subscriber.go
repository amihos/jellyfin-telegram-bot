package database

import (
	"fmt"
	"jellyfin-telegram-bot/pkg/models"

	"gorm.io/gorm"
)

// AddSubscriber adds a new subscriber to the database
func (db *DB) AddSubscriber(chatID int64, username, firstName string) error {
	subscriber := models.Subscriber{
		ChatID:    chatID,
		Username:  username,
		FirstName: firstName,
		IsActive:  true,
	}

	// Use FirstOrCreate to handle duplicate chat_id gracefully
	result := db.Where(models.Subscriber{ChatID: chatID}).FirstOrCreate(&subscriber)
	if result.Error != nil {
		return fmt.Errorf("failed to add subscriber: %w", result.Error)
	}

	// If subscriber was found (not created), ensure they're active
	if result.RowsAffected == 0 {
		// Subscriber already exists, reactivate if needed
		if err := db.Model(&subscriber).Update("is_active", true).Error; err != nil {
			return fmt.Errorf("failed to reactivate subscriber: %w", err)
		}
	}

	return nil
}

// RemoveSubscriber removes or deactivates a subscriber
func (db *DB) RemoveSubscriber(chatID int64) error {
	result := db.Model(&models.Subscriber{}).
		Where("chat_id = ?", chatID).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to remove subscriber: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetAllActiveSubscribers returns a list of all active subscriber chat IDs
func (db *DB) GetAllActiveSubscribers() ([]int64, error) {
	var subscribers []models.Subscriber
	result := db.Where("is_active = ?", true).Find(&subscribers)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active subscribers: %w", result.Error)
	}

	chatIDs := make([]int64, 0, len(subscribers))
	for _, sub := range subscribers {
		chatIDs = append(chatIDs, sub.ChatID)
	}

	return chatIDs, nil
}

// IsSubscribed checks if a user is subscribed and active
func (db *DB) IsSubscribed(chatID int64) (bool, error) {
	var count int64
	result := db.Model(&models.Subscriber{}).
		Where("chat_id = ? AND is_active = ?", chatID, true).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check subscription status: %w", result.Error)
	}

	return count > 0, nil
}
