package application

import (
	"context"
	"fmt"

	"github.com/MHG14/aethoria_marketplace/internal/domain/item"
)

type ListItemsResponse struct {
	Items []item.Item `json:"items"`
}

// ListItems returns all items in the marketplace
func (a *App) ListItems(ctx context.Context) (ListItemsResponse, error) {
	items, err := a.repos.Item.List(ctx)
	if err != nil {
		return ListItemsResponse{}, fmt.Errorf("list items: %w", err)
	}
	return ListItemsResponse{Items: items}, nil
}

type ListItemsByOwnerResponse struct {
	Items []item.Item `json:"items"`
}

// ListItemsByOwner returns all items owned by a specific guild
func (a *App) ListItemsByOwner(ctx context.Context, ownerID int64) (ListItemsByOwnerResponse, error) {
	items, err := a.repos.Item.ListByOwner(ctx, ownerID)
	if err != nil {
		return ListItemsByOwnerResponse{}, fmt.Errorf("list items by owner: %w", err)
	}
	return ListItemsByOwnerResponse{Items: items}, nil
}
