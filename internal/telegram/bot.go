package telegram

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"jellyfin-telegram-bot/internal/config"
	"jellyfin-telegram-bot/pkg/models"

	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

// Bot represents the Telegram bot instance
type Bot struct {
	bot            *bot.Bot
	db             SubscriberDB
	jellyfinClient JellyfinClient
	config         *config.Config
}

// SubscriberDB defines the interface for subscriber operations
type SubscriberDB interface {
	AddSubscriber(chatID int64, username, firstName string) error
	RemoveSubscriber(chatID int64) error
	GetAllActiveSubscribers() ([]int64, error)
	IsSubscribed(chatID int64) (bool, error)
	// Mute operations
	AddMutedSeries(chatID int64, seriesID string, seriesName string) error
	RemoveMutedSeries(chatID int64, seriesID string) error
	GetMutedSeriesByUser(chatID int64) ([]models.MutedSeries, error)
	IsSeriesMuted(chatID int64, seriesID string) (bool, error)
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
func NewBot(token string, db SubscriberDB, jellyfinClient JellyfinClient, cfg *config.Config) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	botInstance := &Bot{
		db:             db,
		jellyfinClient: jellyfinClient,
		config:         cfg,
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, botInstance.handleStart),
		bot.WithMessageTextHandler("/recent", bot.MatchTypeExact, botInstance.handleRecent),
		bot.WithMessageTextHandler("/search", bot.MatchTypePrefix, botInstance.handleSearch),
		bot.WithMessageTextHandler("/mutedlist", bot.MatchTypeExact, botInstance.handleMutedList),
		bot.WithCallbackQueryDataHandler("nav:", bot.MatchTypePrefix, botInstance.handleNavigationCallback),
		bot.WithCallbackQueryDataHandler("mute:", bot.MatchTypePrefix, botInstance.handleMuteCallback),
		bot.WithCallbackQueryDataHandler("undo_mute:", bot.MatchTypePrefix, botInstance.handleUndoMuteCallback),
		bot.WithCallbackQueryDataHandler("unmute:", bot.MatchTypePrefix, botInstance.handleUnmuteCallback),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	botInstance.bot = b

	// Register bot commands for Telegram Menu Button API
	err = botInstance.registerBotCommands(context.Background())
	if err != nil {
		slog.Warn("Failed to register bot commands (non-fatal)",
			"error", err)
		// Non-fatal error - continue with bot initialization
	}

	slog.Info("Telegram bot initialized successfully")

	return botInstance, nil
}

// registerBotCommands registers bot commands with Telegram for Menu Button integration
func (b *Bot) registerBotCommands(ctx context.Context) error {
	commands := []botModels.BotCommand{
		{
			Command:     "start",
			Description: "عضویت در ربات",
		},
		{
			Command:     "recent",
			Description: "مشاهده محتوای اخیر",
		},
		{
			Command:     "search",
			Description: "جستجوی محتوا",
		},
		{
			Command:     "mutedlist",
			Description: "مشاهده سریال‌های مسدود شده",
		},
	}

	_, err := b.bot.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})

	if err != nil {
		return fmt.Errorf("failed to set bot commands: %w", err)
	}

	slog.Info("Bot commands registered successfully",
		"count", len(commands))

	return nil
}

// Start starts the bot with polling mode
func (b *Bot) Start(ctx context.Context) {
	slog.Info("Starting Telegram bot...")
	b.bot.Start(ctx)
}

// defaultHandler handles unknown commands and messages
func defaultHandler(ctx context.Context, b *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	// Send help message in Persian for unknown commands
	helpMessage := `دستور نامعتبر است.

دستورات موجود:
/start - عضویت در ربات
/recent - مشاهده محتوای اخیر
/search - جستجوی محتوا (مثال: /search interstellar)
/mutedlist - مشاهده سریال‌های مسدود شده`

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

// SendMessageWithKeyboard sends a text message with an inline keyboard to a chat
func (b *Bot) SendMessageWithKeyboard(ctx context.Context, chatID int64, text string, keyboard *botModels.InlineKeyboardMarkup) error {
	_, err := b.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	return err
}

// SendPhotoBytes sends a photo with caption to a chat using byte data
func (b *Bot) SendPhotoBytes(ctx context.Context, chatID int64, imageData []byte, caption string) error {
	_, err := b.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:  chatID,
		Photo:   &botModels.InputFileUpload{Data: bytes.NewReader(imageData), Filename: "poster.jpg"},
		Caption: caption,
	})
	return err
}

// SendPhotoBytesWithKeyboard sends a photo with caption and inline keyboard to a chat using byte data
func (b *Bot) SendPhotoBytesWithKeyboard(ctx context.Context, chatID int64, imageData []byte, caption string, keyboard *botModels.InlineKeyboardMarkup) error {
	_, err := b.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      chatID,
		Photo:       &botModels.InputFileUpload{Data: bytes.NewReader(imageData), Filename: "poster.jpg"},
		Caption:     caption,
		ReplyMarkup: keyboard,
	})
	return err
}

// GetBot returns the underlying bot instance (for testing)
func (b *Bot) GetBot() *bot.Bot {
	return b.bot
}
