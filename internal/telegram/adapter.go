package telegram

import (
	"context"

	"jellyfin-telegram-bot/internal/jellyfin"
	"jellyfin-telegram-bot/pkg/models"
)

// JellyfinClientAdapter adapts the jellyfin.Client to the telegram.JellyfinClient interface
type JellyfinClientAdapter struct {
	client *jellyfin.Client
}

// NewJellyfinClientAdapter creates a new adapter
func NewJellyfinClientAdapter(client *jellyfin.Client) *JellyfinClientAdapter {
	return &JellyfinClientAdapter{
		client: client,
	}
}

// GetRecentItems adapts GetRecentItems from jellyfin.Client
func (a *JellyfinClientAdapter) GetRecentItems(ctx context.Context, limit int) ([]ContentItem, error) {
	items, err := a.client.GetRecentItems(ctx, limit)
	if err != nil {
		return nil, err
	}

	return convertToTelegramContentItems(items), nil
}

// SearchContent adapts SearchContent from jellyfin.Client
func (a *JellyfinClientAdapter) SearchContent(ctx context.Context, query string, limit int) ([]ContentItem, error) {
	items, err := a.client.SearchContent(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return convertToTelegramContentItems(items), nil
}

// GetPosterImage adapts GetPosterImage from jellyfin.Client
func (a *JellyfinClientAdapter) GetPosterImage(ctx context.Context, itemID string) ([]byte, error) {
	return a.client.GetPosterImage(ctx, itemID)
}

// convertToTelegramContentItems converts jellyfin models to telegram ContentItems
func convertToTelegramContentItems(items []models.ContentItem) []ContentItem {
	result := make([]ContentItem, len(items))
	for i, item := range items {
		result[i] = ContentItem{
			ItemID:          item.ItemID,
			Name:            item.Name,
			Type:            item.Type,
			Overview:        item.Overview,
			CommunityRating: item.CommunityRating,
			OfficialRating:  item.OfficialRating,
			ProductionYear:  item.ProductionYear,
			SeriesName:      item.SeriesName,
			SeasonNumber:    item.SeasonNumber,
			EpisodeNumber:   item.EpisodeNumber,
		}
	}
	return result
}
