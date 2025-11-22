package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	botModels "github.com/go-telegram/bot/models"
)

// NotificationContent represents content to be broadcasted
type NotificationContent struct {
	ItemID        string
	Type          string // "Movie" or "Episode"
	Title         string
	Overview      string
	Year          int
	Rating        float64
	SeriesName    string
	SeasonNumber  int
	EpisodeNumber int
}

// FormatNotification formats content for notification message
func FormatNotification(content *NotificationContent) string {
	var message strings.Builder

	if content.Type == "Movie" {
		// Movie notification format
		message.WriteString("🎬 فیلم جدید\n\n")
		message.WriteString(fmt.Sprintf("نام: %s", content.Title))

		if content.Year > 0 {
			message.WriteString(fmt.Sprintf("\nسال: %d", content.Year))
		}

		if content.Overview != "" {
			message.WriteString(fmt.Sprintf("\n\nتوضیحات: %s", content.Overview))
		}

		if content.Rating > 0 {
			message.WriteString(fmt.Sprintf("\n\nامتیاز: %.1f/10", content.Rating))
		}
	} else if content.Type == "Episode" {
		// Episode notification format
		message.WriteString("📺 قسمت جدید\n\n")

		if content.SeriesName != "" {
			message.WriteString(fmt.Sprintf("سریال: %s\n", content.SeriesName))
		} else {
			message.WriteString(fmt.Sprintf("سریال: %s\n", content.Title))
		}

		message.WriteString(fmt.Sprintf("فصل %d - قسمت %d", content.SeasonNumber, content.EpisodeNumber))

		if content.Title != "" && content.SeriesName != "" {
			message.WriteString(fmt.Sprintf("\nنام قسمت: %s", content.Title))
		}

		if content.Overview != "" {
			message.WriteString(fmt.Sprintf("\n\nتوضیحات: %s", content.Overview))
		}

		if content.Rating > 0 {
			message.WriteString(fmt.Sprintf("\n\nامتیاز: %.1f/10", content.Rating))
		}
	}

	return message.String()
}

// shouldShowMuteButton checks if mute button should be shown for this content
func shouldShowMuteButton(content *NotificationContent) bool {
	// Only show for episodes, not movies
	if content.Type != "Episode" {
		return false
	}

	// Don't show if series name is empty or "Unknown Series"
	if content.SeriesName == "" || content.SeriesName == "Unknown Series" {
		return false
	}

	return true
}

// createMuteButton creates inline keyboard with mute button
func createMuteButton(seriesName string) *botModels.InlineKeyboardMarkup {
	return &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{
					Text:         "دنبال نکردن",
					CallbackData: fmt.Sprintf("mute:%s", seriesName),
				},
			},
		},
	}
}

// BroadcastNotification sends a notification to all active subscribers
func (b *Bot) BroadcastNotification(ctx context.Context, content *NotificationContent) error {
	// Get all active subscribers
	subscribers, err := b.db.GetAllActiveSubscribers()
	if err != nil {
		return fmt.Errorf("failed to get subscribers: %w", err)
	}

	if len(subscribers) == 0 {
		slog.Info("No active subscribers to notify")
		return nil
	}

	// Check if NotifyOnlyTesters mode is enabled (for debugging/testing)
	filteredSubscribers := subscribers
	if b.config != nil && b.config.Testing.NotifyOnlyTesters {
		// Filter ALL notifications to only testers
		testerSubscribers := make([]int64, 0, len(subscribers))
		for _, chatID := range subscribers {
			if b.config.IsTester(chatID) {
				testerSubscribers = append(testerSubscribers, chatID)
			}
		}
		filteredSubscribers = testerSubscribers
		slog.Info("NotifyOnlyTesters mode enabled - filtered to testers only",
			"item_id", content.ItemID,
			"total_subscribers", len(subscribers),
			"tester_count", len(testerSubscribers))
	} else if len(content.ItemID) >= 5 && content.ItemID[:5] == "test-" {
		// Smart test detection: filter to only testers for test notifications
		if b.config != nil && b.config.Testing.EnableBetaFeatures {
			// Filter to only testers
			testSubscribers := make([]int64, 0, len(subscribers))
			for _, chatID := range subscribers {
				if b.config.IsTester(chatID) {
					testSubscribers = append(testSubscribers, chatID)
				}
			}
			filteredSubscribers = testSubscribers
			slog.Info("Test notification detected - filtered to testers only",
				"item_id", content.ItemID,
				"total_subscribers", len(subscribers),
				"tester_count", len(testSubscribers))
		}
	}

	// Filter out muted users for episode notifications
	mutedCount := 0

	if content.Type == "Episode" && content.SeriesName != "" {
		tempSubscribers := make([]int64, 0, len(filteredSubscribers))
		for _, chatID := range filteredSubscribers {
			isMuted, err := b.db.IsSeriesMuted(chatID, content.SeriesName)
			if err != nil {
				slog.Error("Failed to check if series is muted, including subscriber",
					"chat_id", chatID,
					"series_name", content.SeriesName,
					"error", err)
				// Include subscriber if check fails to avoid missing notifications
				tempSubscribers = append(tempSubscribers, chatID)
				continue
			}

			if !isMuted {
				tempSubscribers = append(tempSubscribers, chatID)
			} else {
				mutedCount++
			}
		}
		filteredSubscribers = tempSubscribers

		if mutedCount > 0 {
			slog.Info("Filtered muted users",
				"muted_count", mutedCount,
				"series_name", content.SeriesName)
		}
	}

	if len(filteredSubscribers) == 0 {
		slog.Info("No subscribers to notify after filtering",
			"total_subscribers", len(subscribers),
			"muted_count", mutedCount)
		return nil
	}

	slog.Info("Broadcasting notification",
		"content_type", content.Type,
		"title", content.Title,
		"subscriber_count", len(filteredSubscribers),
		"filtered_count", mutedCount)

	// Format notification message
	message := FormatNotification(content)

	// Create inline keyboard for episodes with valid series name
	var keyboard *botModels.InlineKeyboardMarkup
	if shouldShowMuteButton(content) {
		keyboard = createMuteButton(content.SeriesName)
	} else if !shouldShowMuteButton(content) && content.Type == "Episode" {
		slog.Debug("Skipping mute button",
			"reason", "invalid series name",
			"series_name", content.SeriesName)
	}

	// Fetch poster image
	var imageData []byte
	if content.ItemID != "" {
		imageData, err = b.jellyfinClient.GetPosterImage(ctx, content.ItemID)
		if err != nil {
			slog.Warn("Failed to fetch poster image for notification",
				"item_id", content.ItemID,
				"error", err)
			// Continue without image
		}
	}

	// Track broadcast statistics
	successCount := 0
	failureCount := 0
	blockedCount := 0

	// Broadcast to all filtered subscribers
	for _, chatID := range filteredSubscribers {
		// Handle Telegram rate limiting (max 30 messages/second)
		// Add small delay to avoid hitting rate limits
		time.Sleep(35 * time.Millisecond)

		var sendErr error
		if imageData != nil && len(imageData) > 0 {
			// Send with image
			if keyboard != nil {
				sendErr = b.SendPhotoBytesWithKeyboard(ctx, chatID, imageData, message, keyboard)
			} else {
				sendErr = b.SendPhotoBytes(ctx, chatID, imageData, message)
			}
		} else {
			// Send text only
			if keyboard != nil {
				sendErr = b.SendMessageWithKeyboard(ctx, chatID, message, keyboard)
			} else {
				sendErr = b.SendMessage(ctx, chatID, message)
			}
		}

		if sendErr != nil {
			// Check if bot was blocked by user
			errorStr := sendErr.Error()
			if strings.Contains(errorStr, "blocked") || strings.Contains(errorStr, "user is deactivated") ||
				strings.Contains(errorStr, "bot was blocked") || strings.Contains(errorStr, "chat not found") {
				slog.Warn("Bot blocked by user or chat not found, marking inactive",
					"chat_id", chatID,
					"error", sendErr)

				// Mark subscriber as inactive
				if err := b.db.RemoveSubscriber(chatID); err != nil {
					slog.Error("Failed to mark subscriber as inactive",
						"chat_id", chatID,
						"error", err)
				}
				blockedCount++
			} else {
				slog.Error("Failed to send notification",
					"chat_id", chatID,
					"error", sendErr)
				failureCount++
			}
		} else {
			successCount++
		}
	}

	slog.Info("Broadcast completed",
		"total_subscribers", len(subscribers),
		"muted_subscribers", mutedCount,
		"sent_to", len(filteredSubscribers),
		"success", successCount,
		"failures", failureCount,
		"blocked", blockedCount)

	return nil
}

// BroadcastNotificationWithRetry sends notification with retry logic for failures
func (b *Bot) BroadcastNotificationWithRetry(ctx context.Context, content *NotificationContent, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			slog.Info("Retrying broadcast after backoff",
				"attempt", attempt,
				"backoff", backoff)
			time.Sleep(backoff)
		}

		err := b.BroadcastNotification(ctx, content)
		if err == nil {
			return nil
		}

		lastErr = err
		slog.Warn("Broadcast attempt failed",
			"attempt", attempt,
			"error", err)
	}

	return fmt.Errorf("broadcast failed after %d attempts: %w", maxRetries, lastErr)
}
