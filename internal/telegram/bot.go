// SPDX-License-Identifier: MIT

package telegram

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"jellyfin-telegram-bot/internal/config"
	"jellyfin-telegram-bot/internal/i18n"
	"jellyfin-telegram-bot/pkg/models"

	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// Bot represents the Telegram bot instance
type Bot struct {
	bot            *bot.Bot
	db             SubscriberDB
	jellyfinClient JellyfinClient
	config         *config.Config
	i18nBundle     *goi18n.Bundle
}

// SubscriberDB defines the interface for subscriber operations
type SubscriberDB interface {
	AddSubscriber(chatID int64, username, firstName string) error
	RemoveSubscriber(chatID int64) error
	GetAllActiveSubscribers() ([]int64, error)
	IsSubscribed(chatID int64) (bool, error)
	// Language operations
	SetLanguage(chatID int64, languageCode string) error
	GetLanguage(chatID int64) (string, error)
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

	// Initialize i18n bundle
	bundle, err := i18n.InitBundle()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize i18n: %w", err)
	}

	botInstance := &Bot{
		db:             db,
		jellyfinClient: jellyfinClient,
		config:         cfg,
		i18nBundle:     bundle,
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(botInstance.defaultHandler),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, botInstance.handleStart),
		bot.WithMessageTextHandler("/recent", bot.MatchTypeExact, botInstance.handleRecent),
		bot.WithMessageTextHandler("/search", bot.MatchTypePrefix, botInstance.handleSearch),
		bot.WithMessageTextHandler("/mutedlist", bot.MatchTypeExact, botInstance.handleMutedList),
		bot.WithMessageTextHandler("/language", bot.MatchTypeExact, botInstance.handleLanguage),
		bot.WithCallbackQueryDataHandler("nav:", bot.MatchTypePrefix, botInstance.handleNavigationCallback),
		bot.WithCallbackQueryDataHandler("mute:", bot.MatchTypePrefix, botInstance.handleMuteCallback),
		bot.WithCallbackQueryDataHandler("undo_mute:", bot.MatchTypePrefix, botInstance.handleUndoMuteCallback),
		bot.WithCallbackQueryDataHandler("unmute:", bot.MatchTypePrefix, botInstance.handleUnmuteCallback),
		bot.WithCallbackQueryDataHandler("lang:", bot.MatchTypePrefix, botInstance.handleLanguageCallback),
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
	// Register commands for each supported language
	for _, langCode := range i18n.SupportedLanguages {
		localizer := i18n.GetLocalizer(b.i18nBundle, langCode)

		commands := []botModels.BotCommand{
			{
				Command:     "start",
				Description: i18n.T(localizer, "command.start.description"),
			},
			{
				Command:     "recent",
				Description: i18n.T(localizer, "command.recent.description"),
			},
			{
				Command:     "search",
				Description: i18n.T(localizer, "command.search.description"),
			},
			{
				Command:     "mutedlist",
				Description: i18n.T(localizer, "command.mutedlist.description"),
			},
			{
				Command:     "language",
				Description: i18n.T(localizer, "command.language.description"),
			},
		}

		_, err := b.bot.SetMyCommands(ctx, &bot.SetMyCommandsParams{
			Commands:     commands,
			LanguageCode: langCode,
		})

		if err != nil {
			slog.Warn("Failed to set bot commands for language",
				"language", langCode,
				"error", err)
			// Continue with other languages even if one fails
		} else {
			slog.Info("Bot commands registered for language",
				"language", langCode,
				"count", len(commands))
		}
	}

	return nil
}

// Start starts the bot with polling mode
func (b *Bot) Start(ctx context.Context) {
	slog.Info("Starting Telegram bot...")
	b.bot.Start(ctx)
}

// defaultHandler handles unknown commands and messages
func (b *Bot) defaultHandler(ctx context.Context, bot *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	// Get user's language preference
	localizer := b.getLocalizerForUser(ctx, chatID, update.Message.From.LanguageCode)

	// Send help message
	helpMessage := i18n.T(localizer, "help.invalid_command")

	err := b.SendMessage(ctx, chatID, helpMessage)
	if err != nil {
		slog.Error("Failed to send help message",
			"chat_id", chatID,
			"error", err)
	}
}

// getLocalizerForUser gets the localizer for a user with fallback chain
// Fallback chain: saved preference → Telegram language → English default
func (b *Bot) getLocalizerForUser(ctx context.Context, chatID int64, telegramLangCode string) *goi18n.Localizer {
	// Try to get saved language preference
	savedLang, err := b.db.GetLanguage(chatID)
	if err == nil && savedLang != "" {
		return i18n.GetLocalizer(b.i18nBundle, savedLang)
	}

	// Fallback to Telegram language code
	if telegramLangCode != "" {
		detectedLang := i18n.DetectLanguage(telegramLangCode, i18n.SupportedLanguages)
		return i18n.GetLocalizer(b.i18nBundle, detectedLang)
	}

	// Final fallback to English
	return i18n.GetLocalizer(b.i18nBundle, i18n.DefaultLanguage)
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
