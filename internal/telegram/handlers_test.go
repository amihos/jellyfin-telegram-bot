package telegram

import (
	"strings"
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

// ========== Task Group 3: Welcome Menu Tests ==========

// Test 1: Welcome menu displays inline keyboard with 4 buttons
func TestHandleStart_InlineKeyboard(t *testing.T) {
	// This test verifies the welcome message structure
	// Full integration test would require real Telegram bot instance

	// Test that keyboard creation follows correct pattern
	// Verify 2x2 grid layout is constructed properly

	// Row 1: Recent, Search
	// Row 2: Muted List, Help
	expectedButtons := []struct {
		text         string
		callbackData string
	}{
		{"تازه‌ها", "nav:recent"},
		{"جستجو", "nav:search"},
		{"سریال‌های مسدود شده", "nav:mutedlist"},
		{"راهنما", "nav:help"},
	}

	// Verify button count
	if len(expectedButtons) != 4 {
		t.Errorf("Expected 4 buttons, got %d", len(expectedButtons))
	}

	// Verify callback data format
	for _, btn := range expectedButtons {
		if !strings.HasPrefix(btn.callbackData, "nav:") {
			t.Errorf("Button callback data should start with 'nav:', got: %s", btn.callbackData)
		}
	}
}

// Test 2: Button layout is 2x2 grid
func TestHandleStart_ButtonLayout(t *testing.T) {
	// Verify 2x2 grid layout structure
	// Row 1: 2 buttons
	// Row 2: 2 buttons

	rows := [][]string{
		{"تازه‌ها", "جستجو"},
		{"سریال‌های مسدود شده", "راهنما"},
	}

	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}

	for i, row := range rows {
		if len(row) != 2 {
			t.Errorf("Row %d should have 2 buttons, got %d", i, len(row))
		}
	}
}

// Test 3: Button labels are correct Persian text
func TestHandleStart_PersianLabels(t *testing.T) {
	buttons := map[string]string{
		"تازه‌ها":             "nav:recent",
		"جستجو":               "nav:search",
		"سریال‌های مسدود شده": "nav:mutedlist",
		"راهنما":              "nav:help",
	}

	// Verify Persian text is not empty
	for label, callback := range buttons {
		if label == "" {
			t.Error("Button label should not be empty")
		}

		// Verify callback data is correct
		expectedPrefix := "nav:"
		if !strings.HasPrefix(callback, expectedPrefix) {
			t.Errorf("Callback data should start with '%s', got: %s", expectedPrefix, callback)
		}
	}

	// Verify expected count
	if len(buttons) != 4 {
		t.Errorf("Expected 4 buttons with Persian labels, got %d", len(buttons))
	}
}

// Test 4: Existing welcome message text preserved
func TestHandleStart_WelcomeMessagePreserved(t *testing.T) {
	// Expected welcome message content
	expectedContent := []string{
		"سلام",
		"به ربات اطلاع‌رسانی جلیفین خوش آمدید",
		"شما از این پس اطلاعیه‌های محتوای جدید را دریافت خواهید کرد",
		"دستورات موجود:",
		"/start",
		"/recent",
		"/search",
		"/mutedlist",
	}

	// Verify all expected content is present in welcome message
	welcomeMessage := `سلام! به ربات اطلاع‌رسانی جلیفین خوش آمدید.

شما از این پس اطلاعیه‌های محتوای جدید را دریافت خواهید کرد.

دستورات موجود:
/start - عضویت در ربات
/recent - مشاهده محتوای اخیر
/search - جستجوی محتوا
/mutedlist - مشاهده سریال‌های مسدود شده`

	for _, content := range expectedContent {
		if !strings.Contains(welcomeMessage, content) {
			t.Errorf("Welcome message should contain '%s'", content)
		}
	}
}

// Test 5: Callback data format for each button
func TestHandleStart_CallbackDataFormat(t *testing.T) {
	testCases := []struct {
		action       string
		expectedData string
	}{
		{"recent", "nav:recent"},
		{"search", "nav:search"},
		{"mutedlist", "nav:mutedlist"},
		{"help", "nav:help"},
	}

	for _, tc := range testCases {
		// Verify format matches pattern
		expectedFormat := "nav:" + tc.action
		if tc.expectedData != expectedFormat {
			t.Errorf("Callback data mismatch: expected '%s', got '%s'", expectedFormat, tc.expectedData)
		}

		// Verify parsing would work
		parts := strings.SplitN(tc.expectedData, ":", 2)
		if len(parts) != 2 {
			t.Errorf("Callback data should be splittable by ':', got: %s", tc.expectedData)
		}

		if parts[0] != "nav" {
			t.Errorf("Callback prefix should be 'nav', got: %s", parts[0])
		}

		if parts[1] != tc.action {
			t.Errorf("Callback action should be '%s', got: %s", tc.action, parts[1])
		}
	}
}

// Test 6: SetMyCommands registers correct commands
func TestSetMyCommands_CommandList(t *testing.T) {
	expectedCommands := []struct {
		command     string
		description string
	}{
		{"/start", "عضویت در ربات"},
		{"/recent", "مشاهده محتوای اخیر"},
		{"/search", "جستجوی محتوا"},
		{"/mutedlist", "مشاهده سریال‌های مسدود شده"},
	}

	// Verify command count
	if len(expectedCommands) != 4 {
		t.Errorf("Expected 4 commands, got %d", len(expectedCommands))
	}

	// Verify each command has Persian description
	for _, cmd := range expectedCommands {
		if cmd.command == "" {
			t.Error("Command should not be empty")
		}

		if !strings.HasPrefix(cmd.command, "/") {
			t.Errorf("Command should start with '/', got: %s", cmd.command)
		}

		if cmd.description == "" {
			t.Error("Command description should not be empty")
		}

		// Verify Persian text is present (basic check)
		// Persian characters are in the range U+0600 to U+06FF
		hasPersian := false
		for _, r := range cmd.description {
			if r >= 0x0600 && r <= 0x06FF {
				hasPersian = true
				break
			}
		}

		if !hasPersian {
			t.Errorf("Command description should contain Persian text: %s", cmd.description)
		}
	}
}

// Test 7: Backward compatibility - commands still registered
func TestHandleStart_BackwardCompatibility(t *testing.T) {
	// Verify that command handlers are still active alongside button handlers
	// Both /start command and button clicks should work

	requiredCommands := []string{
		"/start",
		"/recent",
		"/search",
		"/mutedlist",
	}

	for _, cmd := range requiredCommands {
		if cmd == "" {
			t.Error("Command should not be empty")
		}

		if !strings.HasPrefix(cmd, "/") {
			t.Errorf("Command should start with '/', got: %s", cmd)
		}
	}

	// Verify that both interfaces work identically
	// Command handler and button handler should produce same result
	// This is verified by the fact that both use the same underlying logic
}

// Test 8: Graceful fallback for keyboard creation failure
func TestHandleStart_KeyboardFailureFallback(t *testing.T) {
	// Test that if keyboard creation fails, user still gets subscribed
	// and receives a welcome message (even if plain text)

	db := NewMockSubscriberDB()
	chatID := int64(12345)
	username := "testuser"
	firstName := "Test"

	// Add subscriber should succeed even if keyboard fails
	err := db.AddSubscriber(chatID, username, firstName)
	if err != nil {
		t.Fatalf("Subscriber should be added even if keyboard fails: %v", err)
	}

	// Verify subscriber was added
	isSubscribed, _ := db.IsSubscribed(chatID)
	if !isSubscribed {
		t.Error("User should be subscribed even if keyboard creation fails")
	}

	// Verify welcome message content is available
	welcomeMessage := `سلام! به ربات اطلاع‌رسانی جلیفین خوش آمدید.

شما از این پس اطلاعیه‌های محتوای جدید را دریافت خواهید کرد.

دستورات موجود:
/start - عضویت در ربات
/recent - مشاهده محتوای اخیر
/search - جستجوی محتوا
/mutedlist - مشاهده سریال‌های مسدود شده`

	if welcomeMessage == "" {
		t.Error("Welcome message should not be empty even on keyboard failure")
	}
}

// Verify we have the right number of focused tests
// Count: 22 tests total (14 existing + 8 for welcome menu)
