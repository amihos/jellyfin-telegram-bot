package models

// ContentItem represents a movie or episode from Jellyfin
type ContentItem struct {
	ItemID          string  `json:"Id"`
	Name            string  `json:"Name"`
	Type            string  `json:"Type"` // "Movie" or "Episode"
	Overview        string  `json:"Overview"`
	CommunityRating float64 `json:"CommunityRating"`
	OfficialRating  string  `json:"OfficialRating"`
	ProductionYear  int     `json:"ProductionYear"`

	// Episode-specific fields
	SeriesName    string `json:"SeriesName,omitempty"`
	SeasonNumber  int    `json:"ParentIndexNumber,omitempty"`
	EpisodeNumber int    `json:"IndexNumber,omitempty"`
}

// JellyfinItemsResponse represents the response from Jellyfin Items API
type JellyfinItemsResponse struct {
	Items            []ContentItem `json:"Items"`
	TotalRecordCount int           `json:"TotalRecordCount"`
}

// GetDisplayTitle returns the appropriate title for display
func (c *ContentItem) GetDisplayTitle() string {
	if c.Type == "Episode" && c.SeriesName != "" {
		return c.SeriesName
	}
	return c.Name
}

// GetRatingDisplay returns formatted rating string
func (c *ContentItem) GetRatingDisplay() string {
	if c.CommunityRating > 0 {
		return c.OfficialRating
	}
	return "N/A"
}
