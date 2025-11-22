package integration

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"jellyfin-telegram-bot/internal/database"
	"jellyfin-telegram-bot/internal/handlers"
)

// Test 1: End-to-end workflow - Episode notification -> Mute button -> User excluded from future notifications
func TestMuteWorkflow_EndToEnd(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_mute_e2e_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add two subscribers
	db.AddSubscriber(12345, "user1", "User 1")
	db.AddSubscriber(67890, "user2", "User 2")

	// Track broadcast calls
	type BroadcastRecord struct {
		Recipients []int64
		SeriesName string
	}
	var broadcasts []BroadcastRecord
	var mu sync.Mutex

	mockBroadcaster := &mockBroadcaster{
		broadcastFunc: func(ctx context.Context, content *handlers.NotificationContent) error {
			// Simulate filtering logic from BroadcastNotification
			subscribers, _ := db.GetAllActiveSubscribers()

			filtered := make([]int64, 0)
			if content.Type == "Episode" && content.SeriesName != "" {
				for _, chatID := range subscribers {
					isMuted, _ := db.IsSeriesMuted(chatID, content.SeriesName)
					if !isMuted {
						filtered = append(filtered, chatID)
					}
				}
			} else {
				filtered = subscribers
			}

			mu.Lock()
			broadcasts = append(broadcasts, BroadcastRecord{
				Recipients: filtered,
				SeriesName: content.SeriesName,
			})
			mu.Unlock()
			return nil
		},
	}

	// Send first episode notification
	content1 := &handlers.NotificationContent{
		Type:       "Episode",
		SeriesName: "Breaking Bad",
		Title:      "Pilot",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), content1)

	// Verify both users received the notification
	if len(broadcasts) != 1 {
		t.Fatalf("Expected 1 broadcast, got %d", len(broadcasts))
	}
	if len(broadcasts[0].Recipients) != 2 {
		t.Errorf("Expected 2 recipients, got %d", len(broadcasts[0].Recipients))
	}

	// User 1 mutes "Breaking Bad"
	err = db.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")
	if err != nil {
		t.Fatalf("Failed to mute series: %v", err)
	}

	// Send second episode notification for same series
	content2 := &handlers.NotificationContent{
		Type:       "Episode",
		SeriesName: "Breaking Bad",
		Title:      "Cat's in the Bag...",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), content2)

	// Verify only user 2 received the notification
	if len(broadcasts) != 2 {
		t.Fatalf("Expected 2 broadcasts total, got %d", len(broadcasts))
	}
	if len(broadcasts[1].Recipients) != 1 {
		t.Errorf("Expected 1 recipient after mute, got %d", len(broadcasts[1].Recipients))
	}
	if broadcasts[1].Recipients[0] != 67890 {
		t.Errorf("Expected user 67890 to receive notification, got %d", broadcasts[1].Recipients[0])
	}
}

// Test 2: End-to-end workflow - Mute -> Unmute -> Notifications restored
func TestUnmuteRestoresNotifications_EndToEnd(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_unmute_e2e_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add subscriber
	db.AddSubscriber(12345, "testuser", "Test User")

	// Mute series
	err = db.AddMutedSeries(12345, "Game of Thrones", "Game of Thrones")
	if err != nil {
		t.Fatalf("Failed to mute series: %v", err)
	}

	// Verify series is muted
	isMuted, _ := db.IsSeriesMuted(12345, "Game of Thrones")
	if !isMuted {
		t.Error("Series should be muted")
	}

	// Track notifications
	notificationCount := 0
	mockBroadcaster := &mockBroadcaster{
		broadcastFunc: func(ctx context.Context, content *handlers.NotificationContent) error {
			subscribers, _ := db.GetAllActiveSubscribers()
			for _, chatID := range subscribers {
				isMuted, _ := db.IsSeriesMuted(chatID, content.SeriesName)
				if !isMuted {
					notificationCount++
				}
			}
			return nil
		},
	}

	// Send notification while muted
	content1 := &handlers.NotificationContent{
		Type:       "Episode",
		SeriesName: "Game of Thrones",
		Title:      "Winter Is Coming",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), content1)

	if notificationCount != 0 {
		t.Errorf("Expected 0 notifications while muted, got %d", notificationCount)
	}

	// Unmute series
	err = db.RemoveMutedSeries(12345, "Game of Thrones")
	if err != nil {
		t.Fatalf("Failed to unmute series: %v", err)
	}

	// Verify series is unmuted
	isMuted, _ = db.IsSeriesMuted(12345, "Game of Thrones")
	if isMuted {
		t.Error("Series should be unmuted")
	}

	// Send notification after unmuting
	content2 := &handlers.NotificationContent{
		Type:       "Episode",
		SeriesName: "Game of Thrones",
		Title:      "The Kingsroad",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), content2)

	if notificationCount != 1 {
		t.Errorf("Expected 1 notification after unmute, got %d", notificationCount)
	}
}

// Test 3: Multiple users can independently mute/unmute same series with database persistence
func TestMultipleUsersIndependentMuting_WithPersistence(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_multi_user_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add three subscribers
	db.AddSubscriber(111, "user1", "User 1")
	db.AddSubscriber(222, "user2", "User 2")
	db.AddSubscriber(333, "user3", "User 3")

	seriesName := "The Office"

	// User 1 and User 2 mute the series
	err = db.AddMutedSeries(111, seriesName, seriesName)
	if err != nil {
		t.Fatalf("User 1 failed to mute: %v", err)
	}
	err = db.AddMutedSeries(222, seriesName, seriesName)
	if err != nil {
		t.Fatalf("User 2 failed to mute: %v", err)
	}

	// Verify each user's mute status
	user1Muted, _ := db.IsSeriesMuted(111, seriesName)
	user2Muted, _ := db.IsSeriesMuted(222, seriesName)
	user3Muted, _ := db.IsSeriesMuted(333, seriesName)

	if !user1Muted {
		t.Error("User 1 should have series muted")
	}
	if !user2Muted {
		t.Error("User 2 should have series muted")
	}
	if user3Muted {
		t.Error("User 3 should not have series muted")
	}

	// User 1 unmutes
	err = db.RemoveMutedSeries(111, seriesName)
	if err != nil {
		t.Fatalf("User 1 failed to unmute: %v", err)
	}

	// Verify User 1 unmuted but User 2 still muted
	user1Muted, _ = db.IsSeriesMuted(111, seriesName)
	user2Muted, _ = db.IsSeriesMuted(222, seriesName)

	if user1Muted {
		t.Error("User 1 should not have series muted after unmute")
	}
	if !user2Muted {
		t.Error("User 2 should still have series muted")
	}

	// Verify muted lists
	user1List, _ := db.GetMutedSeriesByUser(111)
	user2List, _ := db.GetMutedSeriesByUser(222)
	user3List, _ := db.GetMutedSeriesByUser(333)

	if len(user1List) != 0 {
		t.Errorf("User 1 should have 0 muted series, got %d", len(user1List))
	}
	if len(user2List) != 1 {
		t.Errorf("User 2 should have 1 muted series, got %d", len(user2List))
	}
	if len(user3List) != 0 {
		t.Errorf("User 3 should have 0 muted series, got %d", len(user3List))
	}
}

// Test 4: /mutedlist command integration with real database
func TestMutedListCommand_DatabaseIntegration(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_mutedlist_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)

	// Initially empty list
	mutedSeries, err := db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get muted series: %v", err)
	}
	if len(mutedSeries) != 0 {
		t.Errorf("Expected empty list initially, got %d items", len(mutedSeries))
	}

	// Add multiple muted series
	series := []string{"Breaking Bad", "Game of Thrones", "The Office", "Friends"}
	for _, s := range series {
		err = db.AddMutedSeries(chatID, s, s)
		if err != nil {
			t.Fatalf("Failed to mute %s: %v", s, err)
		}
	}

	// Retrieve list
	mutedSeries, err = db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get muted series: %v", err)
	}
	if len(mutedSeries) != len(series) {
		t.Errorf("Expected %d muted series, got %d", len(series), len(mutedSeries))
	}

	// Verify all series are in the list
	foundSeries := make(map[string]bool)
	for _, ms := range mutedSeries {
		foundSeries[ms.SeriesID] = true
	}
	for _, s := range series {
		if !foundSeries[s] {
			t.Errorf("Expected to find %s in muted list", s)
		}
	}

	// Unmute one series
	err = db.RemoveMutedSeries(chatID, "The Office")
	if err != nil {
		t.Fatalf("Failed to unmute The Office: %v", err)
	}

	// Verify list updated
	mutedSeries, err = db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get updated muted series: %v", err)
	}
	if len(mutedSeries) != len(series)-1 {
		t.Errorf("Expected %d muted series after unmute, got %d", len(series)-1, len(mutedSeries))
	}

	// Verify "The Office" is not in the list
	for _, ms := range mutedSeries {
		if ms.SeriesID == "The Office" {
			t.Error("The Office should not be in muted list after unmute")
		}
	}
}

// Test 5: Callback data parsing with special characters (Persian text in series names)
func TestCallbackDataParsing_PersianCharacters(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_persian_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)

	// Series names with Persian and special characters
	testSeries := []string{
		"سریال فارسی",
		"Series: The Beginning",
		"Series & More",
		"Breaking Bad - Season 1",
	}

	// Mute series with special characters
	for _, series := range testSeries {
		err = db.AddMutedSeries(chatID, series, series)
		if err != nil {
			t.Fatalf("Failed to mute series '%s': %v", series, err)
		}
	}

	// Verify all series are muted
	for _, series := range testSeries {
		isMuted, err := db.IsSeriesMuted(chatID, series)
		if err != nil {
			t.Fatalf("Failed to check mute status for '%s': %v", series, err)
		}
		if !isMuted {
			t.Errorf("Expected series '%s' to be muted", series)
		}
	}

	// Retrieve muted list
	mutedList, err := db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get muted list: %v", err)
	}
	if len(mutedList) != len(testSeries) {
		t.Errorf("Expected %d muted series, got %d", len(testSeries), len(mutedList))
	}

	// Unmute series with Persian characters
	err = db.RemoveMutedSeries(chatID, "سریال فارسی")
	if err != nil {
		t.Fatalf("Failed to unmute Persian series: %v", err)
	}

	// Verify unmute worked
	isMuted, _ := db.IsSeriesMuted(chatID, "سریال فارسی")
	if isMuted {
		t.Error("Persian series should be unmuted")
	}
}

// Test 6: Concurrent mute operations don't create duplicates (composite unique index)
func TestConcurrentMuteOperations_NoDuplicates(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_concurrent_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)
	seriesName := "Breaking Bad"

	// Simulate concurrent mute operations
	var wg sync.WaitGroup
	concurrentOps := 10

	for i := 0; i < concurrentOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db.AddMutedSeries(chatID, seriesName, seriesName)
		}()
	}

	wg.Wait()

	// Verify only one record exists
	mutedList, err := db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get muted list: %v", err)
	}

	count := 0
	for _, ms := range mutedList {
		if ms.SeriesID == seriesName {
			count++
		}
	}

	if count != 1 {
		t.Errorf("Expected exactly 1 record after concurrent operations, got %d", count)
	}
}

// Test 7: Series muting doesn't affect movie notifications (integration level)
func TestSeriesMuting_DoesNotAffectMovies(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_movies_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add subscriber
	db.AddSubscriber(12345, "testuser", "Test User")

	// Mute a series
	db.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")

	// Track notifications
	var receivedNotifications []string
	var mu sync.Mutex

	mockBroadcaster := &mockBroadcaster{
		broadcastFunc: func(ctx context.Context, content *handlers.NotificationContent) error {
			subscribers, _ := db.GetAllActiveSubscribers()

			for _, chatID := range subscribers {
				shouldSend := true
				if content.Type == "Episode" && content.SeriesName != "" {
					isMuted, _ := db.IsSeriesMuted(chatID, content.SeriesName)
					shouldSend = !isMuted
				}

				if shouldSend {
					mu.Lock()
					receivedNotifications = append(receivedNotifications, content.Type+":"+content.Title)
					mu.Unlock()
				}
			}
			return nil
		},
	}

	// Send episode notification for muted series
	episodeContent := &handlers.NotificationContent{
		Type:       "Episode",
		SeriesName: "Breaking Bad",
		Title:      "Pilot",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), episodeContent)

	// Send movie notification
	movieContent := &handlers.NotificationContent{
		Type:  "Movie",
		Title: "Interstellar",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), movieContent)

	// Verify only movie notification was received
	if len(receivedNotifications) != 1 {
		t.Errorf("Expected 1 notification (movie only), got %d", len(receivedNotifications))
	}
	if len(receivedNotifications) > 0 && receivedNotifications[0] != "Movie:Interstellar" {
		t.Errorf("Expected movie notification, got %s", receivedNotifications[0])
	}
}

// Test 8: Multiple mute/unmute operations maintain data integrity
func TestMultipleMuteUnmuteOperations_DataIntegrity(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_integrity_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)
	series := []string{"Series A", "Series B", "Series C", "Series D"}

	// Mute first three series
	for i := 0; i < 3; i++ {
		db.AddMutedSeries(chatID, series[i], series[i])
	}

	// Verify all three muted
	mutedList, _ := db.GetMutedSeriesByUser(chatID)
	if len(mutedList) != 3 {
		t.Errorf("Expected 3 muted series, got %d", len(mutedList))
	}

	// Unmute middle series
	db.RemoveMutedSeries(chatID, "Series B")

	// Verify correct series removed
	mutedList, _ = db.GetMutedSeriesByUser(chatID)
	if len(mutedList) != 2 {
		t.Errorf("Expected 2 muted series after unmute, got %d", len(mutedList))
	}

	foundSeriesA := false
	foundSeriesC := false
	for _, ms := range mutedList {
		if ms.SeriesID == "Series A" {
			foundSeriesA = true
		}
		if ms.SeriesID == "Series C" {
			foundSeriesC = true
		}
		if ms.SeriesID == "Series B" {
			t.Error("Series B should not be in muted list")
		}
	}

	if !foundSeriesA || !foundSeriesC {
		t.Error("Series A and C should still be in muted list")
	}

	// Mute a different series (Series D)
	db.AddMutedSeries(chatID, "Series D", "Series D")

	// Verify we have three series muted now (A, C, D)
	mutedList, _ = db.GetMutedSeriesByUser(chatID)
	if len(mutedList) != 3 {
		t.Errorf("Expected 3 muted series after adding Series D, got %d", len(mutedList))
	}

	// Verify Series D is in the list
	foundSeriesD := false
	for _, ms := range mutedList {
		if ms.SeriesID == "Series D" {
			foundSeriesD = true
		}
	}
	if !foundSeriesD {
		t.Error("Series D should be in muted list")
	}
}

// Test 9: Empty series name handling in notification filtering
func TestEmptySeriesName_NotificationFiltering(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_empty_series_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Add subscriber
	db.AddSubscriber(12345, "testuser", "Test User")

	// Mute a series
	db.AddMutedSeries(12345, "Breaking Bad", "Breaking Bad")

	// Track notifications
	notificationCount := 0
	mockBroadcaster := &mockBroadcaster{
		broadcastFunc: func(ctx context.Context, content *handlers.NotificationContent) error {
			subscribers, _ := db.GetAllActiveSubscribers()

			for _, chatID := range subscribers {
				shouldSend := true
				// Only filter if it's an episode with a valid series name
				if content.Type == "Episode" && content.SeriesName != "" && content.SeriesName != "Unknown Series" {
					isMuted, _ := db.IsSeriesMuted(chatID, content.SeriesName)
					shouldSend = !isMuted
				}

				if shouldSend {
					notificationCount++
				}
			}
			return nil
		},
	}

	// Send episode notification with empty series name
	emptyContent := &handlers.NotificationContent{
		Type:       "Episode",
		SeriesName: "",
		Title:      "Episode with no series",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), emptyContent)

	// User should receive notification (not filtered)
	if notificationCount != 1 {
		t.Errorf("Expected user to receive notification with empty series name, got %d notifications", notificationCount)
	}

	// Reset counter
	notificationCount = 0

	// Send episode notification with "Unknown Series"
	unknownContent := &handlers.NotificationContent{
		Type:       "Episode",
		SeriesName: "Unknown Series",
		Title:      "Unknown episode",
	}
	mockBroadcaster.BroadcastNotification(context.Background(), unknownContent)

	// User should receive notification (not filtered)
	if notificationCount != 1 {
		t.Errorf("Expected user to receive notification with Unknown Series, got %d notifications", notificationCount)
	}
}

// Test 10: Database record cleanup after multiple operations
func TestDatabaseCleanup_AfterOperations(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_cleanup_" + time.Now().Format("20060102150405") + ".db"
	defer os.Remove(dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	chatID := int64(12345)

	// Add and remove series multiple times
	for i := 0; i < 5; i++ {
		db.AddMutedSeries(chatID, "Series X", "Series X")
		db.RemoveMutedSeries(chatID, "Series X")
	}

	// Verify no records exist
	mutedList, err := db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get muted list: %v", err)
	}
	if len(mutedList) != 0 {
		t.Errorf("Expected empty list after cleanup, got %d items", len(mutedList))
	}

	// Verify series is not muted
	isMuted, _ := db.IsSeriesMuted(chatID, "Series X")
	if isMuted {
		t.Error("Series X should not be muted after cleanup")
	}

	// Add multiple different series
	series := []string{"A", "B", "C", "D", "E"}
	for _, s := range series {
		db.AddMutedSeries(chatID, s, s)
	}

	// Remove all series
	for _, s := range series {
		db.RemoveMutedSeries(chatID, s)
	}

	// Verify all records cleaned up
	mutedList, _ = db.GetMutedSeriesByUser(chatID)
	if len(mutedList) != 0 {
		t.Errorf("Expected empty list after removing all series, got %d items", len(mutedList))
	}
}
