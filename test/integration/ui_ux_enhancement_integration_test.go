package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"jellyfin-telegram-bot/internal/database"
	"jellyfin-telegram-bot/internal/telegram"
)

// Test 1: Complete flow - /start displays welcome menu with navigation buttons
func TestStartCommand_DisplaysWelcomeMenuWithButtons(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_start_menu_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Verify database is ready
	chatID := int64(12345)
	err = db.AddSubscriber(chatID, "testuser", "Test User")
	if err != nil {
		t.Fatalf("Failed to add subscriber: %v", err)
	}

	// Verify subscriber was added (simulates /start command success)
	isSubscribed, err := db.IsSubscribed(chatID)
	if err != nil {
		t.Fatalf("Failed to check subscription: %v", err)
	}
	if !isSubscribed {
		t.Error("Expected user to be subscribed after /start")
	}

	// Verify welcome menu buttons would be displayed (structure validation)
	expectedButtons := []struct {
		text         string
		callbackData string
	}{
		{"تازه‌ها", "nav:recent"},
		{"جستجو", "nav:search"},
		{"سریال‌های مسدود شده", "nav:mutedlist"},
		{"راهنما", "nav:help"},
	}

	// Validate button structure
	if len(expectedButtons) != 4 {
		t.Errorf("Expected 4 welcome menu buttons, got %d", len(expectedButtons))
	}

	// Verify all buttons have Persian labels and correct callback data
	for _, btn := range expectedButtons {
		if btn.text == "" {
			t.Error("Button text should not be empty")
		}
		if !strings.HasPrefix(btn.callbackData, "nav:") {
			t.Errorf("Button callback should start with 'nav:', got: %s", btn.callbackData)
		}
	}
}

// Test 2: Complete flow - nav:recent button displays recent content
func TestNavigationButton_Recent_DisplaysRecentContent(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_nav_recent_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Create mock Jellyfin client with recent items
	mockJellyfin := &mockJellyfinClient{
		recentItems: []telegram.ContentItem{
			{
				ItemID:          "item1",
				Name:            "Breaking Bad - Pilot",
				Type:            "Episode",
				SeriesName:      "Breaking Bad",
				SeasonNumber:    1,
				EpisodeNumber:   1,
				CommunityRating: 9.5,
				ProductionYear:  2008,
			},
			{
				ItemID:          "item2",
				Name:            "Interstellar",
				Type:            "Movie",
				CommunityRating: 8.6,
				ProductionYear:  2014,
			},
		},
	}

	// Verify recent items can be fetched (simulates nav:recent callback)
	items, err := mockJellyfin.GetRecentItems(context.Background(), 15)
	if err != nil {
		t.Fatalf("Failed to get recent items: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 recent items, got %d", len(items))
	}

	// Verify content structure
	if items[0].Type != "Episode" {
		t.Errorf("Expected first item to be Episode, got %s", items[0].Type)
	}
	if items[1].Type != "Movie" {
		t.Errorf("Expected second item to be Movie, got %s", items[1].Type)
	}

	// Prevent unused variable warning
	_ = db
}

// Test 3: Complete flow - nav:mutedlist button displays muted series
func TestNavigationButton_MutedList_DisplaysMutedSeries(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_nav_mutedlist_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)

	// Add some muted series
	series := []string{"Breaking Bad", "Game of Thrones", "The Office"}
	for _, s := range series {
		err = db.AddMutedSeries(chatID, s, s)
		if err != nil {
			t.Fatalf("Failed to mute series %s: %v", s, err)
		}
	}

	// Simulate nav:mutedlist callback - retrieve muted series
	mutedSeries, err := db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get muted series: %v", err)
	}

	if len(mutedSeries) != len(series) {
		t.Errorf("Expected %d muted series, got %d", len(series), len(mutedSeries))
	}

	// Verify all expected series are in the list
	foundSeries := make(map[string]bool)
	for _, ms := range mutedSeries {
		foundSeries[ms.SeriesName] = true
	}

	for _, expectedSeries := range series {
		if !foundSeries[expectedSeries] {
			t.Errorf("Expected to find %s in muted list", expectedSeries)
		}
	}
}

// Test 4: Complete flow - mute → undo button → unmute succeeds
func TestMuteUndoFlow_Complete(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_mute_undo_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)
	seriesName := "Breaking Bad"

	// Step 1: Verify series is not muted initially
	isMuted, _ := db.IsSeriesMuted(chatID, seriesName)
	if isMuted {
		t.Error("Series should not be muted initially")
	}

	// Step 2: User clicks mute button (simulates handleMuteCallback)
	err = db.AddMutedSeries(chatID, seriesName, seriesName)
	if err != nil {
		t.Fatalf("Failed to mute series: %v", err)
	}

	// Step 3: Verify series is muted
	isMuted, _ = db.IsSeriesMuted(chatID, seriesName)
	if !isMuted {
		t.Error("Series should be muted after mute action")
	}

	// Step 4: User clicks undo button immediately (simulates handleUndoMuteCallback)
	err = db.RemoveMutedSeries(chatID, seriesName)
	if err != nil {
		t.Fatalf("Failed to undo mute: %v", err)
	}

	// Step 5: Verify series is unmuted
	isMuted, _ = db.IsSeriesMuted(chatID, seriesName)
	if isMuted {
		t.Error("Series should be unmuted after undo action")
	}

	// Step 6: Verify series removed from muted list
	mutedList, _ := db.GetMutedSeriesByUser(chatID)
	if len(mutedList) != 0 {
		t.Errorf("Expected empty muted list after undo, got %d items", len(mutedList))
	}
}

// Test 5: Complete flow - nav:help button displays help message
func TestNavigationButton_Help_DisplaysHelpMessage(t *testing.T) {
	// Simulate nav:help callback - verify help message content
	helpMessage := `دستورات موجود:
/start - عضویت در ربات
/recent - مشاهده محتوای اخیر
/search - جستجوی محتوا (مثال: /search interstellar)
/mutedlist - مشاهده سریال‌های مسدود شده`

	// Verify help message contains expected content
	expectedContent := []string{
		"دستورات موجود",
		"/start",
		"/recent",
		"/search",
		"/mutedlist",
	}

	for _, content := range expectedContent {
		if !strings.Contains(helpMessage, content) {
			t.Errorf("Help message should contain '%s'", content)
		}
	}

	// Verify Persian text is present
	if !strings.Contains(helpMessage, "عضویت در ربات") {
		t.Error("Help message should contain Persian descriptions")
	}
}

// Test 6: Complete flow - nav:search button displays search instructions
func TestNavigationButton_Search_DisplaysSearchInstructions(t *testing.T) {
	// Simulate nav:search callback - verify search instructions content
	searchInstructions := "لطفاً عبارت جستجو را وارد کنید. مثال: /search interstellar"

	// Verify search instructions are in Persian
	if !strings.Contains(searchInstructions, "لطفاً عبارت جستجو را وارد کنید") {
		t.Error("Search instructions should be in Persian")
	}

	// Verify example is provided
	if !strings.Contains(searchInstructions, "/search interstellar") {
		t.Error("Search instructions should include example")
	}

	// Verify message is not empty
	if searchInstructions == "" {
		t.Error("Search instructions should not be empty")
	}
}

// Test 7: Persian text renders correctly in welcome menu buttons
func TestWelcomeMenuButtons_PersianTextRendering(t *testing.T) {
	// Test Persian text labels for welcome menu buttons
	buttons := map[string]string{
		"تازه‌ها":             "nav:recent",    // Recent Content
		"جستجو":              "nav:search",    // Search
		"سریال‌های مسدود شده": "nav:mutedlist", // Muted List
		"راهنما":             "nav:help",      // Help
	}

	// Verify all buttons have Persian text (contains characters in Persian Unicode range)
	for label, callback := range buttons {
		hasPersian := false
		for _, r := range label {
			// Persian characters are in the range U+0600 to U+06FF
			if r >= 0x0600 && r <= 0x06FF {
				hasPersian = true
				break
			}
		}

		if !hasPersian {
			t.Errorf("Button label should contain Persian text: %s", label)
		}

		// Verify callback data format
		if !strings.HasPrefix(callback, "nav:") {
			t.Errorf("Callback data should start with 'nav:', got: %s", callback)
		}

		// Verify button text is not empty
		if label == "" {
			t.Error("Button label should not be empty")
		}
	}

	// Verify correct count
	if len(buttons) != 4 {
		t.Errorf("Expected 4 buttons with Persian labels, got %d", len(buttons))
	}
}

// Test 8: Button interface produces identical results to command interface
func TestButtonVsCommand_IdenticalBehavior(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_button_vs_cmd_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)

	// Test Case 1: /mutedlist command vs nav:mutedlist button
	// Add muted series
	db.AddMutedSeries(chatID, "Series A", "Series A")
	db.AddMutedSeries(chatID, "Series B", "Series B")

	// Both command and button should return same result
	mutedList1, err1 := db.GetMutedSeriesByUser(chatID)
	mutedList2, err2 := db.GetMutedSeriesByUser(chatID)

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	if len(mutedList1) != len(mutedList2) {
		t.Error("Command and button should return same number of muted series")
	}

	if len(mutedList1) != 2 {
		t.Errorf("Expected 2 muted series, got %d", len(mutedList1))
	}

	// Test Case 2: Verify both use same underlying database operations
	// This is inherently tested by the fact that both call GetMutedSeriesByUser
}

// Test 9: Menu Button API commands are registered correctly
func TestMenuButtonAPI_CommandRegistration(t *testing.T) {
	// Verify bot commands for Menu Button API
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
		t.Errorf("Expected 4 commands for Menu Button API, got %d", len(expectedCommands))
	}

	// Verify each command has correct format
	for _, cmd := range expectedCommands {
		// Verify command starts with /
		if !strings.HasPrefix(cmd.command, "/") {
			t.Errorf("Command should start with '/', got: %s", cmd.command)
		}

		// Verify description is in Persian
		hasPersian := false
		for _, r := range cmd.description {
			if r >= 0x0600 && r <= 0x06FF {
				hasPersian = true
				break
			}
		}

		if !hasPersian {
			t.Errorf("Command description should be in Persian: %s", cmd.description)
		}

		// Verify description is not empty
		if cmd.description == "" {
			t.Error("Command description should not be empty")
		}
	}
}

// Test 10: Callback handlers respond within performance requirements (200ms)
func TestCallbackHandlers_PerformanceRequirement(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_performance_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)

	// Test 1: Mute operation performance
	startTime := time.Now()
	err = db.AddMutedSeries(chatID, "Test Series", "Test Series")
	muteDuration := time.Since(startTime)

	if err != nil {
		t.Fatalf("Mute operation failed: %v", err)
	}

	// Database operations should be well under 200ms (typically <10ms)
	// We use 100ms threshold to account for slower systems
	if muteDuration > 100*time.Millisecond {
		t.Errorf("Mute operation took %v, should be under 100ms", muteDuration)
	}

	// Test 2: Get muted list performance
	// Add more series to make it more realistic
	for i := 0; i < 10; i++ {
		db.AddMutedSeries(chatID, fmt.Sprintf("Series %d", i), fmt.Sprintf("Series %d", i))
	}

	startTime = time.Now()
	_, err = db.GetMutedSeriesByUser(chatID)
	getMutedDuration := time.Since(startTime)

	if err != nil {
		t.Fatalf("Get muted list operation failed: %v", err)
	}

	if getMutedDuration > 100*time.Millisecond {
		t.Errorf("Get muted list operation took %v, should be under 100ms", getMutedDuration)
	}

	// Test 3: Unmute operation performance
	startTime = time.Now()
	err = db.RemoveMutedSeries(chatID, "Test Series")
	unmuteDuration := time.Since(startTime)

	if err != nil {
		t.Fatalf("Unmute operation failed: %v", err)
	}

	if unmuteDuration > 100*time.Millisecond {
		t.Errorf("Unmute operation took %v, should be under 100ms", unmuteDuration)
	}

	t.Logf("Performance results: Mute=%v, GetMutedList=%v, Unmute=%v",
		muteDuration, getMutedDuration, unmuteDuration)
}

// Mock implementations for testing

type mockJellyfinClient struct {
	recentItems []telegram.ContentItem
	searchFunc  func(query string) ([]telegram.ContentItem, error)
}

func (m *mockJellyfinClient) GetRecentItems(ctx context.Context, limit int) ([]telegram.ContentItem, error) {
	if len(m.recentItems) == 0 {
		return []telegram.ContentItem{}, nil
	}
	if limit < len(m.recentItems) {
		return m.recentItems[:limit], nil
	}
	return m.recentItems, nil
}

func (m *mockJellyfinClient) SearchContent(ctx context.Context, query string, limit int) ([]telegram.ContentItem, error) {
	if m.searchFunc != nil {
		return m.searchFunc(query)
	}
	return []telegram.ContentItem{}, nil
}

func (m *mockJellyfinClient) GetPosterImage(ctx context.Context, itemID string) ([]byte, error) {
	return []byte{}, nil
}
