package models

import "gorm.io/gorm"

// ContentCache represents cached content to prevent duplicate notifications
type ContentCache struct {
	gorm.Model
	JellyfinID string `gorm:"uniqueIndex;not null" json:"jellyfin_id"`
	Title      string `json:"title"`
	Type       string `json:"type"` // "Movie" or "Episode"
}

// TableName specifies the table name for ContentCache model
func (ContentCache) TableName() string {
	return "content_cache"
}
