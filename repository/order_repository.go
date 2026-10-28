package repository

import (
	"context"
	"database/sql"
	"fmt"

	"order-management/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (model.Order, error)
	GetByID(ctx context.Context, id int) (model.Order, error)
	List(ctx context.Context) ([]model.Order, error)
	UpdateStatus(ctx context.Context, id int, status model.OrderStatus) (model.Order, error)
}

type postgresOrderRepo struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepo{db: db}
}

func (r *postgresOrderRepo) Create(ctx context.Context, order model.Order) (model.Order, error) {
	query := `
		INSERT INTO orders (customer_id, product_id, quantity, price, total_amount, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, customer_id, product_id, quantity, price, total_amount, status, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		order.CustomerID, order.ProductID, order.Quantity,
		order.Price, order.TotalAmount, order.Status,
	).Scan(
		&order.ID, &order.CustomerID, &order.ProductID,
		&order.Quantity, &order.Price, &order.TotalAmount,
		&order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("create order: %w", err)
	}
	return order, nil
}

func (r *postgresOrderRepo) GetByID(ctx context.Context, id int) (model.Order, error) {
	query := `
		SELECT id, customer_id, product_id, quantity, price, total_amount, status, created_at, updated_at
		FROM orders WHERE id = $1`

	var order model.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.CustomerID, &order.ProductID,
		&order.Quantity, &order.Price, &order.TotalAmount,
		&order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return model.Order{}, fmt.Errorf("order %d not found", id)
	}
	if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}
	return order, nil
}

func (r *postgresOrderRepo) List(ctx context.Context) ([]model.Order, error) {
	query := `
		SELECT id, customer_id, product_id, quantity, price, total_amount, status, created_at, updated_at
		FROM orders ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(
			&o.ID, &o.CustomerID, &o.ProductID,
			&o.Quantity, &o.Price, &o.TotalAmount,
			&o.Status, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *postgresOrderRepo) UpdateStatus(ctx context.Context, id int, status model.OrderStatus) (model.Order, error) {
	query := `
		UPDATE orders SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, customer_id, product_id, quantity, price, total_amount, status, created_at, updated_at`

	var order model.Order
	err := r.db.QueryRowContext(ctx, query, status, id).Scan(
		&order.ID, &order.CustomerID, &order.ProductID,
		&order.Quantity, &order.Price, &order.TotalAmount,
		&order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return model.Order{}, fmt.Errorf("order %d not found", id)
	}
	if err != nil {
		return model.Order{}, fmt.Errorf("update status: %w", err)
	}
	return order, nil
}
