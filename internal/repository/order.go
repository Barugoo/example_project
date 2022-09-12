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

	res, err := r.db.ExecContext(ctx, "INSERT INTO orders (email, created_at) VALUES ($1, $2)",
		order.Email, now)
	if err != nil {
		return nil, err
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	order.ID = orderID
	order.CreatedAt = now

	return order, nil
}
