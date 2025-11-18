package jellyfin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestClientAuthentication tests that the client includes authentication headers
func TestClientAuthentication(t *testing.T) {
	apiKey := "test-api-key-12345"
	var receivedToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken = r.Header.Get("X-Emby-Token")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Items":[],"TotalRecordCount":0}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, apiKey)
	_, err := client.GetRecentItems(context.Background(), 10)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if receivedToken != apiKey {
		t.Errorf("Expected token '%s', got '%s'", apiKey, receivedToken)
	}
}

// TestGetPosterImageSuccess tests successful image fetching
func TestGetPosterImageSuccess(t *testing.T) {
	expectedImage := []byte("fake-image-data")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/Items/") || !strings.Contains(r.URL.Path, "/Images/Primary") {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(expectedImage)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	imageBytes, err := client.GetPosterImage(context.Background(), "item123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if string(imageBytes) != string(expectedImage) {
		t.Errorf("Expected image '%s', got '%s'", expectedImage, imageBytes)
	}
}

// TestGetPosterImageNotFound tests handling of missing images
func TestGetPosterImageNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	_, err := client.GetPosterImage(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("Expected error for 404, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

// TestGetRecentItemsSuccess tests fetching recent content
func TestGetRecentItemsSuccess(t *testing.T) {
	mockResponse := `{
		"Items": [
			{
				"Id": "movie1",
				"Name": "Test Movie",
				"Type": "Movie",
				"Overview": "A test movie",
				"CommunityRating": 8.5,
				"OfficialRating": "PG-13",
				"ProductionYear": 2023
			},
			{
				"Id": "episode1",
				"Name": "Test Episode",
				"Type": "Episode",
				"SeriesName": "Test Series",
				"ParentIndexNumber": 1,
				"IndexNumber": 1,
				"Overview": "A test episode",
				"CommunityRating": 9.0,
				"ProductionYear": 2023
			}
		],
		"TotalRecordCount": 2
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if r.URL.Query().Get("SortBy") != "DateCreated" {
			t.Errorf("Expected SortBy=DateCreated")
		}
		if r.URL.Query().Get("IncludeItemTypes") != "Movie,Episode" {
			t.Errorf("Expected IncludeItemTypes=Movie,Episode")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	items, err := client.GetRecentItems(context.Background(), 10)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	// Verify movie data
	if items[0].Name != "Test Movie" {
		t.Errorf("Expected name 'Test Movie', got '%s'", items[0].Name)
	}
	if items[0].Type != "Movie" {
		t.Errorf("Expected type 'Movie', got '%s'", items[0].Type)
	}
	if items[0].CommunityRating != 8.5 {
		t.Errorf("Expected rating 8.5, got %f", items[0].CommunityRating)
	}

	// Verify episode data
	if items[1].SeriesName != "Test Series" {
		t.Errorf("Expected series name 'Test Series', got '%s'", items[1].SeriesName)
	}
	if items[1].SeasonNumber != 1 {
		t.Errorf("Expected season 1, got %d", items[1].SeasonNumber)
	}
}

// TestSearchContentSuccess tests search functionality
func TestSearchContentSuccess(t *testing.T) {
	mockResponse := `{
		"Items": [
			{
				"Id": "search1",
				"Name": "Interstellar",
				"Type": "Movie",
				"Overview": "Space movie",
				"CommunityRating": 8.6,
				"ProductionYear": 2014
			}
		],
		"TotalRecordCount": 1
	}`

	searchTerm := "interstellar"
	var receivedSearchTerm string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedSearchTerm = r.URL.Query().Get("SearchTerm")
		if r.URL.Query().Get("IncludeItemTypes") != "Movie,Episode" {
			t.Errorf("Expected IncludeItemTypes=Movie,Episode")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	items, err := client.SearchContent(context.Background(), searchTerm, 10)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if receivedSearchTerm != searchTerm {
		t.Errorf("Expected search term '%s', got '%s'", searchTerm, receivedSearchTerm)
	}

	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	if items[0].Name != "Interstellar" {
		t.Errorf("Expected name 'Interstellar', got '%s'", items[0].Name)
	}
}

// TestSearchContentPersian tests search with Persian characters
func TestSearchContentPersian(t *testing.T) {
	persianQuery := "فیلم"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		searchTerm := r.URL.Query().Get("SearchTerm")
		if searchTerm != persianQuery {
			t.Errorf("Persian query not preserved: got '%s'", searchTerm)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Items":[],"TotalRecordCount":0}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	_, err := client.SearchContent(context.Background(), persianQuery, 10)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// TestUnauthorizedError tests handling of authentication errors
func TestUnauthorizedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "invalid-key")
	_, err := client.GetRecentItems(context.Background(), 10)

	if err == nil {
		t.Fatal("Expected error for 401, got nil")
	}

	if !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("Expected 'authentication failed' error, got: %v", err)
	}
}

// TestContextTimeout tests context timeout handling
func TestContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		select {
		case <-r.Context().Done():
			return
		case <-time.After(100 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	// Create client with very short timeout
	httpClient := &http.Client{Timeout: 10 * time.Millisecond}
	client := NewClientWithHTTPClient(server.URL, "test-key", httpClient)

	_, err := client.GetRecentItems(context.Background(), 10)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}
}
