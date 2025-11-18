package jellyfin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"jellyfin-telegram-bot/pkg/models"
)

// Client represents a Jellyfin API client
type Client struct {
	serverURL  string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Jellyfin API client
func NewClient(serverURL, apiKey string) *Client {
	return &Client{
		serverURL: serverURL,
		apiKey:    apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithHTTPClient creates a new Jellyfin API client with a custom HTTP client
func NewClientWithHTTPClient(serverURL, apiKey string, httpClient *http.Client) *Client {
	return &Client{
		serverURL:  serverURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// doRequest performs an HTTP request with authentication headers
func (c *Client) doRequest(ctx context.Context, method, path string, params url.Values) (*http.Response, error) {
	// Build URL with query parameters
	u := c.serverURL + path
	if params != nil {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	req.Header.Set("X-Emby-Token", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Handle HTTP errors
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		return nil, fmt.Errorf("authentication failed: invalid API key")
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, fmt.Errorf("resource not found")
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// GetPosterImage fetches the primary poster image for a given item ID
func (c *Client) GetPosterImage(ctx context.Context, itemID string) ([]byte, error) {
	path := fmt.Sprintf("/Items/%s/Images/Primary", itemID)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch poster image: %w", err)
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imageBytes, nil
}

// GetRecentItems fetches recently added movies and episodes
func (c *Client) GetRecentItems(ctx context.Context, limit int) ([]models.ContentItem, error) {
	params := url.Values{}
	params.Set("Filters", "IsNotFolder")
	params.Set("Recursive", "true")
	params.Set("SortBy", "DateCreated")
	params.Set("SortOrder", "Descending")
	params.Set("IncludeItemTypes", "Movie,Episode")
	params.Set("Limit", strconv.Itoa(limit))
	params.Set("Fields", "Overview,CommunityRating,OfficialRating,ProductionYear")

	resp, err := c.doRequest(ctx, "GET", "/Items", params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent items: %w", err)
	}
	defer resp.Body.Close()

	var result models.JellyfinItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Items, nil
}

// SearchContent searches for movies and episodes matching the query
func (c *Client) SearchContent(ctx context.Context, query string, limit int) ([]models.ContentItem, error) {
	params := url.Values{}
	params.Set("SearchTerm", query)
	params.Set("Recursive", "true")
	params.Set("IncludeItemTypes", "Movie,Episode")
	params.Set("Limit", strconv.Itoa(limit))
	params.Set("Fields", "Overview,CommunityRating,OfficialRating,ProductionYear")

	resp, err := c.doRequest(ctx, "GET", "/Items", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search content: %w", err)
	}
	defer resp.Body.Close()

	var result models.JellyfinItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Items, nil
}
