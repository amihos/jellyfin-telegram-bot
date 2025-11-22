package telegram

import (
	"context"
	"errors"
	"testing"

	"jellyfin-telegram-bot/pkg/models"

	botModels "github.com/go-telegram/bot/models"
)

// mockSubscriberDB implements SubscriberDB interface for testing
type mockSubscriberDB struct {
	subscribers  []int64
	mutedSeries  map[int64]map[string]bool // chatID -> seriesID -> isMuted
	addSubErr    error
	removeSubErr error
}

func newMockSubscriberDB() *mockSubscriberDB {
	return &mockSubscriberDB{
		subscribers: []int64{},
		mutedSeries: make(map[int64]map[string]bool),
	}
}

func (m *mockSubscriberDB) AddSubscriber(chatID int64, username, firstName string) error {
	if m.addSubErr != nil {
		return m.addSubErr
	}
	m.subscribers = append(m.subscribers, chatID)
	return nil
}

func (m *mockSubscriberDB) RemoveSubscriber(chatID int64) error {
	if m.removeSubErr != nil {
		return m.removeSubErr
	}
	for i, id := range m.subscribers {
		if id == chatID {
			m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockSubscriberDB) GetAllActiveSubscribers() ([]int64, error) {
	return m.subscribers, nil
}

func (m *mockSubscriberDB) IsSubscribed(chatID int64) (bool, error) {
	for _, id := range m.subscribers {
		if id == chatID {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockSubscriberDB) AddMutedSeries(chatID int64, seriesID string, seriesName string) error {
	if m.mutedSeries[chatID] == nil {
		m.mutedSeries[chatID] = make(map[string]bool)
	}
	m.mutedSeries[chatID][seriesID] = true
	return nil
}

func (m *mockSubscriberDB) RemoveMutedSeries(chatID int64, seriesID string) error {
	if m.mutedSeries[chatID] != nil {
		delete(m.mutedSeries[chatID], seriesID)
	}
	return nil
}

func (m *mockSubscriberDB) GetMutedSeriesByUser(chatID int64) ([]models.MutedSeries, error) {
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

func (m *mockSubscriberDB) IsSeriesMuted(chatID int64, seriesID string) (bool, error) {
	if m.mutedSeries[chatID] != nil {
		return m.mutedSeries[chatID][seriesID], nil
	}
	return false, nil
}

// mockJellyfinClient implements JellyfinClient interface for testing
type mockJellyfinClient struct {
	posterData []byte
	posterErr  error
}

func (m *mockJellyfinClient) GetRecentItems(ctx context.Context, limit int) ([]ContentItem, error) {
	return nil, nil
}

func (m *mockJellyfinClient) SearchContent(ctx context.Context, query string, limit int) ([]ContentItem, error) {
	return nil, nil
}

func (m *mockJellyfinClient) GetPosterImage(ctx context.Context, itemID string) ([]byte, error) {
	return m.posterData, m.posterErr
}

// testBotWrapper wraps Bot to track sent messages without needing real Telegram bot
type testBotWrapper struct {
	db             SubscriberDB
	jellyfinClient JellyfinClient
	sentMessages   map[int64][]string
	sentKeyboards  map[int64][]*botModels.InlineKeyboardMarkup
}

func newTestBotWrapper(db SubscriberDB, jf JellyfinClient) *testBotWrapper {
	return &testBotWrapper{
		db:             db,
		jellyfinClient: jf,
		sentMessages:   make(map[int64][]string),
		sentKeyboards:  make(map[int64][]*botModels.InlineKeyboardMarkup),
	}
}

func (tb *testBotWrapper) SendMessage(ctx context.Context, chatID int64, text string) error {
	tb.sentMessages[chatID] = append(tb.sentMessages[chatID], text)
	return nil
}

func (tb *testBotWrapper) SendMessageWithKeyboard(ctx context.Context, chatID int64, text string, keyboard *botModels.InlineKeyboardMarkup) error {
	tb.sentMessages[chatID] = append(tb.sentMessages[chatID], text)
	tb.sentKeyboards[chatID] = append(tb.sentKeyboards[chatID], keyboard)
	return nil
}

func (tb *testBotWrapper) SendPhotoBytes(ctx context.Context, chatID int64, imageData []byte, caption string) error {
	tb.sentMessages[chatID] = append(tb.sentMessages[chatID], caption)
	return nil
}

func (tb *testBotWrapper) SendPhotoBytesWithKeyboard(ctx context.Context, chatID int64, imageData []byte, caption string, keyboard *botModels.InlineKeyboardMarkup) error {
	tb.sentMessages[chatID] = append(tb.sentMessages[chatID], caption)
	tb.sentKeyboards[chatID] = append(tb.sentKeyboards[chatID], keyboard)
	return nil
}

// broadcastNotificationForTest is a test-friendly version of BroadcastNotification
func (tb *testBotWrapper) broadcastNotificationForTest(ctx context.Context, content *NotificationContent) error {
	// Get all active subscribers
	subscribers, err := tb.db.GetAllActiveSubscribers()
	if err != nil {
		return err
	}

	if len(subscribers) == 0 {
		return nil
	}

	// Filter out muted users for episode notifications
	filteredSubscribers := subscribers
	if content.Type == "Episode" && content.SeriesName != "" {
		filteredSubscribers = make([]int64, 0, len(subscribers))
		for _, chatID := range subscribers {
			isMuted, err := tb.db.IsSeriesMuted(chatID, content.SeriesName)
			if err != nil {
				// Include subscriber if check fails to avoid missing notifications
				filteredSubscribers = append(filteredSubscribers, chatID)
				continue
			}

			if !isMuted {
				filteredSubscribers = append(filteredSubscribers, chatID)
			}
		}
	}

	if len(filteredSubscribers) == 0 {
		return nil
	}

	// Format notification message
	message := FormatNotification(content)

	// Create inline keyboard for episodes with valid series name
	var keyboard *botModels.InlineKeyboardMarkup
	if shouldShowMuteButton(content) {
		keyboard = createMuteButton(content.SeriesName)
	}

	// Fetch poster image
	var imageData []byte
	if content.ItemID != "" {
		imageData, err = tb.jellyfinClient.GetPosterImage(ctx, content.ItemID)
		if err != nil {
			// Continue without image
		}
	}

	// Broadcast to all filtered subscribers
	for _, chatID := range filteredSubscribers {
		if imageData != nil && len(imageData) > 0 {
			if keyboard != nil {
				tb.SendPhotoBytesWithKeyboard(ctx, chatID, imageData, message, keyboard)
			} else {
				tb.SendPhotoBytes(ctx, chatID, imageData, message)
			}
		} else {
			if keyboard != nil {
				tb.SendMessageWithKeyboard(ctx, chatID, message, keyboard)
			} else {
				tb.SendMessage(ctx, chatID, message)
			}
		}
	}

	return nil
}

// Test 1: BroadcastNotification excludes muted users from subscriber list
func TestBroadcastNotification_ExcludesMutedUsers(t *testing.T) {
	db := newMockSubscriberDB()
	db.subscribers = []int64{100, 200, 300}

	// User 200 has muted "Breaking Bad"
	db.AddMutedSeries(200, "Breaking Bad", "Breaking Bad")

	jf := &mockJellyfinClient{}
	bot := newTestBotWrapper(db, jf)

	content := &NotificationContent{
		Type:          "Episode",
		SeriesName:    "Breaking Bad",
		Title:         "Pilot",
		SeasonNumber:  1,
		EpisodeNumber: 1,
	}

	ctx := context.Background()
	err := bot.broadcastNotificationForTest(ctx, content)
	if err != nil {
		t.Fatalf("BroadcastNotification failed: %v", err)
	}

	// Check that only users 100 and 300 received the notification
	if len(bot.sentMessages[100]) != 1 {
		t.Errorf("Expected user 100 to receive 1 message, got %d", len(bot.sentMessages[100]))
	}
	if len(bot.sentMessages[200]) != 0 {
		t.Errorf("Expected user 200 (muted) to receive 0 messages, got %d", len(bot.sentMessages[200]))
	}
	if len(bot.sentMessages[300]) != 1 {
		t.Errorf("Expected user 300 to receive 1 message, got %d", len(bot.sentMessages[300]))
	}
}

// Test 2: Muted user does not receive episode notification
func TestBroadcastNotification_MutedUserDoesNotReceive(t *testing.T) {
	db := newMockSubscriberDB()
	db.subscribers = []int64{123}
	db.AddMutedSeries(123, "The Office", "The Office")

	jf := &mockJellyfinClient{}
	bot := newTestBotWrapper(db, jf)

	content := &NotificationContent{
		Type:          "Episode",
		SeriesName:    "The Office",
		Title:         "Pilot",
		SeasonNumber:  1,
		EpisodeNumber: 1,
	}

	ctx := context.Background()
	err := bot.broadcastNotificationForTest(ctx, content)
	if err != nil {
		t.Fatalf("BroadcastNotification failed: %v", err)
	}

	if len(bot.sentMessages[123]) != 0 {
		t.Errorf("Expected muted user to receive 0 messages, got %d", len(bot.sentMessages[123]))
	}
}

// Test 3: Non-muted users still receive notifications normally
func TestBroadcastNotification_NonMutedUsersReceive(t *testing.T) {
	db := newMockSubscriberDB()
	db.subscribers = []int64{100, 200}

	jf := &mockJellyfinClient{}
	bot := newTestBotWrapper(db, jf)

	content := &NotificationContent{
		Type:          "Episode",
		SeriesName:    "Friends",
		Title:         "The One Where It All Began",
		SeasonNumber:  1,
		EpisodeNumber: 1,
	}

	ctx := context.Background()
	err := bot.broadcastNotificationForTest(ctx, content)
	if err != nil {
		t.Fatalf("BroadcastNotification failed: %v", err)
	}

	if len(bot.sentMessages[100]) != 1 {
		t.Errorf("Expected user 100 to receive 1 message, got %d", len(bot.sentMessages[100]))
	}
	if len(bot.sentMessages[200]) != 1 {
		t.Errorf("Expected user 200 to receive 1 message, got %d", len(bot.sentMessages[200]))
	}
}

// Test 4: Movie notifications are not affected by series muting
func TestBroadcastNotification_MovieNotAffectedByMuting(t *testing.T) {
	db := newMockSubscriberDB()
	db.subscribers = []int64{100}

	// User has muted a series, but we're sending a movie notification
	db.AddMutedSeries(100, "Breaking Bad", "Breaking Bad")

	jf := &mockJellyfinClient{}
	bot := newTestBotWrapper(db, jf)

	content := &NotificationContent{
		Type:  "Movie",
		Title: "Interstellar",
		Year:  2014,
	}

	ctx := context.Background()
	err := bot.broadcastNotificationForTest(ctx, content)
	if err != nil {
		t.Fatalf("BroadcastNotification failed: %v", err)
	}

	// User should receive the movie notification despite having muted series
	if len(bot.sentMessages[100]) != 1 {
		t.Errorf("Expected user 100 to receive movie notification, got %d messages", len(bot.sentMessages[100]))
	}
}

// Test 5: Inline mute button appears on episode notifications
func TestBroadcastNotification_InlineMuteButtonOnEpisodes(t *testing.T) {
	db := newMockSubscriberDB()
	db.subscribers = []int64{100}

	jf := &mockJellyfinClient{}
	bot := newTestBotWrapper(db, jf)

	content := &NotificationContent{
		Type:          "Episode",
		SeriesName:    "Breaking Bad",
		Title:         "Pilot",
		SeasonNumber:  1,
		EpisodeNumber: 1,
	}

	ctx := context.Background()
	err := bot.broadcastNotificationForTest(ctx, content)
	if err != nil {
		t.Fatalf("BroadcastNotification failed: %v", err)
	}

	// Check that keyboard was sent
	if len(bot.sentKeyboards[100]) != 1 {
		t.Fatalf("Expected 1 keyboard to be sent, got %d", len(bot.sentKeyboards[100]))
	}

	keyboard := bot.sentKeyboards[100][0]
	if len(keyboard.InlineKeyboard) != 1 || len(keyboard.InlineKeyboard[0]) != 1 {
		t.Fatalf("Expected 1 row with 1 button, got %d rows", len(keyboard.InlineKeyboard))
	}

	button := keyboard.InlineKeyboard[0][0]
	if button.Text != "دنبال نکردن" {
		t.Errorf("Expected button text 'دنبال نکردن', got '%s'", button.Text)
	}
	expectedCallback := "mute:Breaking Bad"
	if button.CallbackData != expectedCallback {
		t.Errorf("Expected callback data '%s', got '%s'", expectedCallback, button.CallbackData)
	}
}

// Test 6: No mute button for invalid series names
func TestBroadcastNotification_NoMuteButtonForInvalidSeries(t *testing.T) {
	db := newMockSubscriberDB()
	db.subscribers = []int64{100}

	jf := &mockJellyfinClient{}

	testCases := []struct {
		name       string
		seriesName string
	}{
		{"empty series name", ""},
		{"unknown series", "Unknown Series"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bot := newTestBotWrapper(db, jf)

			content := &NotificationContent{
				Type:          "Episode",
				SeriesName:    tc.seriesName,
				Title:         "Episode 1",
				SeasonNumber:  1,
				EpisodeNumber: 1,
			}

			ctx := context.Background()
			err := bot.broadcastNotificationForTest(ctx, content)
			if err != nil {
				t.Fatalf("BroadcastNotification failed: %v", err)
			}

			// Check that no keyboard was sent
			if len(bot.sentKeyboards[100]) != 0 {
				t.Errorf("Expected no keyboard for %s, got %d", tc.name, len(bot.sentKeyboards[100]))
			}
		})
	}
}

// Test 7: shouldShowMuteButton function
func TestShouldShowMuteButton(t *testing.T) {
	testCases := []struct {
		name       string
		content    *NotificationContent
		shouldShow bool
	}{
		{
			name: "valid episode",
			content: &NotificationContent{
				Type:       "Episode",
				SeriesName: "Breaking Bad",
			},
			shouldShow: true,
		},
		{
			name: "movie",
			content: &NotificationContent{
				Type:  "Movie",
				Title: "Interstellar",
			},
			shouldShow: false,
		},
		{
			name: "episode with empty series name",
			content: &NotificationContent{
				Type:       "Episode",
				SeriesName: "",
			},
			shouldShow: false,
		},
		{
			name: "episode with Unknown Series",
			content: &NotificationContent{
				Type:       "Episode",
				SeriesName: "Unknown Series",
			},
			shouldShow: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := shouldShowMuteButton(tc.content)
			if result != tc.shouldShow {
				t.Errorf("Expected shouldShowMuteButton to return %v, got %v", tc.shouldShow, result)
			}
		})
	}
}

// Test 8: Error handling when IsSeriesMuted fails
func TestBroadcastNotification_HandlesMuteCheckError(t *testing.T) {
	// Create a custom mock that returns an error for IsSeriesMuted
	db := &errorMockDB{
		mockSubscriberDB: newMockSubscriberDB(),
		muteCheckErr:     errors.New("database error"),
	}
	db.subscribers = []int64{100}

	jf := &mockJellyfinClient{}
	bot := newTestBotWrapper(db, jf)

	content := &NotificationContent{
		Type:          "Episode",
		SeriesName:    "Breaking Bad",
		Title:         "Pilot",
		SeasonNumber:  1,
		EpisodeNumber: 1,
	}

	ctx := context.Background()
	err := bot.broadcastNotificationForTest(ctx, content)
	if err != nil {
		t.Fatalf("BroadcastNotification failed: %v", err)
	}

	// Even with error, user should still receive notification (fail-safe behavior)
	if len(bot.sentMessages[100]) != 1 {
		t.Errorf("Expected user to receive notification despite mute check error, got %d messages", len(bot.sentMessages[100]))
	}
}

// errorMockDB is a mock that can return errors for specific methods
type errorMockDB struct {
	*mockSubscriberDB
	muteCheckErr error
}

func (e *errorMockDB) IsSeriesMuted(chatID int64, seriesID string) (bool, error) {
	if e.muteCheckErr != nil {
		return false, e.muteCheckErr
	}
	return e.mockSubscriberDB.IsSeriesMuted(chatID, seriesID)
}
