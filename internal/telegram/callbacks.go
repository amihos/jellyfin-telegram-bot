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
	"gorm.io/gorm"
)

// handleNavigationCallback handles navigation button callbacks from the welcome menu
func (b *Bot) handleNavigationCallback(ctx context.Context, botInstance *bot.Bot, update *botModels.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callbackQuery := update.CallbackQuery

	// Check if message exists
	if callbackQuery.Message.Message == nil {
		slog.Warn("Callback query message is nil")
		return
	}

	chatID := callbackQuery.Message.Message.Chat.ID
	callbackData := callbackQuery.Data

	slog.Info("Processing navigation callback",
		"chat_id", chatID,
		"callback_data", callbackData)

	// Get localizer for user
	telegramLangCode := ""
	if callbackQuery.From.ID != 0 {
		telegramLangCode = callbackQuery.From.LanguageCode
	}
	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Parse action from callback data (format: "nav:{action}")
	parts := strings.SplitN(callbackData, ":", 2)
	if len(parts) != 2 {
		slog.Error("Invalid navigation callback data format",
			"callback_data", callbackData)

		// Always answer callback query to prevent stuck loading state
		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            i18n.T(localizer, "error.request_processing"),
			ShowAlert:       false,
		})
		return
	}

	action := parts[1]

	// Route to appropriate logic based on action
	switch action {
	case "recent":
		b.handleNavigationRecent(ctx, botInstance, chatID, callbackQuery, localizer)
	case "search":
		b.handleNavigationSearch(ctx, botInstance, chatID, callbackQuery, localizer)
	case "mutedlist":
		b.handleNavigationMutedList(ctx, botInstance, chatID, callbackQuery, localizer)
	case "help":
		b.handleNavigationHelp(ctx, botInstance, chatID, callbackQuery, localizer)
	default:
		slog.Warn("Unknown navigation action",
			"action", action,
			"chat_id", chatID)

		// Answer callback query with error
		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            i18n.T(localizer, "error.invalid_callback"),
			ShowAlert:       false,
		})
	}
}

// handleNavigationRecent handles nav:recent callback
func (b *Bot) handleNavigationRecent(ctx context.Context, botInstance *bot.Bot, chatID int64, callbackQuery *botModels.CallbackQuery, localizer *goi18n.Localizer) {
	slog.Info("Processing nav:recent", "chat_id", chatID)

	// Answer callback query immediately to remove loading state
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		ShowAlert:       false,
	})

	// Fetch recent items from Jellyfin (reuse handleRecent logic)
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

	slog.Info("Sent recent items via navigation",
		"chat_id", chatID,
		"count", len(items))
}

// handleNavigationSearch handles nav:search callback
func (b *Bot) handleNavigationSearch(ctx context.Context, botInstance *bot.Bot, chatID int64, callbackQuery *botModels.CallbackQuery, localizer *goi18n.Localizer) {
	slog.Info("Processing nav:search", "chat_id", chatID)

	// Answer callback query immediately to remove loading state
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		ShowAlert:       false,
	})

	// Send search instructions
	searchInstructions := i18n.T(localizer, "search.prompt")
	if err := b.SendMessage(ctx, chatID, searchInstructions); err != nil {
		slog.Error("Failed to send search instructions",
			"chat_id", chatID,
			"error", err)
	}

	slog.Info("Sent search instructions via navigation", "chat_id", chatID)
}

// handleNavigationMutedList handles nav:mutedlist callback
func (b *Bot) handleNavigationMutedList(ctx context.Context, botInstance *bot.Bot, chatID int64, callbackQuery *botModels.CallbackQuery, localizer *goi18n.Localizer) {
	slog.Info("Processing nav:mutedlist", "chat_id", chatID)

	// Answer callback query immediately to remove loading state
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		ShowAlert:       false,
	})

	// Get all muted series for this user (reuse handleMutedList logic)
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
		slog.Error("Failed to send muted list via navigation",
			"chat_id", chatID,
			"error", err)
	}

	slog.Info("Sent muted list via navigation",
		"chat_id", chatID,
		"count", len(mutedSeries))
}

// handleNavigationHelp handles nav:help callback
func (b *Bot) handleNavigationHelp(ctx context.Context, botInstance *bot.Bot, chatID int64, callbackQuery *botModels.CallbackQuery, localizer *goi18n.Localizer) {
	slog.Info("Processing nav:help", "chat_id", chatID)

	// Answer callback query immediately to remove loading state
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		ShowAlert:       false,
	})

	// Send help message
	helpMessage := i18n.T(localizer, "help.message")

	if err := b.SendMessage(ctx, chatID, helpMessage); err != nil {
		slog.Error("Failed to send help message via navigation",
			"chat_id", chatID,
			"error", err)
	}

	slog.Info("Sent help message via navigation", "chat_id", chatID)
}

// handleMuteCallback handles the mute button callback
func (b *Bot) handleMuteCallback(ctx context.Context, botInstance *bot.Bot, update *botModels.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callbackQuery := update.CallbackQuery

	// Check if message exists
	if callbackQuery.Message.Message == nil {
		slog.Warn("Callback query message is nil")
		return
	}

	chatID := callbackQuery.Message.Message.Chat.ID
	callbackData := callbackQuery.Data

	slog.Info("Processing mute callback",
		"chat_id", chatID,
		"callback_data", callbackData)

	// Get localizer for user
	telegramLangCode := ""
	if callbackQuery.From.ID != 0 {
		telegramLangCode = callbackQuery.From.LanguageCode
	}
	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Parse series name from callback data (format: "mute:{SeriesName}")
	parts := strings.SplitN(callbackData, ":", 2)
	if len(parts) != 2 {
		slog.Error("Invalid callback data format",
			"callback_data", callbackData)
		return
	}

	seriesName := parts[1]

	// Add muted series to database
	err := b.db.AddMutedSeries(chatID, seriesName, seriesName)
	if err != nil {
		slog.Error("Failed to add muted series",
			"chat_id", chatID,
			"series_name", seriesName,
			"error", err)

		// Answer callback query with error
		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            i18n.T(localizer, "mute.error"),
			ShowAlert:       false,
		})
		return
	}

	// Answer callback query
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		Text:            i18n.T(localizer, "mute.callback_success"),
		ShowAlert:       false,
	})

	// Send confirmation message with undo button
	confirmationMsg := i18n.TWithData(localizer, "mute.success", map[string]interface{}{
		"SeriesName": seriesName,
	})

	// Create inline keyboard with undo button
	undoKeyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{
					Text:         i18n.T(localizer, "button.undo_mute"),
					CallbackData: fmt.Sprintf("undo_mute:%s", seriesName),
				},
			},
		},
	}

	if err := b.SendMessageWithKeyboard(ctx, chatID, confirmationMsg, undoKeyboard); err != nil {
		slog.Error("Failed to send confirmation message with undo button",
			"chat_id", chatID,
			"error", err)
	}

	// Edit original message to disable button (change button text to muted state)
	originalMessage := callbackQuery.Message.Message

	// Create new keyboard with disabled button
	disabledKeyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{
					Text:         i18n.T(localizer, "button.muted"),
					CallbackData: "muted", // Inactive callback data
				},
			},
		},
	}

	// Edit message to update button
	_, err = botInstance.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      chatID,
		MessageID:   originalMessage.ID,
		ReplyMarkup: disabledKeyboard,
	})
	if err != nil {
		slog.Warn("Failed to edit message markup",
			"chat_id", chatID,
			"message_id", originalMessage.ID,
			"error", err)
	}

	slog.Info("Successfully muted series",
		"chat_id", chatID,
		"series_name", seriesName)
}

// handleUndoMuteCallback handles the undo mute button callback
func (b *Bot) handleUndoMuteCallback(ctx context.Context, botInstance *bot.Bot, update *botModels.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callbackQuery := update.CallbackQuery

	// Check if message exists
	if callbackQuery.Message.Message == nil {
		slog.Warn("Callback query message is nil")
		return
	}

	chatID := callbackQuery.Message.Message.Chat.ID
	callbackData := callbackQuery.Data

	slog.Info("Processing undo mute callback",
		"chat_id", chatID,
		"callback_data", callbackData)

	// Get localizer for user
	telegramLangCode := ""
	if callbackQuery.From.ID != 0 {
		telegramLangCode = callbackQuery.From.LanguageCode
	}
	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Parse series name from callback data (format: "undo_mute:{SeriesName}")
	parts := strings.SplitN(callbackData, ":", 2)
	if len(parts) != 2 {
		slog.Error("Invalid callback data format",
			"callback_data", callbackData)

		// Always answer callback query
		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            i18n.T(localizer, "error.request_processing"),
			ShowAlert:       false,
		})
		return
	}

	seriesName := parts[1]

	// Remove muted series from database (reuse unmute logic from handleUnmuteCallback)
	err := b.db.RemoveMutedSeries(chatID, seriesName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Debug("Series not found in muted list",
				"chat_id", chatID,
				"series_name", seriesName)

			// Answer callback query
			botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: callbackQuery.ID,
				Text:            i18n.T(localizer, "unmute.not_found"),
				ShowAlert:       false,
			})
			return
		}

		slog.Error("Failed to remove muted series",
			"chat_id", chatID,
			"series_name", seriesName,
			"error", err)

		// Answer callback query with error
		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            i18n.T(localizer, "unmute.error"),
			ShowAlert:       false,
		})
		return
	}

	// Answer callback query
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		Text:            i18n.T(localizer, "unmute.callback_success"),
		ShowAlert:       false,
	})

	// Send confirmation message
	confirmationMsg := i18n.TWithData(localizer, "unmute.success", map[string]interface{}{
		"SeriesName": seriesName,
	})
	if err := b.SendMessage(ctx, chatID, confirmationMsg); err != nil {
		slog.Error("Failed to send confirmation message",
			"chat_id", chatID,
			"error", err)
	}

	// Update button state to show success
	originalMessage := callbackQuery.Message.Message

	// Create new keyboard with success button
	successKeyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{
					Text:         i18n.T(localizer, "unmute.callback_success"),
					CallbackData: "unmuted", // Inactive callback data
				},
			},
		},
	}

	// Edit message to update button
	_, err = botInstance.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      chatID,
		MessageID:   originalMessage.ID,
		ReplyMarkup: successKeyboard,
	})
	if err != nil {
		slog.Warn("Failed to edit message markup",
			"chat_id", chatID,
			"message_id", originalMessage.ID,
			"error", err)
	}

	slog.Info("Successfully unmuted series via undo",
		"chat_id", chatID,
		"series_name", seriesName)
}

// handleUnmuteCallback handles the unmute button callback
func (b *Bot) handleUnmuteCallback(ctx context.Context, botInstance *bot.Bot, update *botModels.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callbackQuery := update.CallbackQuery

	// Check if message exists
	if callbackQuery.Message.Message == nil {
		slog.Warn("Callback query message is nil")
		return
	}

	chatID := callbackQuery.Message.Message.Chat.ID
	callbackData := callbackQuery.Data

	slog.Info("Processing unmute callback",
		"chat_id", chatID,
		"callback_data", callbackData)

	// Get localizer for user
	telegramLangCode := ""
	if callbackQuery.From.ID != 0 {
		telegramLangCode = callbackQuery.From.LanguageCode
	}
	localizer := b.getLocalizerForUser(ctx, chatID, telegramLangCode)

	// Parse series ID from callback data (format: "unmute:{SeriesID}")
	parts := strings.SplitN(callbackData, ":", 2)
	if len(parts) != 2 {
		slog.Error("Invalid callback data format",
			"callback_data", callbackData)
		return
	}

	seriesID := parts[1]

	// Get series name before removal (for confirmation message)
	mutedSeries, err := b.db.GetMutedSeriesByUser(chatID)
	seriesName := seriesID // Default to ID if we can't find the name
	if err == nil {
		for _, ms := range mutedSeries {
			if ms.SeriesID == seriesID {
				seriesName = ms.SeriesName
				break
			}
		}
	}

	// Remove muted series from database
	err = b.db.RemoveMutedSeries(chatID, seriesID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Debug("Series not found in muted list",
				"chat_id", chatID,
				"series_id", seriesID)

			// Answer callback query
			botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: callbackQuery.ID,
				Text:            i18n.T(localizer, "unmute.not_found"),
				ShowAlert:       false,
			})
			return
		}

		slog.Error("Failed to remove muted series",
			"chat_id", chatID,
			"series_id", seriesID,
			"error", err)

		// Answer callback query with error
		botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQuery.ID,
			Text:            i18n.T(localizer, "unmute.error"),
			ShowAlert:       false,
		})
		return
	}

	// Answer callback query
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		Text:            i18n.T(localizer, "unmute.callback_success"),
		ShowAlert:       false,
	})

	// Send confirmation message
	confirmationMsg := i18n.TWithData(localizer, "unmute.success", map[string]interface{}{
		"SeriesName": seriesName,
	})
	if err := b.SendMessage(ctx, chatID, confirmationMsg); err != nil {
		slog.Error("Failed to send confirmation message",
			"chat_id", chatID,
			"error", err)
	}

	// Refresh /mutedlist message by regenerating it
	// Get updated muted series list
	updatedMutedSeries, err := b.db.GetMutedSeriesByUser(chatID)
	if err != nil {
		slog.Error("Failed to get updated muted series list",
			"chat_id", chatID,
			"error", err)
		return
	}

	// Generate new message
	var messageText string
	var keyboard *botModels.InlineKeyboardMarkup

	if len(updatedMutedSeries) == 0 {
		messageText = i18n.T(localizer, "mutedlist.empty")
	} else {
		messageText = i18n.T(localizer, "mutedlist.title") + "\n\n"
		var buttons [][]botModels.InlineKeyboardButton

		for i, series := range updatedMutedSeries {
			messageText += fmt.Sprintf("%d. %s\n", i+1, series.SeriesName)
			// Create unmute button for each series
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

		keyboard = &botModels.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		}
	}

	// Edit original message
	_, err = botInstance.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   callbackQuery.Message.Message.ID,
		Text:        messageText,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		slog.Warn("Failed to refresh mutedlist message",
			"chat_id", chatID,
			"error", err)
	}

	slog.Info("Successfully unmuted series",
		"chat_id", chatID,
		"series_id", seriesID,
		"series_name", seriesName)
}
