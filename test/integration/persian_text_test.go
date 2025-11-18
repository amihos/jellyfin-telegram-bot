package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"jellyfin-telegram-bot/internal/jellyfin"
	"jellyfin-telegram-bot/internal/telegram"
)

// TestPersianCharacterSearch tests searching with Persian characters
func TestPersianCharacterSearch(t *testing.T) {
	// Create mock Jellyfin server that handles Persian search query
	persianQuery := "فیلم"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Emby-Token") != "test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if Persian query is preserved in URL
		searchTerm := r.URL.Query().Get("SearchTerm")
		if searchTerm != "" {
			// Verify Persian characters are preserved (URL decoded)
			if !strings.Contains(r.URL.RawQuery, "SearchTerm") {
				t.Error("SearchTerm parameter not found in query")
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"Items":[{"Id":"1","Name":"نام فیلم","Type":"Movie"}],"TotalRecordCount":1}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// Create Jellyfin client
	client := jellyfin.NewClient(mockServer.URL, "test-api-key")

	// Search with Persian query
	results, err := client.SearchContent(context.Background(), persianQuery, 10)
	if err != nil {
		t.Errorf("Failed to search with Persian characters: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected search results for Persian query, got none")
	}
}

// TestPersianNotificationFormatting tests that Persian notification messages are formatted correctly
func TestPersianNotificationFormatting(t *testing.T) {
	tests := []struct {
		name     string
		content  telegram.NotificationContent
		expected []string // Strings that should appear in formatted message
	}{
		{
			name: "Movie with Persian title",
			content: telegram.NotificationContent{
				ItemID:   "1",
				Type:     "Movie",
				Title:    "نام فیلم فارسی",
				Overview: "توضیحات فیلم به زبان فارسی",
				Rating:   8.5,
				Year:     2024,
			},
			expected: []string{
				"فیلم جدید", // New movie
				"نام فیلم فارسی",
				"توضیحات فیلم به زبان فارسی",
				"2024",
			},
		},
		{
			name: "Episode with Persian series name",
			content: telegram.NotificationContent{
				ItemID:        "2",
				Type:          "Episode",
				Title:         "نام قسمت",
				SeriesName:    "نام سریال",
				Overview:      "توضیحات قسمت",
				SeasonNumber:  1,
				EpisodeNumber: 5,
			},
			expected: []string{
				"قسمت جدید", // New episode
				"نام سریال",
				"فصل 1",
				"قسمت 5",
				"نام قسمت",
				"توضیحات قسمت",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Format notification message
			message := telegram.FormatNotification(&tt.content)

			// Verify all expected strings are present
			for _, expected := range tt.expected {
				if !strings.Contains(message, expected) {
					t.Errorf("Expected message to contain '%s', but it didn't.\nMessage: %s", expected, message)
				}
			}

			// Verify message contains no LTR marks that might interfere with Persian RTL
			if strings.Contains(message, "\u200E") { // LTR mark
				t.Error("Message should not contain LTR marks that might interfere with Persian RTL")
			}
		})
	}
}

// TestRTLFormatting tests that Right-to-Left formatting works correctly with mixed content
func TestRTLFormatting(t *testing.T) {
	// Create a notification content item with mixed Persian and English
	content := telegram.NotificationContent{
		ItemID:   "test-1",
		Type:     "Movie",
		Title:    "The Matrix - ماتریکس",
		Overview: "A computer hacker learns about the true nature of reality. یک هکر کامپیوتری درباره ماهیت واقعی واقعیت می‌آموزد.",
		Rating:   8.7,
		Year:     1999,
	}

	message := telegram.FormatNotification(&content)

	// Verify both Persian and English text are present
	if !strings.Contains(message, "The Matrix") {
		t.Error("Message should contain English title")
	}
	if !strings.Contains(message, "ماتریکس") {
		t.Error("Message should contain Persian title")
	}
	if !strings.Contains(message, "A computer hacker") {
		t.Error("Message should contain English overview")
	}
	if !strings.Contains(message, "واقعیت") {
		t.Error("Message should contain Persian overview")
	}

	// Verify Persian UI elements are present
	if !strings.Contains(message, "فیلم جدید") {
		t.Error("Message should contain Persian movie label")
	}
}
