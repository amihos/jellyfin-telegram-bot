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

// TestWebhookToNotificationPipeline tests the complete flow from webhook to notification
func TestWebhookToNotificationPipeline(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_integration_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add a test subscriber
	err = db.AddSubscriber(12345, "testuser", "Test User")
	if err != nil {
		t.Fatalf("Failed to add subscriber: %v", err)
	}

	// Create a mock broadcaster that records calls
	type BroadcastCall struct {
		ItemID   string
		Title    string
		ItemType string
	}
	var broadcastCalls []BroadcastCall
	mockBroadcaster := &mockBroadcaster{
		broadcastFunc: func(ctx context.Context, content *handlers.NotificationContent) error {
			broadcastCalls = append(broadcastCalls, BroadcastCall{
				ItemID:   content.ItemID,
				Title:    content.Title,
				ItemType: content.Type,
			})
			return nil
		},
	}

	// Create webhook handler
	webhookHandler := handlers.NewWebhookHandler(db, "test-secret")
	webhookHandler.SetBroadcaster(mockBroadcaster)

	// Create test webhook payload (movie)
	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Movie",
		ItemName:         "Test Movie",
		Overview:         "A test movie for integration testing",
		Year:             2024,
		ItemID:           "test-item-123",
	}

	payloadBytes, _ := json.Marshal(payload)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Secret", "test-secret")

	// Create response recorder
	w := httptest.NewRecorder()

	// Handle webhook
	webhookHandler.HandleWebhook(w, req)

	// Wait for async broadcast to complete
	time.Sleep(100 * time.Millisecond)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify broadcast was called
	if len(broadcastCalls) != 1 {
		t.Fatalf("Expected 1 broadcast call, got %d", len(broadcastCalls))
	}

	// Verify broadcast content
	call := broadcastCalls[0]
	if call.ItemID != "test-item-123" {
		t.Errorf("Expected ItemID 'test-item-123', got '%s'", call.ItemID)
	}
	if call.Title != "Test Movie" {
		t.Errorf("Expected Title 'Test Movie', got '%s'", call.Title)
	}
	if call.ItemType != "Movie" {
		t.Errorf("Expected ItemType 'Movie', got '%s'", call.ItemType)
	}

	// Verify content was marked as notified in database
	isNotified, err := db.IsContentNotified("test-item-123")
	if err != nil {
		t.Fatalf("Failed to check content notification: %v", err)
	}
	if !isNotified {
		t.Error("Content should be marked as notified")
	}
}

// TestWebhookDuplicatePrevention tests that duplicate webhooks don't trigger multiple notifications
func TestWebhookDuplicatePrevention(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_duplicate_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add a test subscriber
	err = db.AddSubscriber(12345, "testuser", "Test User")
	if err != nil {
		t.Fatalf("Failed to add subscriber: %v", err)
	}

	// Create a mock broadcaster that counts calls
	callCount := 0
	mockBroadcaster := &mockBroadcaster{
		broadcastFunc: func(ctx context.Context, content *handlers.NotificationContent) error {
			callCount++
			return nil
		},
	}

	// Create webhook handler
	webhookHandler := handlers.NewWebhookHandler(db, "test-secret")
	webhookHandler.SetBroadcaster(mockBroadcaster)

	// Create test webhook payload
	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Movie",
		ItemName:         "Duplicate Test Movie",
		ItemID:           "duplicate-123",
	}

	payloadBytes, _ := json.Marshal(payload)

	// Send webhook first time
	req1 := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payloadBytes))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-Webhook-Secret", "test-secret")
	w1 := httptest.NewRecorder()
	webhookHandler.HandleWebhook(w1, req1)

	// Wait for async broadcast
	time.Sleep(100 * time.Millisecond)

	// Send webhook second time (duplicate)
	req2 := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payloadBytes))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Webhook-Secret", "test-secret")
	w2 := httptest.NewRecorder()
	webhookHandler.HandleWebhook(w2, req2)

	// Wait for any potential async broadcast
	time.Sleep(100 * time.Millisecond)

	// Verify only one broadcast occurred
	if callCount != 1 {
		t.Errorf("Expected 1 broadcast call for duplicate webhooks, got %d", callCount)
	}
}

// TestEpisodeNotificationFlow tests the complete flow for TV episode notifications
func TestEpisodeNotificationFlow(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_episode_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add a test subscriber
	err = db.AddSubscriber(12345, "testuser", "Test User")
	if err != nil {
		t.Fatalf("Failed to add subscriber: %v", err)
	}

	// Create a mock broadcaster that records calls
	type BroadcastCall struct {
		Title    string
		ItemType string
	}
	var broadcastCalls []BroadcastCall
	mockBroadcaster := &mockBroadcaster{
		broadcastFunc: func(ctx context.Context, content *handlers.NotificationContent) error {
			broadcastCalls = append(broadcastCalls, BroadcastCall{
				Title:    content.Title,
				ItemType: content.Type,
			})
			return nil
		},
	}

	// Create webhook handler
	webhookHandler := handlers.NewWebhookHandler(db, "test-secret")
	webhookHandler.SetBroadcaster(mockBroadcaster)

	// Create test webhook payload (episode)
	payload := models.JellyfinWebhook{
		NotificationType: "ItemAdded",
		ItemType:         "Episode",
		ItemName:         "Episode Title",
		SeriesName:       "Test Series",
		SeasonNumber:     1,
		EpisodeNumber:    5,
		ItemID:           "episode-123",
	}

	payloadBytes, _ := json.Marshal(payload)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Secret", "test-secret")

	// Create response recorder
	w := httptest.NewRecorder()

	// Handle webhook
	webhookHandler.HandleWebhook(w, req)

	// Wait for async broadcast
	time.Sleep(100 * time.Millisecond)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify broadcast was called
	if len(broadcastCalls) != 1 {
		t.Fatalf("Expected 1 broadcast call, got %d", len(broadcastCalls))
	}

	// Verify episode notification
	call := broadcastCalls[0]
	if call.ItemType != "Episode" {
		t.Errorf("Expected ItemType 'Episode', got '%s'", call.ItemType)
	}
}

// TestJellyfinAPIIntegration tests the integration with Jellyfin API (with mock server)
func TestJellyfinAPIIntegration(t *testing.T) {
	// Create mock Jellyfin server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication header
		if r.Header.Get("X-Emby-Token") != "test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Mock poster image endpoint
		if r.URL.Path == "/Items/test-item-123/Images/Primary" {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write([]byte("fake-image-data"))
			return
		}

		// Mock recent items endpoint
		if r.URL.Path == "/Items" && r.URL.Query().Get("SortBy") == "DateCreated" {
			response := map[string]interface{}{
				"Items": []map[string]interface{}{
					{
						"Id":   "item-1",
						"Name": "Recent Movie",
						"Type": "Movie",
					},
				},
				"TotalRecordCount": 1,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Mock search endpoint
		if r.URL.Path == "/Items" && r.URL.Query().Get("SearchTerm") != "" {
			response := map[string]interface{}{
				"Items": []map[string]interface{}{
					{
						"Id":   "search-1",
						"Name": "Search Result",
						"Type": "Movie",
					},
				},
				"TotalRecordCount": 1,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// Create Jellyfin client pointing to mock server
	client := jellyfin.NewClient(mockServer.URL, "test-api-key")

	// Test image fetching
	imageData, err := client.GetPosterImage(context.Background(), "test-item-123")
	if err != nil {
		t.Errorf("Failed to fetch image: %v", err)
	}
	if len(imageData) == 0 {
		t.Error("Expected image data, got empty response")
	}

	// Test recent items query
	recentItems, err := client.GetRecentItems(context.Background(), 10)
	if err != nil {
		t.Errorf("Failed to get recent items: %v", err)
	}
	if len(recentItems) != 1 {
		t.Errorf("Expected 1 recent item, got %d", len(recentItems))
	}

	// Test search query
	searchResults, err := client.SearchContent(context.Background(), "test", 10)
	if err != nil {
		t.Errorf("Failed to search content: %v", err)
	}
	if len(searchResults) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(searchResults))
	}
}

// Mock broadcaster for testing
type mockBroadcaster struct {
	broadcastFunc func(ctx context.Context, content *handlers.NotificationContent) error
}

func (m *mockBroadcaster) BroadcastNotification(ctx context.Context, content *handlers.NotificationContent) error {
	if m.broadcastFunc != nil {
		return m.broadcastFunc(ctx, content)
	}
	return nil
}
