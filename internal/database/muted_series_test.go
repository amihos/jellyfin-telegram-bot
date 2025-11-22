package database

import (
	"testing"

	"gorm.io/gorm"
)

// Test 1: Add muted series successfully
func TestAddMutedSeries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID := int64(123456789)
	seriesID := "Breaking Bad"
	seriesName := "Breaking Bad"

	err := db.AddMutedSeries(chatID, seriesID, seriesName)
	if err != nil {
		t.Fatalf("Failed to add muted series: %v", err)
	}

	// Verify series was muted
	isMuted, err := db.IsSeriesMuted(chatID, seriesID)
	if err != nil {
		t.Fatalf("Failed to check if series is muted: %v", err)
	}
	if !isMuted {
		t.Error("Expected series to be muted")
	}
}

// Test 2: Composite unique constraint prevents duplicates
func TestAddMutedSeriesDuplicate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID := int64(987654321)
	seriesID := "The Wire"
	seriesName := "The Wire"

	// Add series first time
	err := db.AddMutedSeries(chatID, seriesID, seriesName)
	if err != nil {
		t.Fatalf("Failed to add muted series first time: %v", err)
	}

	// Add same series again - should be handled gracefully
	err = db.AddMutedSeries(chatID, seriesID, seriesName)
	if err != nil {
		t.Fatalf("Failed to handle duplicate muted series: %v", err)
	}

	// Verify only one record exists
	mutedSeries, err := db.GetMutedSeriesByUser(chatID)
	if err != nil {
		t.Fatalf("Failed to get muted series: %v", err)
	}
	if len(mutedSeries) != 1 {
		t.Errorf("Expected 1 muted series, got %d", len(mutedSeries))
	}
}

// Test 3: Different users can mute the same series independently
func TestAddMutedSeriesDifferentUsers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	seriesID := "Game of Thrones"
	seriesName := "Game of Thrones"

	// Two different users mute the same series
	err := db.AddMutedSeries(111, seriesID, seriesName)
	if err != nil {
		t.Fatalf("Failed to add muted series for user 1: %v", err)
	}

	err = db.AddMutedSeries(222, seriesID, seriesName)
	if err != nil {
		t.Fatalf("Failed to add muted series for user 2: %v", err)
	}

	// Verify both users have the series muted
	isMuted1, _ := db.IsSeriesMuted(111, seriesID)
	isMuted2, _ := db.IsSeriesMuted(222, seriesID)

	if !isMuted1 || !isMuted2 {
		t.Error("Expected both users to have series muted")
	}
}

// Test 4: Remove muted series successfully
func TestRemoveMutedSeries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID := int64(333444555)
	seriesID := "Stranger Things"
	seriesName := "Stranger Things"

	// Add muted series
	err := db.AddMutedSeries(chatID, seriesID, seriesName)
	if err != nil {
		t.Fatalf("Failed to add muted series: %v", err)
	}

	// Remove muted series
	err = db.RemoveMutedSeries(chatID, seriesID)
	if err != nil {
		t.Fatalf("Failed to remove muted series: %v", err)
	}

	// Verify series is not muted
	isMuted, err := db.IsSeriesMuted(chatID, seriesID)
	if err != nil {
		t.Fatalf("Failed to check if series is muted: %v", err)
	}
	if isMuted {
		t.Error("Expected series to be unmuted")
	}
}

// Test 5: Remove non-existent muted series returns error
func TestRemoveNonExistentMutedSeries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := db.RemoveMutedSeries(999999999, "NonExistentSeries")
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

// Test 6: Get muted series by user returns correct filtered list
func TestGetMutedSeriesByUser(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID1 := int64(111)
	chatID2 := int64(222)

	// User 1 mutes multiple series
	db.AddMutedSeries(chatID1, "Breaking Bad", "Breaking Bad")
	db.AddMutedSeries(chatID1, "The Wire", "The Wire")
	db.AddMutedSeries(chatID1, "The Sopranos", "The Sopranos")

	// User 2 mutes one series
	db.AddMutedSeries(chatID2, "Friends", "Friends")

	// Get muted series for user 1
	mutedSeries1, err := db.GetMutedSeriesByUser(chatID1)
	if err != nil {
		t.Fatalf("Failed to get muted series for user 1: %v", err)
	}
	if len(mutedSeries1) != 3 {
		t.Errorf("Expected 3 muted series for user 1, got %d", len(mutedSeries1))
	}

	// Get muted series for user 2
	mutedSeries2, err := db.GetMutedSeriesByUser(chatID2)
	if err != nil {
		t.Fatalf("Failed to get muted series for user 2: %v", err)
	}
	if len(mutedSeries2) != 1 {
		t.Errorf("Expected 1 muted series for user 2, got %d", len(mutedSeries2))
	}
}

// Test 7: Get muted series for user with no muted series returns empty slice
func TestGetMutedSeriesByUserEmpty(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mutedSeries, err := db.GetMutedSeriesByUser(999999999)
	if err != nil {
		t.Fatalf("Failed to get muted series: %v", err)
	}
	if len(mutedSeries) != 0 {
		t.Errorf("Expected 0 muted series, got %d", len(mutedSeries))
	}
}

// Test 8: IsSeriesMuted returns correct boolean for muted/unmuted state
func TestIsSeriesMuted(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID := int64(666777888)
	mutedSeriesID := "The Office"
	unmutedSeriesID := "Parks and Recreation"

	// Mute one series
	db.AddMutedSeries(chatID, mutedSeriesID, "The Office")

	// Check muted series
	isMuted, err := db.IsSeriesMuted(chatID, mutedSeriesID)
	if err != nil {
		t.Fatalf("Failed to check muted series: %v", err)
	}
	if !isMuted {
		t.Error("Expected series to be muted")
	}

	// Check unmuted series
	isUnmuted, err := db.IsSeriesMuted(chatID, unmutedSeriesID)
	if err != nil {
		t.Fatalf("Failed to check unmuted series: %v", err)
	}
	if isUnmuted {
		t.Error("Expected series to be unmuted")
	}
}
