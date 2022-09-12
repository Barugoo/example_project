package repository

import (
	"context"
	"database/sql"

	"example_project/internal/models"
)

type ItemRepository interface {
	GetItem(ctx context.Context, itemID int64) (*models.Item, error)
	ListItemsAll(ctx context.Context) ([]*models.Item, error)
}

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db}
}

func (r *itemRepository) ListItemsAll(ctx context.Context) (res []*models.Item, err error) {

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, color, created_at FROM items")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i models.Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Color, &i.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, &i)
	}
	return res, nil
}

func (r *itemRepository) GetItem(ctx context.Context, itemID int64) (*models.Item, error) {
	var res models.Item

	row := r.db.QueryRowContext(ctx, "SELECT id, name, color, created_at FROM items WHERE id = $1 LIMIT 1", itemID)
	if row.Err() == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	err := row.Scan(&res.ID, &res.Name, &res.Color, &res.CreatedAt)
	return &res, err
}
