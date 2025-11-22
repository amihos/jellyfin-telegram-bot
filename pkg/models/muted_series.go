package models

import "gorm.io/gorm"

// MutedSeries represents a series that a user has muted
type MutedSeries struct {
	gorm.Model
	ChatID     int64  `gorm:"uniqueIndex:idx_chat_series;not null" json:"chat_id"`
	SeriesID   string `gorm:"uniqueIndex:idx_chat_series;not null" json:"series_id"`
	SeriesName string `json:"series_name"`
}

// TableName specifies the table name for MutedSeries model
func (MutedSeries) TableName() string {
	return "muted_series"
}
