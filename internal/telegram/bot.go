package telegram

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Bot represents the Telegram bot instance
type Bot struct {
	bot            *bot.Bot
	db             SubscriberDB
	jellyfinClient JellyfinClient
}

// SubscriberDB defines the interface for subscriber operations
type SubscriberDB interface {
	AddSubscriber(chatID int64, username, firstName string) error
	RemoveSubscriber(chatID int64) error
	GetAllActiveSubscribers() ([]int64, error)
	IsSubscribed(chatID int64) (bool, error)
}

// JellyfinClient defines the interface for Jellyfin API operations
type JellyfinClient interface {
	GetRecentItems(ctx context.Context, limit int) ([]ContentItem, error)
	SearchContent(ctx context.Context, query string, limit int) ([]ContentItem, error)
	GetPosterImage(ctx context.Context, itemID string) ([]byte, error)
}

// ContentItem represents content from Jellyfin (local interface to avoid circular imports)
type ContentItem struct {
	ItemID          string
	Name            string
	Type            string
	Overview        string
	CommunityRating float64
	OfficialRating  string
	ProductionYear  int
	SeriesName      string
	SeasonNumber    int
	EpisodeNumber   int
}

// NewBot creates a new Telegram bot instance
func NewBot(token string, db SubscriberDB, jellyfinClient JellyfinClient) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	botInstance := &Bot{
		bot:            b,
		db:             db,
		jellyfinClient: jellyfinClient,
	}

	// Register command handlers
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, botInstance.handleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/recent", bot.MatchTypeExact, botInstance.handleRecent)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypePrefix, botInstance.handleSearch)

	slog.Info("Telegram bot initialized successfully")

	return botInstance, nil
}

// Start starts the bot with polling mode
func (b *Bot) Start(ctx context.Context) {
	slog.Info("Starting Telegram bot...")
	b.bot.Start(ctx)
}

// defaultHandler handles unknown commands and messages
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	// Send help message in Persian for unknown commands
	helpMessage := `دستور نامعتبر است.

دستورات موجود:
/start - عضویت در ربات
/recent - مشاهده محتوای اخیر
/search - جستجوی محتوا (مثال: /search interstellar)`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpMessage,
	})

	if err != nil {
		slog.Error("Failed to send help message",
			"chat_id", update.Message.Chat.ID,
			"error", err)
	}
}

// SendMessage sends a text message to a chat
func (b *Bot) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := b.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	return err
}

// SendPhotoBytes sends a photo with caption to a chat using byte data
func (b *Bot) SendPhotoBytes(ctx context.Context, chatID int64, imageData []byte, caption string) error {
	_, err := b.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:  chatID,
		Photo:   &models.InputFileUpload{Data: bytes.NewReader(imageData), Filename: "poster.jpg"},
		Caption: caption,
	})
	return err
}

// GetBot returns the underlying bot instance (for testing)
func (b *Bot) GetBot() *bot.Bot {
	return b.bot
}
