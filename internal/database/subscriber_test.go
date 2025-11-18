package database

import (
	"os"
	"testing"

	"gorm.io/gorm"
)

// setupTestDB creates a temporary test database
func setupTestDB(t *testing.T) (*DB, func()) {
	// Create temporary database file
	tmpDB := "/tmp/test_jellyfin_bot.db"

	// Clean up any existing test database
	os.Remove(tmpDB)

	db, err := NewDB(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(tmpDB)
	}

	return db, cleanup
}

// Test 1: Add subscriber successfully
func TestAddSubscriber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := db.AddSubscriber(123456789, "testuser", "Test")
	if err != nil {
		t.Fatalf("Failed to add subscriber: %v", err)
	}

	// Verify subscriber was added
	isSubscribed, err := db.IsSubscribed(123456789)
	if err != nil {
		t.Fatalf("Failed to check subscription: %v", err)
	}
	if !isSubscribed {
		t.Error("Expected subscriber to be active")
	}
}

// Test 2: Prevent duplicate subscribers (idempotency)
func TestAddSubscriberDuplicate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID := int64(987654321)

	// Add subscriber first time
	err := db.AddSubscriber(chatID, "user1", "User One")
	if err != nil {
		t.Fatalf("Failed to add subscriber first time: %v", err)
	}

	// Add same subscriber again
	err = db.AddSubscriber(chatID, "user1", "User One")
	if err != nil {
		t.Fatalf("Failed to add subscriber second time: %v", err)
	}

	// Verify only one subscriber exists
	subscribers, err := db.GetAllActiveSubscribers()
	if err != nil {
		t.Fatalf("Failed to get subscribers: %v", err)
	}
	if len(subscribers) != 1 {
		t.Errorf("Expected 1 subscriber, got %d", len(subscribers))
	}
}

// Test 3: Remove subscriber successfully
func TestRemoveSubscriber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID := int64(111222333)

	// Add subscriber
	err := db.AddSubscriber(chatID, "testuser", "Test")
	if err != nil {
		t.Fatalf("Failed to add subscriber: %v", err)
	}

	// Remove subscriber
	err = db.RemoveSubscriber(chatID)
	if err != nil {
		t.Fatalf("Failed to remove subscriber: %v", err)
	}

	// Verify subscriber is not active
	isSubscribed, err := db.IsSubscribed(chatID)
	if err != nil {
		t.Fatalf("Failed to check subscription: %v", err)
	}
	if isSubscribed {
		t.Error("Expected subscriber to be inactive")
	}
}

// Test 4: Remove non-existent subscriber returns error
func TestRemoveNonExistentSubscriber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := db.RemoveSubscriber(999999999)
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

// Test 5: Get all active subscribers
func TestGetAllActiveSubscribers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Add multiple subscribers
	db.AddSubscriber(111, "user1", "User One")
	db.AddSubscriber(222, "user2", "User Two")
	db.AddSubscriber(333, "user3", "User Three")

	// Remove one subscriber
	db.RemoveSubscriber(222)

	// Get active subscribers
	subscribers, err := db.GetAllActiveSubscribers()
	if err != nil {
		t.Fatalf("Failed to get active subscribers: %v", err)
	}

	// Should have 2 active subscribers
	if len(subscribers) != 2 {
		t.Errorf("Expected 2 active subscribers, got %d", len(subscribers))
	}

	// Verify correct chat IDs
	expectedIDs := map[int64]bool{111: true, 333: true}
	for _, chatID := range subscribers {
		if !expectedIDs[chatID] {
			t.Errorf("Unexpected chat ID in active subscribers: %d", chatID)
		}
	}
}

// Test 6: Reactivate removed subscriber
func TestReactivateSubscriber(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	chatID := int64(444555666)

	// Add and remove subscriber
	db.AddSubscriber(chatID, "testuser", "Test")
	db.RemoveSubscriber(chatID)

	// Verify inactive
	isSubscribed, _ := db.IsSubscribed(chatID)
	if isSubscribed {
		t.Error("Expected subscriber to be inactive")
	}

	// Re-add subscriber (should reactivate)
	err := db.AddSubscriber(chatID, "testuser", "Test")
	if err != nil {
		t.Fatalf("Failed to reactivate subscriber: %v", err)
	}

	// Verify active again
	isSubscribed, err = db.IsSubscribed(chatID)
	if err != nil {
		t.Fatalf("Failed to check subscription: %v", err)
	}
	if !isSubscribed {
		t.Error("Expected subscriber to be reactivated")
	}
}
