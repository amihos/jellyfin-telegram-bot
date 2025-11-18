package database

import (
	"os"
	"testing"
)

// TestDatabasePersistence verifies database persists across connections
func TestDatabasePersistence(t *testing.T) {
	tmpDB := "/tmp/test_persistence.db"
	defer os.Remove(tmpDB)

	// First connection: Add data
	db1, err := NewDB(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create first database connection: %v", err)
	}

	err = db1.AddSubscriber(999888777, "persisttest", "Persist Test")
	if err != nil {
		t.Fatalf("Failed to add subscriber: %v", err)
	}

	err = db1.MarkContentNotified("persist-content-123", "Test Content", "Movie")
	if err != nil {
		t.Fatalf("Failed to mark content: %v", err)
	}

	// Close first connection
	db1.Close()

	// Second connection: Verify data persisted
	db2, err := NewDB(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create second database connection: %v", err)
	}
	defer db2.Close()

	// Check subscriber persisted
	isSubscribed, err := db2.IsSubscribed(999888777)
	if err != nil {
		t.Fatalf("Failed to check subscription: %v", err)
	}
	if !isSubscribed {
		t.Error("Expected subscriber to persist across connections")
	}

	// Check content persisted
	isNotified, err := db2.IsContentNotified("persist-content-123")
	if err != nil {
		t.Fatalf("Failed to check content: %v", err)
	}
	if !isNotified {
		t.Error("Expected content to persist across connections")
	}
}
