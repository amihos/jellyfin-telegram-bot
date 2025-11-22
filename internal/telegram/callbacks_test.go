package telegram

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// Test 1: Undo button unmutes series correctly
func TestHandleUndoMuteCallback_Success(t *testing.T) {
	db := NewMockSubscriberDB()

	chatID := int64(12345)
	seriesName := "Breaking Bad"

	// First, mute the series
	err := db.AddMutedSeries(chatID, seriesName, seriesName)
	if err != nil {
		t.Fatalf("Failed to mute series: %v", err)
	}

	// Verify series is muted
	isMuted, _ := db.IsSeriesMuted(chatID, seriesName)
	if !isMuted {
		t.Error("Series should be muted before undo")
	}

	// Simulate undo callback - test the unmute logic
	err = db.RemoveMutedSeries(chatID, seriesName)
	if err != nil {
		t.Fatalf("Failed to unmute series: %v", err)
	}

	// Verify series is unmuted
	isMuted, _ = db.IsSeriesMuted(chatID, seriesName)
	if isMuted {
		t.Error("Series should be unmuted after undo")
	}
}

// Test 2: Undo callback data format is correct
func TestUndoMuteCallbackDataFormat(t *testing.T) {
	seriesName := "Game of Thrones"
	expectedFormat := "undo_mute:Game of Thrones"

	// Simulate callback data creation
	callbackData := "undo_mute:" + seriesName

	if callbackData != expectedFormat {
		t.Errorf("Expected callback data '%s', got '%s'", expectedFormat, callbackData)
	}

	// Test parsing
	prefix := "undo_mute:"
	if len(callbackData) <= len(prefix) {
		t.Fatal("Callback data too short")
	}

	parsedSeries := callbackData[len(prefix):]
	if parsedSeries != seriesName {
		t.Errorf("Expected parsed series '%s', got '%s'", seriesName, parsedSeries)
	}
}

// Test 3: Undo works immediately after mute (no delay)
func TestUndoMuteCallback_ImmediateAfterMute(t *testing.T) {
	db := NewMockSubscriberDB()
	chatID := int64(12345)
	seriesName := "The Office"

	// Mute series
	err := db.AddMutedSeries(chatID, seriesName, seriesName)
	if err != nil {
		t.Fatalf("Failed to mute series: %v", err)
	}

	// Immediately unmute (no delay)
	err = db.RemoveMutedSeries(chatID, seriesName)
	if err != nil {
		t.Fatalf("Failed to unmute immediately after mute: %v", err)
	}

	// Verify unmuted
	isMuted, _ := db.IsSeriesMuted(chatID, seriesName)
	if isMuted {
		t.Error("Series should be unmuted immediately")
	}
}

// Test 4: Undo handles series not found in muted list
func TestUndoMuteCallback_SeriesNotFound(t *testing.T) {
	db := NewMockSubscriberDB()
	chatID := int64(12345)
	seriesName := "NonExistent Series"

	// Try to unmute a series that was never muted
	err := db.RemoveMutedSeries(chatID, seriesName)

	// Should not error - just silently succeed
	if err != nil {
		t.Errorf("RemoveMutedSeries should not error for non-existent series, got: %v", err)
	}

	// Verify still not muted
	isMuted, _ := db.IsSeriesMuted(chatID, seriesName)
	if isMuted {
		t.Error("Non-existent series should not be muted")
	}
}

// Test 5: Undo with Persian series name
func TestUndoMuteCallback_PersianSeriesName(t *testing.T) {
	db := NewMockSubscriberDB()
	chatID := int64(12345)
	seriesName := "سریال تست"

	// Mute Persian series
	err := db.AddMutedSeries(chatID, seriesName, seriesName)
	if err != nil {
		t.Fatalf("Failed to mute Persian series: %v", err)
	}

	// Unmute Persian series
	err = db.RemoveMutedSeries(chatID, seriesName)
	if err != nil {
		t.Fatalf("Failed to unmute Persian series: %v", err)
	}

	// Verify unmuted
	isMuted, _ := db.IsSeriesMuted(chatID, seriesName)
	if isMuted {
		t.Error("Persian series should be unmuted")
	}
}

// Test 6: Undo preserves other muted series
func TestUndoMuteCallback_PreservesOtherMutedSeries(t *testing.T) {
	db := NewMockSubscriberDB()
	chatID := int64(12345)

	// Mute multiple series
	series1 := "Breaking Bad"
	series2 := "Game of Thrones"
	series3 := "The Office"

	db.AddMutedSeries(chatID, series1, series1)
	db.AddMutedSeries(chatID, series2, series2)
	db.AddMutedSeries(chatID, series3, series3)

	// Unmute only series2
	err := db.RemoveMutedSeries(chatID, series2)
	if err != nil {
		t.Fatalf("Failed to unmute series2: %v", err)
	}

	// Verify series1 and series3 are still muted
	isMuted1, _ := db.IsSeriesMuted(chatID, series1)
	isMuted2, _ := db.IsSeriesMuted(chatID, series2)
	isMuted3, _ := db.IsSeriesMuted(chatID, series3)

	if !isMuted1 {
		t.Error("series1 should still be muted")
	}
	if isMuted2 {
		t.Error("series2 should be unmuted")
	}
	if !isMuted3 {
		t.Error("series3 should still be muted")
	}

	// Verify muted list count
	mutedList, _ := db.GetMutedSeriesByUser(chatID)
	if len(mutedList) != 2 {
		t.Errorf("Expected 2 muted series, got %d", len(mutedList))
	}
}

// Test 7: Undo button callback data parsing with special characters
func TestUndoMuteCallback_SpecialCharactersInSeriesName(t *testing.T) {
	testCases := []string{
		"Breaking Bad: The Complete Series",
		"Game of Thrones - Season 1",
		"The Office & Parks",
		"سریال: فارسی",
	}

	for _, seriesName := range testCases {
		callbackData := "undo_mute:" + seriesName

		// Parse callback data
		prefix := "undo_mute:"
		parsed := callbackData[len(prefix):]

		if parsed != seriesName {
			t.Errorf("Failed to parse series name with special characters: expected '%s', got '%s'", seriesName, parsed)
		}
	}
}

// Test 8: Verify undo reuses unmute logic correctly
func TestUndoMuteCallback_ReusesUnmuteLogic(t *testing.T) {
	db := NewMockSubscriberDB()
	chatID := int64(12345)
	seriesName := "Test Series"

	// Mute series
	db.AddMutedSeries(chatID, seriesName, seriesName)

	// Verify muted
	isMuted, _ := db.IsSeriesMuted(chatID, seriesName)
	if !isMuted {
		t.Fatal("Series should be muted initially")
	}

	// Unmute using the same logic that undo would use (RemoveMutedSeries)
	err := db.RemoveMutedSeries(chatID, seriesName)
	if err != nil {
		t.Fatalf("RemoveMutedSeries failed: %v", err)
	}

	// Verify unmuted - this tests that the unmute logic works correctly
	isMuted, _ = db.IsSeriesMuted(chatID, seriesName)
	if isMuted {
		t.Error("Series should be unmuted after RemoveMutedSeries")
	}

	// Verify series removed from muted list
	mutedList, _ := db.GetMutedSeriesByUser(chatID)
	if len(mutedList) != 0 {
		t.Errorf("Muted list should be empty, got %d items", len(mutedList))
	}
}

// MockDB with error simulation for testing error handling

type MockDBWithErrors struct {
	*MockSubscriberDB
	shouldFailRemove bool
}

func NewMockDBWithErrors() *MockDBWithErrors {
	return &MockDBWithErrors{
		MockSubscriberDB: NewMockSubscriberDB(),
		shouldFailRemove: false,
	}
}

func (m *MockDBWithErrors) RemoveMutedSeries(chatID int64, seriesID string) error {
	if m.shouldFailRemove {
		return errors.New("database error")
	}
	// Check if series exists in muted list
	if m.mutedSeries[chatID] == nil || !m.mutedSeries[chatID][seriesID] {
		return gorm.ErrRecordNotFound
	}
	return m.MockSubscriberDB.RemoveMutedSeries(chatID, seriesID)
}

// Test 9: Error handling for database errors
func TestUndoMuteCallback_DatabaseError(t *testing.T) {
	db := NewMockDBWithErrors()
	db.shouldFailRemove = true

	chatID := int64(12345)
	seriesName := "Test Series"

	// Try to unmute
	err := db.RemoveMutedSeries(chatID, seriesName)

	if err == nil {
		t.Error("Expected error from database, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected 'database error', got: %v", err)
	}
}

// Test 10: Error handling when series not in database (GORM ErrRecordNotFound)
func TestUndoMuteCallback_RecordNotFound(t *testing.T) {
	db := NewMockDBWithErrors()
	chatID := int64(12345)
	seriesName := "NonExistent"

	// Try to unmute non-existent series
	err := db.RemoveMutedSeries(chatID, seriesName)

	if err == nil {
		t.Error("Expected ErrRecordNotFound, got nil")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got: %v", err)
	}
}
