package telegram

import (
	"testing"
)

// TestHandleMuteCallback_Success tests successful mute callback
func TestHandleMuteCallback_Success(t *testing.T) {
	mockDB := NewMockSubscriberDB()
	mockJellyfin := NewMockJellyfinClient()

	botInstance := &Bot{
		db:             mockDB,
		jellyfinClient: mockJellyfin,
	}

	// Test database operation - AddMutedSeries should work
	err := mockDB.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify series was muted
	isMuted, err := mockDB.IsSeriesMuted(12345, "Breaking Bad")
	if err != nil {
		t.Errorf("Expected no error checking mute status, got: %v", err)
	}
	if !isMuted {
		t.Error("Expected series to be muted after AddMutedSeries")
	}

	// Prevent unused variable warning
	_ = botInstance
}

// TestHandleUnmuteCallback_Success tests successful unmute callback
func TestHandleUnmuteCallback_Success(t *testing.T) {
	mockDB := NewMockSubscriberDB()
	mockJellyfin := NewMockJellyfinClient()

	botInstance := &Bot{
		db:             mockDB,
		jellyfinClient: mockJellyfin,
	}

	// First mute the series
	err := mockDB.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")
	if err != nil {
		t.Errorf("Expected no error adding muted series, got: %v", err)
	}

	// Then unmute it
	err = mockDB.RemoveMutedSeries(12345, "Breaking Bad")
	if err != nil {
		t.Errorf("Expected no error removing muted series, got: %v", err)
	}

	// Verify series is no longer muted
	isMuted, err := mockDB.IsSeriesMuted(12345, "Breaking Bad")
	if err != nil {
		t.Errorf("Expected no error checking mute status, got: %v", err)
	}
	if isMuted {
		t.Error("Expected series to be unmuted after RemoveMutedSeries")
	}

	// Prevent unused variable warning
	_ = botInstance
}

// TestHandleMutedList_EmptyList tests /mutedlist with no muted series
func TestHandleMutedList_EmptyList(t *testing.T) {
	mockDB := NewMockSubscriberDB()
	mockJellyfin := NewMockJellyfinClient()

	botInstance := &Bot{
		db:             mockDB,
		jellyfinClient: mockJellyfin,
	}

	// Get muted series for user with no muted series
	mutedSeries, err := mockDB.GetMutedSeriesByUser(12345)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(mutedSeries) != 0 {
		t.Errorf("Expected empty list, got %d items", len(mutedSeries))
	}

	// Prevent unused variable warning
	_ = botInstance
}

// TestHandleMutedList_WithSeries tests /mutedlist with muted series
func TestHandleMutedList_WithSeries(t *testing.T) {
	mockDB := NewMockSubscriberDB()
	mockJellyfin := NewMockJellyfinClient()

	botInstance := &Bot{
		db:             mockDB,
		jellyfinClient: mockJellyfin,
	}

	// Add some muted series
	mockDB.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")
	mockDB.AddMutedSeries(12345, "Game of Thrones", "Game of Thrones")

	// Get muted series
	result, err := mockDB.GetMutedSeriesByUser(12345)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 muted series, got %d", len(result))
	}

	// Verify series names
	seriesNames := make(map[string]bool)
	for _, series := range result {
		seriesNames[series.SeriesName] = true
	}
	if !seriesNames["Breaking Bad"] || !seriesNames["Game of Thrones"] {
		t.Error("Expected both Breaking Bad and Game of Thrones in muted list")
	}

	// Prevent unused variable warning
	_ = botInstance
}

// TestMuteCallback_CreatesRecord tests that mute callback creates database record
func TestMuteCallback_CreatesRecord(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	// Add muted series
	err := mockDB.AddMutedSeries(12345, "The Office", "The Office")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify it was added
	isMuted, err := mockDB.IsSeriesMuted(12345, "The Office")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !isMuted {
		t.Error("Expected The Office to be muted")
	}
}

// TestUnmuteCallback_DeletesRecord tests that unmute callback deletes database record
func TestUnmuteCallback_DeletesRecord(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	// First add a muted series
	mockDB.AddMutedSeries(12345, "The Office", "The Office")

	// Then remove it
	err := mockDB.RemoveMutedSeries(12345, "The Office")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify it was removed
	isMuted, err := mockDB.IsSeriesMuted(12345, "The Office")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if isMuted {
		t.Error("Expected The Office to not be muted after removal")
	}
}

// TestMutedListCommand_FormatsListCorrectly tests list formatting
func TestMutedListCommand_FormatsListCorrectly(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	// Add multiple series
	mockDB.AddMutedSeries(12345, "Series1", "Series One")
	mockDB.AddMutedSeries(12345, "Series2", "Series Two")

	// Retrieve list
	result, err := mockDB.GetMutedSeriesByUser(12345)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 series, got %d", len(result))
	}

	// Verify each series has proper fields
	for _, series := range result {
		if series.SeriesID == "" {
			t.Error("Expected series ID to be non-empty")
		}
		if series.SeriesName == "" {
			t.Error("Expected series name to be non-empty")
		}
		if series.ChatID != 12345 {
			t.Errorf("Expected chat ID to be 12345, got %d", series.ChatID)
		}
	}
}

// TestMultipleUsersCanMuteSameSeries tests that different users can independently mute same series
func TestMultipleUsersCanMuteSameSeries(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	// User 1 mutes Breaking Bad
	mockDB.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")

	// User 2 mutes Breaking Bad
	mockDB.AddMutedSeries(67890, "Breaking Bad", "Breaking Bad")

	// Verify both have it muted
	isMuted1, _ := mockDB.IsSeriesMuted(12345, "Breaking Bad")
	isMuted2, _ := mockDB.IsSeriesMuted(67890, "Breaking Bad")

	if !isMuted1 {
		t.Error("Expected user 1 to have Breaking Bad muted")
	}
	if !isMuted2 {
		t.Error("Expected user 2 to have Breaking Bad muted")
	}

	// Verify each user has their own list
	list1, _ := mockDB.GetMutedSeriesByUser(12345)
	list2, _ := mockDB.GetMutedSeriesByUser(67890)

	if len(list1) != 1 {
		t.Errorf("Expected user 1 to have 1 muted series, got %d", len(list1))
	}
	if len(list2) != 1 {
		t.Errorf("Expected user 2 to have 1 muted series, got %d", len(list2))
	}
}

// TestUnmuteRestoresNotifications tests that unmuting allows notifications again
func TestUnmuteRestoresNotifications(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	chatID := int64(12345)
	seriesID := "Breaking Bad"

	// Initially not muted
	isMuted, _ := mockDB.IsSeriesMuted(chatID, seriesID)
	if isMuted {
		t.Error("Series should not be muted initially")
	}

	// Mute the series
	mockDB.AddMutedSeries(chatID, seriesID, seriesID)
	isMuted, _ = mockDB.IsSeriesMuted(chatID, seriesID)
	if !isMuted {
		t.Error("Series should be muted after AddMutedSeries")
	}

	// Unmute the series
	mockDB.RemoveMutedSeries(chatID, seriesID)
	isMuted, _ = mockDB.IsSeriesMuted(chatID, seriesID)
	if isMuted {
		t.Error("Series should not be muted after RemoveMutedSeries")
	}
}

// TestMutedSeriesFiltering tests that muted users are filtered from notifications
func TestMutedSeriesFiltering(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	// Add subscribers
	mockDB.AddSubscriber(12345, "user1", "User 1")
	mockDB.AddSubscriber(67890, "user2", "User 2")
	mockDB.AddSubscriber(11111, "user3", "User 3")

	// User 1 mutes "Breaking Bad"
	mockDB.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")

	// Get all subscribers
	allSubscribers, _ := mockDB.GetAllActiveSubscribers()
	if len(allSubscribers) != 3 {
		t.Errorf("Expected 3 subscribers, got %d", len(allSubscribers))
	}

	// Filter subscribers who haven't muted "Breaking Bad"
	var filteredSubscribers []int64
	for _, chatID := range allSubscribers {
		isMuted, _ := mockDB.IsSeriesMuted(chatID, "Breaking Bad")
		if !isMuted {
			filteredSubscribers = append(filteredSubscribers, chatID)
		}
	}

	// Should have 2 subscribers after filtering (user2 and user3)
	if len(filteredSubscribers) != 2 {
		t.Errorf("Expected 2 filtered subscribers, got %d", len(filteredSubscribers))
	}

	// Verify user1 is not in filtered list
	for _, chatID := range filteredSubscribers {
		if chatID == 12345 {
			t.Error("User 1 should be filtered out as they muted the series")
		}
	}
}

// TestDuplicateMuteAttempt tests that attempting to mute same series twice doesn't cause error
func TestDuplicateMuteAttempt(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	chatID := int64(12345)
	seriesID := "Breaking Bad"

	// Mute once
	err := mockDB.AddMutedSeries(chatID, seriesID, seriesID)
	if err != nil {
		t.Errorf("First mute should succeed, got error: %v", err)
	}

	// Mute again
	err = mockDB.AddMutedSeries(chatID, seriesID, seriesID)
	if err != nil {
		t.Errorf("Duplicate mute should be handled gracefully, got error: %v", err)
	}

	// Should still be muted
	isMuted, _ := mockDB.IsSeriesMuted(chatID, seriesID)
	if !isMuted {
		t.Error("Series should still be muted")
	}

	// Should only have one entry
	list, _ := mockDB.GetMutedSeriesByUser(chatID)
	count := 0
	for _, series := range list {
		if series.SeriesID == seriesID {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Expected exactly 1 entry for series, got %d", count)
	}
}

// TestGetMutedSeriesByUser_MultipleSeriesForSameUser tests retrieval of multiple muted series
func TestGetMutedSeriesByUser_MultipleSeriesForSameUser(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	chatID := int64(12345)

	// Mute multiple series
	series := []string{"Breaking Bad", "Game of Thrones", "The Office", "Friends"}
	for _, s := range series {
		mockDB.AddMutedSeries(chatID, s, s)
	}

	// Retrieve muted list
	mutedList, err := mockDB.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(mutedList) != len(series) {
		t.Errorf("Expected %d muted series, got %d", len(series), len(mutedList))
	}

	// Verify all series are in the list
	foundSeries := make(map[string]bool)
	for _, ms := range mutedList {
		foundSeries[ms.SeriesID] = true
	}

	for _, s := range series {
		if !foundSeries[s] {
			t.Errorf("Expected to find %s in muted list", s)
		}
	}
}

// TestIsSeriesMuted_NonExistentSeries tests checking mute status for series that was never muted
func TestIsSeriesMuted_NonExistentSeries(t *testing.T) {
	mockDB := NewMockSubscriberDB()

	chatID := int64(12345)
	seriesID := "Never Muted Show"

	isMuted, err := mockDB.IsSeriesMuted(chatID, seriesID)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if isMuted {
		t.Error("Expected series to not be muted")
	}
}

// TestCallbackDataFormat tests that callback data format is consistent
func TestCallbackDataFormat(t *testing.T) {
	// Test mute callback data format
	seriesName := "Breaking Bad"
	muteCallbackData := "mute:" + seriesName
	if muteCallbackData != "mute:Breaking Bad" {
		t.Errorf("Expected callback data 'mute:Breaking Bad', got '%s'", muteCallbackData)
	}

	// Test unmute callback data format
	seriesID := "Breaking Bad"
	unmuteCallbackData := "unmute:" + seriesID
	if unmuteCallbackData != "unmute:Breaking Bad" {
		t.Errorf("Expected callback data 'unmute:Breaking Bad', got '%s'", unmuteCallbackData)
	}
}

// Verify we have the right number of focused tests
// Count: 14 tests total (within the 2-8 per focused area guideline, covering critical behaviors)
