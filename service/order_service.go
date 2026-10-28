package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"order-management/model"
	"order-management/repository"
	"order-management/resilience"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req model.CreateOrderRequest) (model.Order, error)
	GetOrder(ctx context.Context, id int) (model.Order, error)
	ListOrders(ctx context.Context) ([]model.Order, error)
	UpdateOrderStatus(ctx context.Context, id int, status string) (model.Order, error)
}

type orderService struct {
	repo       repository.OrderRepository
	httpClient *http.Client
	productURL string
	cb         *resilience.CircuitBreaker
}

func NewOrderService(repo repository.OrderRepository, productURL string) OrderService {
	return &orderService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		productURL: productURL,
		cb:         resilience.NewCircuitBreaker("product-service", 3),
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	price := s.fetchProductPrice(ctx, req.ProductID)
	totalAmount := price * float64(req.Quantity)

	order := model.Order{
		CustomerID:  req.CustomerID,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		Price:       price,
		TotalAmount: totalAmount,
		Status:      model.StatusCreated,
	}

	return s.repo.Create(ctx, order)
}

func (s *orderService) GetOrder(ctx context.Context, id int) (model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) ListOrders(ctx context.Context) ([]model.Order, error) {
	return s.repo.List(ctx)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id int, status string) (model.Order, error) {
	orderStatus := model.OrderStatus(status)

	switch orderStatus {
	case model.StatusCreated, model.StatusConfirmed, model.StatusCancelled:
	default:
		return model.Order{}, fmt.Errorf("invalid status: %s (allowed: CREATED, CONFIRMED, CANCELLED)", status)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	if existing.Status == model.StatusCancelled {
		return model.Order{}, fmt.Errorf("cannot update a cancelled order")
	}

	return s.repo.UpdateStatus(ctx, id, orderStatus)
}

func (s *orderService) fetchProductPrice(ctx context.Context, productID string) float64 {
	var price float64 = 29.99

	err := s.cb.Execute(func() error {
		url := fmt.Sprintf("%s/products/%s", s.productURL, productID)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}

		resp, err := resilience.DoWithRetry(s.httpClient, req, resilience.DefaultRetryConfig())
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var product model.ProductPrice
		if err := json.Unmarshal(body, &product); err != nil {
			return err
		}
		price = product.Price
		return nil
	})

	if err != nil {
		// Graceful degradation — use fallback price
		price = 29.99
	}

	return price
}
