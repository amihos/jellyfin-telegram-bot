package models

import "time"

// JellyfinWebhook represents the payload received from Jellyfin webhook plugin
type JellyfinWebhook struct {
	NotificationType string    `json:"NotificationType"`
	Timestamp        time.Time `json:"Timestamp"`
	ServerID         string    `json:"ServerId"`
	ServerName       string    `json:"ServerName"`
	ServerURL        string    `json:"ServerUrl"`
	ServerVersion    string    `json:"ServerVersion"`
	ItemID           string    `json:"ItemId"`
	ItemName         string    `json:"ItemName"`
	ItemType         string    `json:"ItemType"`
	Year             int       `json:"Year"`
	Overview         string    `json:"Overview"`
	ItemPath         string    `json:"ItemPath"`
	UserName         string    `json:"UserName"`
	UserID           string    `json:"UserId"`

	// Episode-specific fields
	SeriesName    string `json:"SeriesName,omitempty"`
	SeasonNumber  int    `json:"SeasonNumber,omitempty"`
	EpisodeNumber int    `json:"EpisodeNumber,omitempty"`
}

// IsMovie returns true if the webhook is for a movie
func (w *JellyfinWebhook) IsMovie() bool {
	return w.ItemType == "Movie"
}

// IsEpisode returns true if the webhook is for a TV episode
func (w *JellyfinWebhook) IsEpisode() bool {
	return w.ItemType == "Episode"
}

// IsItemAdded returns true if the notification type is ItemAdded
func (w *JellyfinWebhook) IsItemAdded() bool {
	return w.NotificationType == "ItemAdded"
}

// IsValid returns true if the webhook is valid for processing
func (w *JellyfinWebhook) IsValid() bool {
	return w.IsItemAdded() && (w.IsMovie() || w.IsEpisode())
}
