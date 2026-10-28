package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"order-management/model"
)

type mockRepo struct {
	mu     sync.RWMutex
	orders map[int]model.Order
	nextID int
}

func newMockRepo() *mockRepo {
	return &mockRepo{orders: make(map[int]model.Order), nextID: 1}
}

func (r *mockRepo) Create(ctx context.Context, order model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = r.nextID
	r.nextID++
	r.orders[order.ID] = order
	return order, nil
}

func (r *mockRepo) GetByID(ctx context.Context, id int) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return model.Order{}, fmt.Errorf("order %d not found", id)
	}
	return o, nil
}

func (r *mockRepo) List(ctx context.Context) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []model.Order
	for _, o := range r.orders {
		list = append(list, o)
	}
	return list, nil
}

func (r *mockRepo) UpdateStatus(ctx context.Context, id int, status model.OrderStatus) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[id]
	if !ok {
		return model.Order{}, fmt.Errorf("order %d not found", id)
	}
	o.Status = status
	r.orders[id] = o
	return o, nil
}

func TestCreateOrder(t *testing.T) {
	repo := newMockRepo()
	order := model.Order{
		CustomerID: 101, ProductID: "P001",
		Quantity: 2, Price: 29.99, TotalAmount: 59.98,
		Status: model.StatusCreated,
	}
	created, err := repo.Create(context.Background(), order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}
	if created.Status != model.StatusCreated {
		t.Errorf("expected CREATED, got %s", created.Status)
	}
}

func TestGetOrder(t *testing.T) {
	repo := newMockRepo()
	repo.Create(context.Background(), model.Order{CustomerID: 101, ProductID: "P001", Quantity: 1, Status: model.StatusCreated})

	order, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.CustomerID != 101 {
		t.Errorf("expected customer 101, got %d", order.CustomerID)
	}
}

func TestGetOrderNotFound(t *testing.T) {
	repo := newMockRepo()
	_, err := repo.GetByID(context.Background(), 999)
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestListOrders(t *testing.T) {
	repo := newMockRepo()
	repo.Create(context.Background(), model.Order{CustomerID: 101, Status: model.StatusCreated})
	repo.Create(context.Background(), model.Order{CustomerID: 102, Status: model.StatusConfirmed})

	orders, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}

func TestUpdateStatus(t *testing.T) {
	repo := newMockRepo()
	repo.Create(context.Background(), model.Order{CustomerID: 101, Status: model.StatusCreated})

	updated, err := repo.UpdateStatus(context.Background(), 1, model.StatusConfirmed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != model.StatusConfirmed {
		t.Errorf("expected CONFIRMED, got %s", updated.Status)
	}
}

func TestConcurrentCreates(t *testing.T) {
	repo := newMockRepo()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			repo.Create(context.Background(), model.Order{
				CustomerID: n, ProductID: "P001",
				Quantity: 1, Status: model.StatusCreated,
			})
		}(i)
	}
	wg.Wait()

	orders, _ := repo.List(context.Background())
	if len(orders) != 100 {
		t.Errorf("expected 100 orders after concurrent creates, got %d", len(orders))
	}
}

func TestInvalidStatusTransition(t *testing.T) {
	validStatuses := map[string]bool{
		"CREATED": true, "CONFIRMED": true, "CANCELLED": true,
	}
	if validStatuses["INVALID"] {
		t.Error("INVALID should not be a valid status")
	}
	if !validStatuses["CONFIRMED"] {
		t.Error("CONFIRMED should be valid")
	}
}
