package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jellyfin-telegram-bot/pkg/models"
)

// TestWebhookHandler_ValidMoviePayload tests handling of a valid movie webhook
func TestWebhookHandler_ValidMoviePayload(t *testing.T) {
	// Create a mock database that tracks if content was marked as notified
	db := &MockDB{
		contentNotified: make(map[string]bool),
	}

	handler := NewWebhookHandler(db, "")

	// Create a valid movie payload
	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Movie",
		ItemID:           "movie123",
		ItemName:         "Interstellar",
		Year:             2014,
		Overview:         "A team of explorers travel through a wormhole in space.",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleWebhook(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Assert content was marked as notified
	if !db.contentNotified["movie123"] {
		t.Error("Expected content to be marked as notified")
	}
}

// TestWebhookHandler_ValidEpisodePayload tests handling of a valid episode webhook
func TestWebhookHandler_ValidEpisodePayload(t *testing.T) {
	db := &MockDB{
		contentNotified: make(map[string]bool),
	}

	handler := NewWebhookHandler(db, "")

	// Create a valid episode payload
	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Episode",
		ItemID:           "episode456",
		ItemName:         "Pilot",
		SeriesName:       "Breaking Bad",
		SeasonNumber:     1,
		EpisodeNumber:    1,
		Overview:         "A high school chemistry teacher turned meth cook.",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleWebhook(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Assert content was marked as notified
	if !db.contentNotified["episode456"] {
		t.Error("Expected content to be marked as notified")
	}
}

// TestWebhookHandler_FilterInvalidContentType tests rejection of invalid content types
func TestWebhookHandler_FilterInvalidContentType(t *testing.T) {
	db := &MockDB{
		contentNotified: make(map[string]bool),
	}

	handler := NewWebhookHandler(db, "")

	testCases := []struct {
		name     string
		itemType string
	}{
		{"Series", "Series"},
		{"Season", "Season"},
		{"Audio", "Audio"},
		{"Book", "Book"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload := models.JellyfinWebhook{
				NotificationType: "ItemAdded",
				ItemType:         tc.itemType,
				ItemID:           "test123",
				ItemName:         "Test Item",
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleWebhook(w, req)

			// Should still return 200 but not process
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			// Content should NOT be marked as notified
			if db.contentNotified["test123"] {
				t.Errorf("Content type %s should not be processed", tc.itemType)
			}
		})
	}
}

// TestWebhookHandler_DuplicateContent tests duplicate content detection
func TestWebhookHandler_DuplicateContent(t *testing.T) {
	db := &MockDB{
		contentNotified: make(map[string]bool),
	}

	// Mark content as already notified
	db.contentNotified["movie789"] = true

	handler := NewWebhookHandler(db, "")

	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Movie",
		ItemID:           "movie789",
		ItemName:         "The Matrix",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleWebhook(w, req)

	// Should return 200 (already processed)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// markCount should still be 0 (not called again)
	if db.markCount > 0 {
		t.Error("Duplicate content should not be marked again")
	}
}

// TestWebhookHandler_InvalidJSON tests handling of malformed JSON
func TestWebhookHandler_InvalidJSON(t *testing.T) {
	db := &MockDB{
		contentNotified: make(map[string]bool),
	}

	handler := NewWebhookHandler(db, "")

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleWebhook(w, req)

	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestWebhookHandler_WrongNotificationType tests filtering by notification type
func TestWebhookHandler_WrongNotificationType(t *testing.T) {
	db := &MockDB{
		contentNotified: make(map[string]bool),
	}

	handler := NewWebhookHandler(db, "")

	payload := models.JellyfinWebhook{
		NotificationType: "ItemUpdated", // Not ItemAdded
		ItemType:         "Movie",
		ItemID:           "movie999",
		ItemName:         "Test Movie",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleWebhook(w, req)

	// Should return 200 but not process
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Content should NOT be marked as notified
	if db.contentNotified["movie999"] {
		t.Error("Non-ItemAdded notification should not be processed")
	}
}

// TestWebhookHandler_WithSecret tests webhook security with secret token
func TestWebhookHandler_WithSecret(t *testing.T) {
	db := &MockDB{
		contentNotified: make(map[string]bool),
	}

	secret := "mysecrettoken"
	handler := NewWebhookHandler(db, secret)

	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Movie",
		ItemID:           "movie555",
		ItemName:         "Test Movie",
	}

	body, _ := json.Marshal(payload)

	// Test without secret
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleWebhook(w, req)

	// Should return 401 Unauthorized
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	// Test with correct secret
	req = httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Secret", secret)
	w = httptest.NewRecorder()

	handler.HandleWebhook(w, req)

	// Should return 200
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// MockDB is a mock database for testing
type MockDB struct {
	contentNotified map[string]bool
	markCount       int
}

func (m *MockDB) IsContentNotified(jellyfinID string) (bool, error) {
	return m.contentNotified[jellyfinID], nil
}

func (m *MockDB) MarkContentNotified(jellyfinID, title, contentType string) error {
	m.contentNotified[jellyfinID] = true
	m.markCount++
	return nil
}
