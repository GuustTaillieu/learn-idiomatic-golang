package memory

import (
	"context"
	"sync"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type itemStore struct {
	mu    sync.RWMutex
	items map[domain.ItemID]*domain.Item
}

func NewItemStore() *itemStore {
	return &itemStore{
		items: make(map[domain.ItemID]*domain.Item),
	}
}

func (s *itemStore) Save(ctx context.Context, item *domain.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[item.ID] = item
	return nil
}

func (s *itemStore) Get(ctx context.Context, id domain.ItemID) (*domain.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, domain.ErrItemNotFound
	}

	return item, nil
}

func (s *itemStore) GetAll(ctx context.Context) ([]*domain.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]*domain.Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return items, nil
}
