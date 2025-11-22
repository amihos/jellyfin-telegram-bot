package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"jellyfin-telegram-bot/pkg/models"
)

// ContentTracker defines the interface for content tracking operations
type ContentTracker interface {
	IsContentNotified(jellyfinID string) (bool, error)
	MarkContentNotified(jellyfinID, title, contentType string) error
}

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

// NotificationBroadcaster defines the interface for broadcasting notifications
type NotificationBroadcaster interface {
	BroadcastNotification(ctx context.Context, content *NotificationContent) error
}

// WebhookHandler handles incoming Jellyfin webhook requests
type WebhookHandler struct {
	db          ContentTracker
	secret      string
	broadcaster NotificationBroadcaster
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(db ContentTracker, secret string) *WebhookHandler {
	return &WebhookHandler{
		db:          db,
		secret:      secret,
		broadcaster: nil,
	}
}

// SetBroadcaster sets the notification broadcaster
func (h *WebhookHandler) SetBroadcaster(broadcaster NotificationBroadcaster) {
	h.broadcaster = broadcaster
}

// HandleWebhook processes incoming webhook requests from Jellyfin
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Validate request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate webhook secret if configured
	if h.secret != "" {
		providedSecret := r.Header.Get("X-Webhook-Secret")
		if providedSecret != h.secret {
			slog.Warn("Webhook request with invalid or missing secret",
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Parse webhook payload
	var payload models.JellyfinWebhook
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("Failed to parse webhook payload",
			"error", err,
			"remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Log received webhook
	slog.Info("Received webhook",
		"notification_type", payload.NotificationType,
		"item_type", payload.ItemType,
		"item_id", payload.ItemID,
		"item_name", payload.ItemName)

	// Validate webhook - must be ItemAdded and Movie or Episode
	if !payload.IsValid() {
		slog.Debug("Webhook ignored - invalid content",
			"notification_type", payload.NotificationType,
			"item_type", payload.ItemType,
			"item_id", payload.ItemID)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if content already notified
	notified, err := h.db.IsContentNotified(payload.ItemID)
	if err != nil {
		slog.Error("Failed to check content notification status",
			"error", err,
			"item_id", payload.ItemID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if notified {
		slog.Info("Content already notified, skipping",
			"item_id", payload.ItemID,
			"item_name", payload.ItemName)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract metadata
	metadata := h.extractMetadata(&payload)

	// Log what will be notified
	slog.Info("New content ready for notification",
		"item_id", payload.ItemID,
		"type", metadata.Type,
		"title", metadata.Title,
		"year", metadata.Year)

	if payload.IsEpisode() {
		slog.Info("Episode details",
			"series_name", metadata.SeriesName,
			"season", metadata.SeasonNumber,
			"episode", metadata.EpisodeNumber)
	}

	// Mark content as notified to prevent duplicates
	contentType := "Movie"
	if payload.IsEpisode() {
		contentType = "Episode"
	}

	if err := h.db.MarkContentNotified(payload.ItemID, payload.ItemName, contentType); err != nil {
		slog.Error("Failed to mark content as notified",
			"error", err,
			"item_id", payload.ItemID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Content marked as notified",
		"item_id", payload.ItemID,
		"item_name", payload.ItemName)

	// Broadcast notification to subscribers
	if h.broadcaster != nil {
		// Extract all metadata for notification
		content := &NotificationContent{
			ItemID:        payload.ItemID,
			Type:          contentType,
			Title:         metadata.Title,  // Use metadata.Title which has "Unknown" fallback
			Overview:      metadata.Overview,  // Use metadata.Overview for consistency
			Year:          payload.Year,
			Rating:        0, // Webhook doesn't include rating - could fetch from API if needed
			SeriesName:    metadata.SeriesName,  // Use metadata.SeriesName for "Unknown Series" fallback
			SeasonNumber:  payload.SeasonNumber,
			EpisodeNumber: payload.EpisodeNumber,
		}

		// Broadcast asynchronously to avoid blocking webhook response
		go func() {
			ctx := context.Background()
			if err := h.broadcaster.BroadcastNotification(ctx, content); err != nil {
				slog.Error("Failed to broadcast notification",
					"item_id", payload.ItemID,
					"error", err)
			}
		}()

		slog.Info("Notification broadcast initiated",
			"item_id", payload.ItemID)
	} else {
		slog.Warn("No broadcaster configured, notification not sent")
	}

	w.WriteHeader(http.StatusOK)
}

// ContentMetadata represents extracted metadata for notifications
type ContentMetadata struct {
	Type          string
	Title         string
	Overview      string
	Year          int
	ItemID        string
	SeriesName    string
	SeasonNumber  int
	EpisodeNumber int
}

// extractMetadata extracts relevant metadata from webhook payload
func (h *WebhookHandler) extractMetadata(payload *models.JellyfinWebhook) *ContentMetadata {
	metadata := &ContentMetadata{
		Title:    payload.ItemName,
		Overview: payload.Overview,
		Year:     payload.Year,
		ItemID:   payload.ItemID,
	}

	if payload.IsMovie() {
		metadata.Type = "Movie"
	} else if payload.IsEpisode() {
		metadata.Type = "Episode"
		metadata.SeriesName = payload.SeriesName
		metadata.SeasonNumber = payload.SeasonNumber
		metadata.EpisodeNumber = payload.EpisodeNumber
	}

	// Handle missing fields gracefully
	if metadata.Title == "" {
		metadata.Title = "Unknown"
	}
	if metadata.Overview == "" {
		metadata.Overview = "No description available"
	}
	if payload.IsEpisode() && metadata.SeriesName == "" {
		metadata.SeriesName = "Unknown Series"
	}

	return metadata
}

// StartWebhookServer starts the HTTP server for webhook endpoint
func StartWebhookServer(port string, handler *WebhookHandler) error {
	http.HandleFunc("/webhook", handler.HandleWebhook)
	http.HandleFunc("/health", HealthCheckHandler)

	addr := fmt.Sprintf(":%s", port)
	slog.Info("Starting webhook server", "address", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("webhook server failed: %w", err)
	}

	return nil
}
