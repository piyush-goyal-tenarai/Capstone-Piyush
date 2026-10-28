package model

import "time"

type OrderStatus string

const (
	StatusCreated   OrderStatus = "CREATED"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID          int         `json:"id"`
	CustomerID  int         `json:"customerId"`
	ProductID   string      `json:"productId"`
	Quantity    int         `json:"quantity"`
	Price       float64     `json:"price"`
	TotalAmount float64     `json:"totalAmount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type CreateOrderRequest struct {
	CustomerID int    `json:"customerId" binding:"required,min=1"`
	ProductID  string `json:"productId" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"requestId,omitempty"`
}

type ProductPrice struct {
	ProductID string  `json:"productId"`
	Price     float64 `json:"price"`
}
