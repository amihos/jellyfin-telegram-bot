package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

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
			Text:            "خطا در مسدود کردن سریال",
			ShowAlert:       false,
		})
		return
	}

	// Answer callback query
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		Text:            "✓ مسدود شد",
		ShowAlert:       false,
	})

	// Send confirmation message
	confirmationMsg := fmt.Sprintf("✓ شما دیگر اعلان‌های %s را دریافت نخواهید کرد", seriesName)
	if err := b.SendMessage(ctx, chatID, confirmationMsg); err != nil {
		slog.Error("Failed to send confirmation message",
			"chat_id", chatID,
			"error", err)
	}

	// Edit original message to disable button (change button text to "✓ مسدود شده")
	originalMessage := callbackQuery.Message.Message

	// Create new keyboard with disabled button
	disabledKeyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{
					Text:         "✓ مسدود شده",
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
				Text:            "سریال در لیست مسدودی‌ها یافت نشد",
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
			Text:            "خطا در رفع مسدودیت سریال",
			ShowAlert:       false,
		})
		return
	}

	// Answer callback query
	botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		Text:            "✓ رفع مسدودیت شد",
		ShowAlert:       false,
	})

	// Send confirmation message
	confirmationMsg := fmt.Sprintf("✓ %s از لیست مسدودی‌ها حذف شد", seriesName)
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
		messageText = "شما هیچ سریالی را مسدود نکرده‌اید"
	} else {
		messageText = "سریال‌های مسدود شده:\n\n"
		var buttons [][]botModels.InlineKeyboardButton

		for i, series := range updatedMutedSeries {
			messageText += fmt.Sprintf("%d. %s\n", i+1, series.SeriesName)
			// Create unmute button for each series
			buttons = append(buttons, []botModels.InlineKeyboardButton{
				{
					Text:         fmt.Sprintf("رفع مسدودیت: %s", series.SeriesName),
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
