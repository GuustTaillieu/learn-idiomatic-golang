package store

import (
	"context"
	"sync"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type ItemMemoryStore struct {
	mu    sync.RWMutex
	items map[domain.ItemID]*domain.Item
}

func NewItemMemoryStore() *ItemMemoryStore {
	return &ItemMemoryStore{
		items: make(map[domain.ItemID]*domain.Item),
	}
}

func (s *ItemMemoryStore) Save(ctx context.Context, item *domain.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[item.ID] = item
	return nil
}

func (s *ItemMemoryStore) Get(ctx context.Context, id domain.ItemID) (*domain.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, domain.ErrItemNotFound
	}

	return item, nil
}

func (s *ItemMemoryStore) GetAll(ctx context.Context) ([]*domain.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]*domain.Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return items, nil
}
