package telegram

import (
	"context"

	"jellyfin-telegram-bot/internal/handlers"
)

// BroadcasterAdapter adapts the Bot to work as a NotificationBroadcaster
// This allows the webhook handler to call broadcast methods on the bot
type BroadcasterAdapter struct {
	bot *Bot
}

// NewBroadcasterAdapter creates a new broadcaster adapter
func NewBroadcasterAdapter(bot *Bot) *BroadcasterAdapter {
	return &BroadcasterAdapter{
		bot: bot,
	}
}

// BroadcastNotification implements the NotificationBroadcaster interface
// It accepts webhook notification content and converts it for the bot
func (a *BroadcasterAdapter) BroadcastNotification(ctx context.Context, content *handlers.NotificationContent) error {
	// Convert the webhook content to bot notification content
	notifContent := &NotificationContent{
		ItemID:        content.ItemID,
		Type:          content.Type,
		Title:         content.Title,
		Overview:      content.Overview,
		Year:          content.Year,
		Rating:        content.Rating,
		SeriesName:    content.SeriesName,
		SeasonNumber:  content.SeasonNumber,
		EpisodeNumber: content.EpisodeNumber,
	}

	// Call the bot's broadcast method
	return a.bot.BroadcastNotification(ctx, notifContent)
}
