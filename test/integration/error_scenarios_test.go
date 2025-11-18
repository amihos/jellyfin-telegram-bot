package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"jellyfin-telegram-bot/internal/database"
	"jellyfin-telegram-bot/internal/handlers"
	"jellyfin-telegram-bot/internal/jellyfin"
	"jellyfin-telegram-bot/pkg/models"
)

// TestJellyfinAPIErrors tests handling of various Jellyfin API errors
func TestJellyfinAPIErrors(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
	}{
		{
			name:          "Unauthorized (401)",
			statusCode:    http.StatusUnauthorized,
			expectedError: true,
		},
		{
			name:          "Not Found (404)",
			statusCode:    http.StatusNotFound,
			expectedError: true,
		},
		{
			name:          "Server Error (500)",
			statusCode:    http.StatusInternalServerError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server that returns error status
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer mockServer.Close()

			// Create client
			client := jellyfin.NewClient(mockServer.URL, "test-api-key")

			// Try to fetch image
			_, err := client.GetPosterImage(context.Background(), "test-item")

			// Verify error occurred
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

// TestWebhookInvalidPayloads tests handling of invalid webhook payloads
func TestWebhookInvalidPayloads(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_invalid_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Create webhook handler
	webhookHandler := handlers.NewWebhookHandler(db, "test-secret")

	tests := []struct {
		name       string
		payload    string
		statusCode int
	}{
		{
			name:       "Invalid JSON",
			payload:    `{invalid json}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Empty payload",
			payload:    `{}`,
			statusCode: http.StatusOK, // Empty payload is valid but filtered out
		},
		{
			name: "Invalid item type (Audio)",
			payload: `{
				"NotificationType": "ItemAdded",
				"ItemType": "Audio",
				"Name": "Test Audio"
			}`,
			statusCode: http.StatusOK, // Filtered out but not an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/webhook", bytes.NewReader([]byte(tt.payload)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Webhook-Secret", "test-secret")

			w := httptest.NewRecorder()
			webhookHandler.HandleWebhook(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

// TestWebhookSecurityValidation tests webhook security validation
func TestWebhookSecurityValidation(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_security_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Create webhook handler with secret
	webhookHandler := handlers.NewWebhookHandler(db, "correct-secret")

	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Movie",
		ItemName:         "Test Movie",
		ItemID:           "test-123",
	}
	payloadBytes, _ := json.Marshal(payload)

	tests := []struct {
		name       string
		secret     string
		statusCode int
	}{
		{
			name:       "Correct secret",
			secret:     "correct-secret",
			statusCode: http.StatusOK,
		},
		{
			name:       "Wrong secret",
			secret:     "wrong-secret",
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "Missing secret",
			secret:     "",
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.secret != "" {
				req.Header.Set("X-Webhook-Secret", tt.secret)
			}

			w := httptest.NewRecorder()
			webhookHandler.HandleWebhook(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
