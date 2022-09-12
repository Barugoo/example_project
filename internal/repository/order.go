package repository

import (
	"context"
	"database/sql"
	"time"

	"example_project/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, o *models.Order) (*models.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	now := time.Now()

	var id int64
	row := r.db.QueryRowContext(ctx, "INSERT INTO orders (email, created_at) VALUES ($1, $2) RETURNING id",
		order.Email, now)
	if err := row.Err(); err != nil {
		return nil, err
	}
	row.Scan(&id)

	order.ID = id
	order.CreatedAt = now

	return order, nil
}
