package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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

	slog.Info("Processing /start command",
		"chat_id", chatID,
		"username", username,
		"first_name", firstName)

	// Add subscriber to database
	err := b.db.AddSubscriber(chatID, username, firstName)
	if err != nil {
		slog.Error("Failed to add subscriber",
			"chat_id", chatID,
			"error", err)

		errorMsg := "متأسفانه خطایی رخ داده. لطفاً دوباره تلاش کنید."
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Send welcome message in Persian
	welcomeMessage := `سلام! به ربات اطلاع‌رسانی جلیفین خوش آمدید.

شما از این پس اطلاعیه‌های محتوای جدید را دریافت خواهید کرد.

دستورات موجود:
/start - عضویت در ربات
/recent - مشاهده محتوای اخیر
/search - جستجوی محتوا
/mutedlist - مشاهده سریال‌های مسدود شده`

	// Create inline keyboard with 2x2 button grid
	keyboard := &botModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]botModels.InlineKeyboardButton{
			{
				{Text: "تازه‌ها", CallbackData: "nav:recent"},
				{Text: "جستجو", CallbackData: "nav:search"},
			},
			{
				{Text: "سریال‌های مسدود شده", CallbackData: "nav:mutedlist"},
				{Text: "راهنما", CallbackData: "nav:help"},
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
		"username", username)
}

// handleRecent handles the /recent command
func (b *Bot) handleRecent(ctx context.Context, _ *bot.Bot, update *botModels.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	slog.Info("Processing /recent command", "chat_id", chatID)

	// Fetch recent items from Jellyfin
	items, err := b.jellyfinClient.GetRecentItems(ctx, 15)
	if err != nil {
		slog.Error("Failed to fetch recent items",
			"chat_id", chatID,
			"error", err)

		errorMsg := "خطا در دریافت محتوای اخیر. لطفاً بعداً تلاش کنید."
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Handle empty results
	if len(items) == 0 {
		b.SendMessage(ctx, chatID, "محتوای اخیری یافت نشد")
		return
	}

	// Send each item with poster and formatted message
	for _, item := range items {
		b.sendContentItem(ctx, chatID, &item)
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

	slog.Info("Processing /search command",
		"chat_id", chatID,
		"text", text)

	// Extract search query (remove "/search " prefix)
	query := strings.TrimSpace(strings.TrimPrefix(text, "/search"))

	// Check if query is empty
	if query == "" {
		helpMsg := "لطفاً عبارت جستجو را وارد کنید. مثال: /search interstellar"
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

		errorMsg := "خطا در جستجوی محتوا. لطفاً بعداً تلاش کنید."
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Handle empty results
	if len(items) == 0 {
		noResultsMsg := fmt.Sprintf("نتیجه‌ای برای '%s' یافت نشد", query)
		b.SendMessage(ctx, chatID, noResultsMsg)
		return
	}

	// Send each item with poster and formatted message
	for _, item := range items {
		b.sendContentItem(ctx, chatID, &item)
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

	slog.Info("Processing /mutedlist command", "chat_id", chatID)

	// Get all muted series for this user
	mutedSeries, err := b.db.GetMutedSeriesByUser(chatID)
	if err != nil {
		slog.Error("Failed to get muted series",
			"chat_id", chatID,
			"error", err)

		errorMsg := "خطا در دریافت لیست سریال‌های مسدود شده. لطفاً بعداً تلاش کنید."
		b.SendMessage(ctx, chatID, errorMsg)
		return
	}

	// Handle empty list case
	if len(mutedSeries) == 0 {
		emptyMsg := "شما هیچ سریالی را مسدود نکرده‌اید"
		b.SendMessage(ctx, chatID, emptyMsg)
		return
	}

	// Format response message
	var messageText strings.Builder
	messageText.WriteString("سریال‌های مسدود شده:\n\n")

	// Create inline keyboard with unmute button for each series
	var buttons [][]botModels.InlineKeyboardButton

	for i, series := range mutedSeries {
		// Add series to message
		messageText.WriteString(fmt.Sprintf("%d. %s\n", i+1, series.SeriesName))

		// Create unmute button for this series
		buttons = append(buttons, []botModels.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("رفع مسدودیت: %s", series.SeriesName),
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

// sendContentItem sends a single content item with poster and formatted message
func (b *Bot) sendContentItem(ctx context.Context, chatID int64, item *ContentItem) {
	// Format message
	message := FormatContentMessage(item)

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

// FormatContentMessage formats a content item for display
func FormatContentMessage(item *ContentItem) string {
	var message strings.Builder

	// Type indicator and title
	if item.Type == "Movie" {
		message.WriteString("🎬 فیلم\n\n")
		message.WriteString(fmt.Sprintf("نام: %s", item.Name))
	} else if item.Type == "Episode" {
		message.WriteString("📺 قسمت\n\n")
		if item.SeriesName != "" {
			message.WriteString(fmt.Sprintf("سریال: %s\n", item.SeriesName))
		}
		message.WriteString(fmt.Sprintf("فصل %d - قسمت %d", item.SeasonNumber, item.EpisodeNumber))
		if item.Name != "" {
			message.WriteString(fmt.Sprintf("\nنام قسمت: %s", item.Name))
		}
	}

	// Production year
	if item.ProductionYear > 0 {
		message.WriteString(fmt.Sprintf("\nسال: %d", item.ProductionYear))
	}

	// Description
	if item.Overview != "" {
		message.WriteString(fmt.Sprintf("\n\nتوضیحات: %s", item.Overview))
	}

	// Rating
	if item.CommunityRating > 0 {
		message.WriteString(fmt.Sprintf("\n\nامتیاز: %.1f/10", item.CommunityRating))
	} else if item.OfficialRating != "" {
		message.WriteString(fmt.Sprintf("\n\nرده سنی: %s", item.OfficialRating))
	}

	return message.String()
}
