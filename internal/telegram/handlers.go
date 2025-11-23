// SPDX-License-Identifier: MIT

package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"jellyfin-telegram-bot/internal/i18n"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

// handleStart handles the /start command
func (b *Bot) handleStart(ctx context.Context, _ *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	firstName := update.Message.From.FirstName
	telegramLangCode := update.Message.From.LanguageCode

	slog.Info("Processing /start command",
		"chat_id", chatID,
		"username", username,
		"first_name", firstName,
		"telegram_lang", telegramLangCode)

	// Add subscriber to database
	err := b.db.AddSubscriber(chatID, username, firstName)
	if err != nil {
		slog.Error("Failed to add subscriber",
			"chat_id", chatID,
			"error", err)

		localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)
		errorMsg := i18n.T(localizer, "error.generic")
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Detect and save user's language preference
	detectedLang := i18n.DetectLanguage(telegramLangCode, i18n.SupportedLanguages)
	if err := b.db.SetLanguage(chatID, detectedLang); err != nil {
		slog.Warn("Failed to set language preference",
			"chat_id", chatID,
			"detected_lang", detectedLang,
			"error", err)
	}

	// Get localizer for user
	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Send welcome message
	welcomeMessage := i18n.T(localizer, "welcome.message")

	// Create inline keyboard with 2x2 button grid
	keyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{Text: i18n.T(localizer, "button.recent"), CallbackData: "nav:recent"},
				{Text: i18n.T(localizer, "button.search"), CallbackData: "nav:search"},
			},
			{
				{Text: i18n.T(localizer, "button.mutedlist"), CallbackData: "nav:mutedlist"},
				{Text: i18n.T(localizer, "button.help"), CallbackData: "nav:help"},
			},
		},
	}

	// Send message with inline keyboard
	err = b.SendMessageWithKeyboard(ctx, chatID, welcomeMessage, keyboard)
	if err != nil {
		slog.Error("Failed to send welcome message with keyboard",
			"chat_id", chatID,
			"error", err)

		// Graceful fallback: send plain text message if keyboard fails
		err = b.SendMessage(ctx, chatID, welcomeMessage)
		if err != nil {
			slog.Error("Failed to send fallback welcome message",
				"chat_id", chatID,
				"error", err)
		}
		return
	}

	slog.Info("User subscribed successfully",
		"chat_id", chatID,
		"username", username,
		"language", detectedLang)
}

// handleRecent handles the /recent command
func (b *Bot) handleRecent(ctx context.Context, _ *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramLangCode := update.Message.From.LanguageCode

	slog.Info("Processing /recent command", "chat_id", chatID)

	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Fetch recent items from Jellyfin
	items, err := b.jellyfinClient.GetRecentItems(ctx, 15)
	if err != nil {
		slog.Error("Failed to fetch recent items",
			"chat_id", chatID,
			"error", err)

		errorMsg := i18n.T(localizer, "recent.error")
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Handle empty results
	if len(items) == 0 {
		b.SendMessage(ctx, chatID, i18n.T(localizer, "recent.no_results"))
		return
	}

	// Send each item with poster and formatted message
	for _, item := range items {
		b.sendContentItem(ctx, chatID, &item, localizer)
	}

	slog.Info("Sent recent items",
		"chat_id", chatID,
		"count", len(items))
}

// handleSearch handles the /search command
func (b *Bot) handleSearch(ctx context.Context, _ *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text
	telegramLangCode := update.Message.From.LanguageCode

	slog.Info("Processing /search command",
		"chat_id", chatID,
		"text", text)

	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Extract search query (remove "/search " prefix)
	query := strings.TrimSpace(strings.TrimPrefix(text, "/search"))

	// Check if query is empty
	if query == "" {
		helpMsg := i18n.T(localizer, "search.prompt")
		b.SendMessage(ctx, chatID, helpMsg)
		return
	}

	// Search content in Jellyfin
	items, err := b.jellyfinClient.SearchContent(ctx, query, 10)
	if err != nil {
		slog.Error("Failed to search content",
			"chat_id", chatID,
			"query", query,
			"error", err)

		errorMsg := i18n.T(localizer, "search.error")
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Handle empty results
	if len(items) == 0 {
		noResultsMsg := i18n.TWithData(localizer, "search.no_results", map[string]interface{}{
			"Query": query,
		})
		b.SendMessage(ctx, chatID, noResultsMsg)
		return
	}

	// Send each item with poster and formatted message
	for _, item := range items {
		b.sendContentItem(ctx, chatID, &item, localizer)
	}

	slog.Info("Sent search results",
		"chat_id", chatID,
		"query", query,
		"count", len(items))
}

// handleMutedList handles the /mutedlist command
func (b *Bot) handleMutedList(ctx context.Context, _ *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramLangCode := update.Message.From.LanguageCode

	slog.Info("Processing /mutedlist command", "chat_id", chatID)

	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Get all muted series for this user
	mutedSeries, err := b.db.GetMutedSeriesByUser(chatID)
	if err != nil {
		slog.Error("Failed to get muted series",
			"chat_id", chatID,
			"error", err)

		errorMsg := i18n.T(localizer, "mutedlist.error")
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Handle empty list case
	if len(mutedSeries) == 0 {
		emptyMsg := i18n.T(localizer, "mutedlist.empty")
		b.SendMessage(ctx, chatID, emptyMsg)
		return
	}

	// Format response message
	var messageText strings.Builder
	messageText.WriteString(i18n.T(localizer, "mutedlist.title"))
	messageText.WriteString("\n\n")

	// Create inline keyboard with unmute button for each series
	var buttons [][]botModels.InlineKeyboardButton

	for i, series := range mutedSeries {
		// Add series to message
		messageText.WriteString(fmt.Sprintf("%d. %s\n", i+1, series.SeriesName))

		// Create unmute button for this series
		buttonText := i18n.TWithData(localizer, "button.unmute", map[string]interface{}{
			"SeriesName": series.SeriesName,
		})
		buttons = append(buttons, []botModels.InlineKeyboardButton{
			{
				Text:         buttonText,
				CallbackData: fmt.Sprintf("unmute:%s", series.SeriesID),
			},
		})
	}

	keyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	// Send message with inline keyboard
	err = b.SendMessageWithKeyboard(ctx, chatID, messageText.String(), keyboard)
	if err != nil {
		slog.Error("Failed to send muted list",
			"chat_id", chatID,
			"error", err)
	}

	slog.Info("Sent muted list",
		"chat_id", chatID,
		"count", len(mutedSeries))
}

// handleLanguage handles the /language command
func (b *Bot) handleLanguage(ctx context.Context, _ *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramLangCode := update.Message.From.LanguageCode

	slog.Info("Processing /language command", "chat_id", chatID)

	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Create language selection keyboard
	keyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{
					Text:         i18n.T(localizer, "language.button.english"),
					CallbackData: "lang:en",
				},
			},
			{
				{
					Text:         i18n.T(localizer, "language.button.persian"),
					CallbackData: "lang:fa",
				},
			},
		},
	}

	// Send language selection prompt
	promptMsg := i18n.T(localizer, "language.select")
	err := b.SendMessageWithKeyboard(ctx, chatID, promptMsg, keyboard)
	if err != nil {
		slog.Error("Failed to send language selection",
			"chat_id", chatID,
			"error", err)
	}
}

// handleLanguageCallback handles language selection from inline keyboard
func (b *Bot) handleLanguageCallback(ctx context.Context, botInstance *bot.Bot, update *botModels.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callbackQuery := update.CallbackQuery
	if callbackQuery.Message.Message == nil {
		slog.Warn("Callback query message is nil")
		return
	}

	chatID := callbackQuery.Message.Message.Chat.ID
	callbackData := callbackQuery.Data

	slog.Info("Processing language callback",
		"chat_id", chatID,
		"callback_data", callbackData)

	// Parse language from callback data (format: "lang:{code}")
	parts := strings.SplitN(callbackData, ":", 2)
	if len(parts) != 2 {
		slog.Error("Invalid language callback data format",
			"callback_data", callbackData)

		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            "Error processing request",
			ShowAlert:       false,
		})
		return
	}

	selectedLang := parts[1]

	// Validate selected language
	if !i18n.IsSupportedLanguage(selectedLang) {
		slog.Error("Unsupported language selected",
			"language", selectedLang)

		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            "Unsupported language",
			ShowAlert:       false,
		})
		return
	}

	// Save language preference to database
	if err := b.db.SetLanguage(chatID, selectedLang); err != nil {
		slog.Error("Failed to set language preference",
			"chat_id", chatID,
			"language", selectedLang,
			"error", err)

		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            "Error saving language preference",
			ShowAlert:       false,
		})
		return
	}

	// Get localizer in NEW language
	localizer := i18n.GetLocalizer(b.i18nBundle, selectedLang)

	// Answer callback query
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		ShowAlert:       false,
	})

	// Send confirmation in NEW language
	confirmationMsg := i18n.T(localizer, "language.selected")
	if err := b.SendMessage(ctx, chatID, confirmationMsg); err != nil {
		slog.Error("Failed to send language confirmation",
			"chat_id", chatID,
			"error", err)
	}

	slog.Info("Language preference updated",
		"chat_id", chatID,
		"language", selectedLang)
}

// sendContentItem sends a single content item with poster and formatted message
func (b *Bot) sendContentItem(ctx context.Context, chatID int64, item *ContentItem, localizer *goi18n.Localizer) {
	// Format message using i18n
	message := FormatContentMessage(item, localizer)

	// Try to fetch and send poster image
	imageData, err := b.jellyfinClient.GetPosterImage(ctx, item.ItemID)
	if err != nil {
		slog.Warn("Failed to fetch poster image, sending text only",
			"item_id", item.ItemID,
			"error", err)

		// Send text message only if image fetch fails
		if err := b.SendMessage(ctx, chatID, message); err != nil {
			slog.Error("Failed to send content message",
				"chat_id", chatID,
				"item_id", item.ItemID,
				"error", err)
		}
		return
	}

	// Send photo with caption
	if err := b.SendPhotoBytes(ctx, chatID, imageData, message); err != nil {
		slog.Error("Failed to send content photo",
			"chat_id", chatID,
			"item_id", item.ItemID,
			"error", err)

		// Fallback to text message if photo send fails
		if err := b.SendMessage(ctx, chatID, message); err != nil {
			slog.Error("Failed to send fallback content message",
				"chat_id", chatID,
				"item_id", item.ItemID,
				"error", err)
		}
	}
}

// FormatContentMessage formats a content item for display using i18n
func FormatContentMessage(item *ContentItem, localizer *goi18n.Localizer) string {
	var message strings.Builder

	// Type indicator and title
	if item.Type == "Movie" {
		message.WriteString(i18n.T(localizer, "content.field.movie"))
		message.WriteString("\n\n")
		message.WriteString(i18n.TWithData(localizer, "content.field.name", map[string]interface{}{
			"Name": item.Name,
		}))
	} else if item.Type == "Episode" {
		message.WriteString(i18n.T(localizer, "content.field.episode"))
		message.WriteString("\n\n")
		if item.SeriesName != "" {
			message.WriteString(i18n.TWithData(localizer, "content.field.series", map[string]interface{}{
				"SeriesName": item.SeriesName,
			}))
			message.WriteString("\n")
		}
		message.WriteString(i18n.TWithData(localizer, "content.field.episode_number", map[string]interface{}{
			"SeasonNumber":  item.SeasonNumber,
			"EpisodeNumber": item.EpisodeNumber,
		}))
		if item.Name != "" {
			message.WriteString("\n")
			message.WriteString(i18n.TWithData(localizer, "content.field.episode_name", map[string]interface{}{
				"Name": item.Name,
			}))
		}
	}

	// Production year
	if item.ProductionYear > 0 {
		message.WriteString("\n")
		message.WriteString(i18n.TWithData(localizer, "content.field.year", map[string]interface{}{
			"Year": item.ProductionYear,
		}))
	}

	// Description
	if item.Overview != "" {
		message.WriteString("\n\n")
		message.WriteString(i18n.TWithData(localizer, "content.field.description", map[string]interface{}{
			"Description": item.Overview,
		}))
	}

	// Rating
	if item.CommunityRating > 0 {
		message.WriteString("\n\n")
		message.WriteString(i18n.TWithData(localizer, "content.field.rating", map[string]interface{}{
			"Rating": fmt.Sprintf("%.1f", item.CommunityRating),
		}))
	} else if item.OfficialRating != "" {
		message.WriteString("\n\n")
		message.WriteString(i18n.TWithData(localizer, "content.field.official_rating", map[string]interface{}{
			"Rating": item.OfficialRating,
		}))
	}

	return message.String()
}
