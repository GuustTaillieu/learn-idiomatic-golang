package memory

import (
	"context"
	"sync"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type orderStore struct {
	mu     sync.RWMutex
	orders map[domain.OrderID]*domain.Order
}

func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[domain.OrderID]*domain.Order),
	}
}

func (s *orderStore) Save(ctx context.Context, order *domain.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.ID] = order
	return nil
}

func (s *orderStore) Get(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

func (s *orderStore) GetAll(ctx context.Context) ([]*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *orderStore) GetPendingOrders(ctx context.Context) ([]*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pendingOrders := make([]*domain.Order, 0)
	for _, order := range s.orders {
		if order.Status == domain.OrderStatusPending {
			pendingOrders = append(pendingOrders, order)
		}
	}

	return pendingOrders, nil
}
