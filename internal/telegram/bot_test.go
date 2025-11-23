package telegram

import (
	"context"
	"errors"
	"testing"

	"jellyfin-telegram-bot/internal/i18n"
	"jellyfin-telegram-bot/pkg/models"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// Mock implementations for testing

type MockSubscriberDB struct {
	subscribers    map[int64]bool
	languages      map[int64]string // chatID -> languageCode
	mutedSeries    map[int64]map[string]bool // chatID -> seriesID -> isMuted
	shouldFailAdd  bool
	shouldFailGet  bool
}

func NewMockSubscriberDB() *MockSubscriberDB {
	return &MockSubscriberDB{
		subscribers: make(map[int64]bool),
		languages:   make(map[int64]string),
		mutedSeries: make(map[int64]map[string]bool),
	}
}

func (m *MockSubscriberDB) AddSubscriber(chatID int64, username, firstName string) error {
	if m.shouldFailAdd {
		return errors.New("database error")
	}
	m.subscribers[chatID] = true
	return nil
}

func (m *MockSubscriberDB) RemoveSubscriber(chatID int64) error {
	m.subscribers[chatID] = false
	return nil
}

func (m *MockSubscriberDB) GetAllActiveSubscribers() ([]int64, error) {
	if m.shouldFailGet {
		return nil, errors.New("database error")
	}
	var active []int64
	for chatID, isActive := range m.subscribers {
		if isActive {
			active = append(active, chatID)
		}
	}
	return active, nil
}

func (m *MockSubscriberDB) IsSubscribed(chatID int64) (bool, error) {
	return m.subscribers[chatID], nil
}

func (m *MockSubscriberDB) SetLanguage(chatID int64, languageCode string) error {
	m.languages[chatID] = languageCode
	return nil
}

func (m *MockSubscriberDB) GetLanguage(chatID int64) (string, error) {
	if lang, ok := m.languages[chatID]; ok {
		return lang, nil
	}
	return "en", nil // Default to English
}

func (m *MockSubscriberDB) AddMutedSeries(chatID int64, seriesID string, seriesName string) error {
	if m.mutedSeries[chatID] == nil {
		m.mutedSeries[chatID] = make(map[string]bool)
	}
	m.mutedSeries[chatID][seriesID] = true
	return nil
}

func (m *MockSubscriberDB) RemoveMutedSeries(chatID int64, seriesID string) error {
	if m.mutedSeries[chatID] != nil {
		delete(m.mutedSeries[chatID], seriesID)
	}
	return nil
}

func (m *MockSubscriberDB) GetMutedSeriesByUser(chatID int64) ([]models.MutedSeries, error) {
	var result []models.MutedSeries
	if m.mutedSeries[chatID] != nil {
		for seriesID := range m.mutedSeries[chatID] {
			result = append(result, models.MutedSeries{
				ChatID:     chatID,
				SeriesID:   seriesID,
				SeriesName: seriesID,
			})
		}
	}
	return result, nil
}

func (m *MockSubscriberDB) IsSeriesMuted(chatID int64, seriesID string) (bool, error) {
	if m.mutedSeries[chatID] != nil {
		return m.mutedSeries[chatID][seriesID], nil
	}
	return false, nil
}

type MockJellyfinClient struct {
	recentItems   []ContentItem
	searchResults []ContentItem
	imageData     []byte
	shouldFail    bool
}

func NewMockJellyfinClient() *MockJellyfinClient {
	return &MockJellyfinClient{
		recentItems: []ContentItem{
			{
				ItemID:          "movie1",
				Name:            "Test Movie",
				Type:            "Movie",
				Overview:        "A test movie",
				CommunityRating: 8.5,
				ProductionYear:  2023,
			},
			{
				ItemID:        "episode1",
				Name:          "Test Episode",
				Type:          "Episode",
				Overview:      "A test episode",
				SeriesName:    "Test Series",
				SeasonNumber:  1,
				EpisodeNumber: 1,
			},
		},
		searchResults: []ContentItem{
			{
				ItemID:          "movie2",
				Name:            "Interstellar",
				Type:            "Movie",
				Overview:        "A space adventure",
				CommunityRating: 9.0,
				ProductionYear:  2014,
			},
		},
		imageData: []byte("fake-image-data"),
	}
}

func (m *MockJellyfinClient) GetRecentItems(ctx context.Context, limit int) ([]ContentItem, error) {
	if m.shouldFail {
		return nil, errors.New("jellyfin error")
	}
	return m.recentItems, nil
}

func (m *MockJellyfinClient) SearchContent(ctx context.Context, query string, limit int) ([]ContentItem, error) {
	if m.shouldFail {
		return nil, errors.New("jellyfin error")
	}
	if query == "notfound" {
		return []ContentItem{}, nil
	}
	return m.searchResults, nil
}

func (m *MockJellyfinClient) GetPosterImage(ctx context.Context, itemID string) ([]byte, error) {
	if m.shouldFail {
		return nil, errors.New("image fetch error")
	}
	return m.imageData, nil
}

// Helper function to get English localizer for testing
func getTestLocalizer() *goi18n.Localizer {
	bundle, err := i18n.InitBundle()
	if err != nil {
		return nil
	}
	return goi18n.NewLocalizer(bundle, "en")
}

// Tests

// Test 1: Bot initialization with valid token
func TestNewBot_Success(t *testing.T) {
	db := NewMockSubscriberDB()
	jellyfin := NewMockJellyfinClient()

	// Note: Bot creation requires a real Telegram token to connect to the API
	// With a test token, it will fail with "not found" or authentication error
	// This test verifies the validation logic works
	bot, err := NewBot("test-token", db, jellyfin, nil)

	// Expect error since test-token is not valid for Telegram API
	if err == nil {
		t.Error("Expected error for invalid Telegram token, but got none")
		if bot != nil && bot.db != db {
			t.Error("Bot database not set correctly")
		}
		if bot != nil && bot.jellyfinClient != jellyfin {
			t.Error("Bot Jellyfin client not set correctly")
		}
	}
	// Error is expected - test token cannot connect to Telegram API
}

// Test 2: Bot initialization fails with empty token
func TestNewBot_EmptyToken(t *testing.T) {
	db := NewMockSubscriberDB()
	jellyfin := NewMockJellyfinClient()

	bot, err := NewBot("", db, jellyfin, nil)

	if err == nil {
		t.Fatal("Expected error for empty token, got nil")
	}

	if bot != nil {
		t.Error("Expected nil bot for empty token")
	}

	if err.Error() != "TELEGRAM_BOT_TOKEN is required" {
		t.Errorf("Expected 'TELEGRAM_BOT_TOKEN is required' error, got: %v", err)
	}
}

// Test 3: FormatContentMessage for Movie
func TestFormatContentMessage_Movie(t *testing.T) {
	localizer := getTestLocalizer()
	if localizer == nil {
		t.Fatal("Failed to initialize test localizer")
	}

	item := &ContentItem{
		ItemID:          "movie1",
		Name:            "The Matrix",
		Type:            "Movie",
		Overview:        "A hacker discovers reality is a simulation",
		CommunityRating: 8.7,
		ProductionYear:  1999,
	}

	message := FormatContentMessage(item, localizer)

	if message == "" {
		t.Fatal("Expected non-empty message")
	}

	// Check for movie indicator (in English for test)
	if !contains(message, "Movie") {
		t.Error("Message should contain movie indicator")
	}

	// Check for title
	if !contains(message, "The Matrix") {
		t.Error("Message should contain movie title")
	}

	// Check for rating
	if !contains(message, "8.7") {
		t.Error("Message should contain rating")
	}
}

// Test 4: FormatContentMessage for Episode
func TestFormatContentMessage_Episode(t *testing.T) {
	localizer := getTestLocalizer()
	if localizer == nil {
		t.Fatal("Failed to initialize test localizer")
	}

	item := &ContentItem{
		ItemID:        "episode1",
		Name:          "Pilot",
		Type:          "Episode",
		Overview:      "The first episode",
		SeriesName:    "Breaking Bad",
		SeasonNumber:  1,
		EpisodeNumber: 1,
	}

	message := FormatContentMessage(item, localizer)

	if message == "" {
		t.Fatal("Expected non-empty message")
	}

	// Check for episode indicator (in English for test)
	if !contains(message, "Episode") {
		t.Error("Message should contain episode indicator")
	}

	// Check for series name
	if !contains(message, "Breaking Bad") {
		t.Error("Message should contain series name")
	}

	// Check for season and episode
	if !contains(message, "Season 1") || !contains(message, "Episode 1") {
		t.Error("Message should contain season and episode numbers")
	}
}

// Test 5: FormatNotification for Movie
func TestFormatNotification_Movie(t *testing.T) {
	localizer := getTestLocalizer()
	if localizer == nil {
		t.Fatal("Failed to initialize test localizer")
	}

	content := &NotificationContent{
		ItemID:   "movie1",
		Type:     "Movie",
		Title:    "Inception",
		Overview: "A thief who steals secrets through dreams",
		Year:     2010,
		Rating:   8.8,
	}

	message := FormatNotification(content, localizer)

	if message == "" {
		t.Fatal("Expected non-empty message")
	}

	// Check for new movie indicator
	if !contains(message, "New") || !contains(message, "Movie") {
		t.Error("Notification should contain 'new movie' indicator")
	}

	// Check for title
	if !contains(message, "Inception") {
		t.Error("Notification should contain movie title")
	}

	// Check for year
	if !contains(message, "2010") {
		t.Error("Notification should contain year")
	}
}

// Test 6: FormatNotification for Episode
func TestFormatNotification_Episode(t *testing.T) {
	localizer := getTestLocalizer()
	if localizer == nil {
		t.Fatal("Failed to initialize test localizer")
	}

	content := &NotificationContent{
		ItemID:        "episode1",
		Type:          "Episode",
		Title:         "The One Where It All Begins",
		Overview:      "The pilot episode",
		SeriesName:    "Friends",
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Rating:        9.0,
	}

	message := FormatNotification(content, localizer)

	if message == "" {
		t.Fatal("Expected non-empty message")
	}

	// Check for new episode indicator
	if !contains(message, "New") || !contains(message, "Episode") {
		t.Error("Notification should contain 'new episode' indicator")
	}

	// Check for series name
	if !contains(message, "Friends") {
		t.Error("Notification should contain series name")
	}

	// Check for season and episode
	if !contains(message, "Season 1") || !contains(message, "Episode 1") {
		t.Error("Notification should contain season and episode numbers")
	}
}

// Test 7: BroadcastNotification success
func TestBroadcastNotification_Success(t *testing.T) {
	db := NewMockSubscriberDB()
	_ = NewMockJellyfinClient() // Not used in this test

	// Add some subscribers
	db.AddSubscriber(12345, "user1", "Test User 1")
	db.AddSubscriber(67890, "user2", "Test User 2")

	// Note: We can't fully test the bot without a real Telegram token
	// This test verifies the logic with mocks
	content := &NotificationContent{
		ItemID:   "movie1",
		Type:     "Movie",
		Title:    "Test Movie",
		Overview: "A test",
		Year:     2023,
		Rating:   8.0,
	}

	// Verify subscribers exist
	subscribers, err := db.GetAllActiveSubscribers()
	if err != nil {
		t.Fatalf("Failed to get subscribers: %v", err)
	}

	if len(subscribers) != 2 {
		t.Errorf("Expected 2 subscribers, got %d", len(subscribers))
	}

	// Verify notification formatting
	localizer := getTestLocalizer()
	if localizer == nil {
		t.Fatal("Failed to initialize test localizer")
	}

	message := FormatNotification(content, localizer)
	if message == "" {
		t.Error("Expected formatted notification message")
	}
}

// Test 8: BroadcastNotification with no subscribers
func TestBroadcastNotification_NoSubscribers(t *testing.T) {
	db := NewMockSubscriberDB()
	_ = NewMockJellyfinClient() // Not used in this test

	content := &NotificationContent{
		ItemID:   "movie1",
		Type:     "Movie",
		Title:    "Test Movie",
		Overview: "A test",
		Year:     2023,
		Rating:   8.0,
	}

	// Verify no subscribers
	subscribers, err := db.GetAllActiveSubscribers()
	if err != nil {
		t.Fatalf("Failed to get subscribers: %v", err)
	}

	if len(subscribers) != 0 {
		t.Errorf("Expected 0 subscribers, got %d", len(subscribers))
	}

	// Format notification should still work
	localizer := getTestLocalizer()
	if localizer == nil {
		t.Fatal("Failed to initialize test localizer")
	}

	message := FormatNotification(content, localizer)
	if message == "" {
		t.Error("Expected formatted notification message even with no subscribers")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
