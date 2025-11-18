package database

import (
	"testing"
)

// Test 7: Mark content as notified and check status
func TestMarkContentNotified(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	jellyfinID := "test-movie-123"
	title := "Test Movie"
	contentType := "Movie"

	// Verify content not notified initially
	isNotified, err := db.IsContentNotified(jellyfinID)
	if err != nil {
		t.Fatalf("Failed to check content notification status: %v", err)
	}
	if isNotified {
		t.Error("Expected content to not be notified initially")
	}

	// Mark content as notified
	err = db.MarkContentNotified(jellyfinID, title, contentType)
	if err != nil {
		t.Fatalf("Failed to mark content as notified: %v", err)
	}

	// Verify content is now notified
	isNotified, err = db.IsContentNotified(jellyfinID)
	if err != nil {
		t.Fatalf("Failed to check content notification status: %v", err)
	}
	if !isNotified {
		t.Error("Expected content to be marked as notified")
	}
}

// Test 8: Prevent duplicate content notifications
func TestPreventDuplicateNotifications(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	jellyfinID := "duplicate-test-456"

	// Mark content as notified
	err := db.MarkContentNotified(jellyfinID, "Duplicate Test", "Episode")
	if err != nil {
		t.Fatalf("Failed to mark content as notified: %v", err)
	}

	// Check if content is notified
	isNotified, err := db.IsContentNotified(jellyfinID)
	if err != nil {
		t.Fatalf("Failed to check content notification status: %v", err)
	}
	if !isNotified {
		t.Fatal("Content should be marked as notified")
	}

	// Try to mark same content again (should fail due to unique constraint)
	err = db.MarkContentNotified(jellyfinID, "Duplicate Test", "Episode")
	if err == nil {
		t.Error("Expected error when marking duplicate content, got nil")
	}
}
