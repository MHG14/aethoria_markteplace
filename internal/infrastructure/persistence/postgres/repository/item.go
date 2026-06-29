package repository

import (
	"context"

	"github.com/MHG14/aethoria_marketplace/internal/domain/item"
	db "github.com/MHG14/aethoria_marketplace/internal/infrastructure/persistence/postgres/sqlc"
)

func (r *itemRepo) Create(ctx context.Context, i item.Item) (item.Item, error) {
	row, err := r.q.CreateItem(ctx, &db.CreateItemParams{
		Name:    i.Name,
		Type:    db.ItemType(i.Type),
		Status:  db.ItemStatus(i.Status),
		OwnerID: i.OwnerID,
	})
	if err != nil {
		return item.Item{}, err
	}
	return toItem(row), nil
}

func (r *itemRepo) GetByID(ctx context.Context, id int64) (item.Item, error) {
	row, err := r.q.GetItem(ctx, id)
	if err != nil {
		return item.Item{}, err
	}
	return toItem(row), nil
}

func (r *itemRepo) GetByIDForUpdate(ctx context.Context, id int64) (item.Item, error) {
	row, err := r.q.GetItemForUpdate(ctx, id)
	if err != nil {
		return item.Item{}, err
	}
	return toItem(row), nil
}

func (r *itemRepo) List(ctx context.Context) ([]item.Item, error) {
	rows, err := r.q.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]item.Item, len(rows))
	for i, row := range rows {
		items[i] = toItem(row)
	}
	return items, nil
}

func (r *itemRepo) ListByOwner(ctx context.Context, ownerID int64) ([]item.Item, error) {
	rows, err := r.q.ListItemsByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	items := make([]item.Item, len(rows))
	for i, row := range rows {
		items[i] = toItem(row)
	}
	return items, nil
}

func (r *itemRepo) UpdateStatus(ctx context.Context, id int64, status item.Status) (item.Item, error) {
	row, err := r.q.UpdateItemStatus(ctx, &db.UpdateItemStatusParams{
		ID:     id,
		Status: db.ItemStatus(status),
	})
	if err != nil {
		return item.Item{}, err
	}
	return toItem(row), nil
}

func (r *itemRepo) UpdateOwner(ctx context.Context, id int64, ownerID int64) (item.Item, error) {
	row, err := r.q.UpdateItemOwner(ctx, &db.UpdateItemOwnerParams{
		ID:      id,
		OwnerID: ownerID,
	})
	if err != nil {
		return item.Item{}, err
	}
	return toItem(row), nil
}

// toItem converts a sqlc db.Item to a domain item.Item
func toItem(row *db.Item) item.Item {
	return item.Item{
		ID:      row.ID,
		Name:    row.Name,
		Type:    item.Type(row.Type),
		Status:  item.Status(row.Status),
		OwnerID: row.OwnerID,
	}
}
