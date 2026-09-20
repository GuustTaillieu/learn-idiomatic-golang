package store

import (
	"context"
	"sync"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type OrderMemoryStore struct {
	mu     sync.RWMutex
	orders map[domain.OrderID]*domain.Order
}

func NewOrderMemoryStore() *OrderMemoryStore {
	return &OrderMemoryStore{
		orders: make(map[domain.OrderID]*domain.Order),
	}
}

func (s *OrderMemoryStore) Save(ctx context.Context, order *domain.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.ID] = order
	return nil
}

func (s *OrderMemoryStore) Get(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderMemoryStore) GetAll(ctx context.Context) ([]*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}

	return orders, nil
}
