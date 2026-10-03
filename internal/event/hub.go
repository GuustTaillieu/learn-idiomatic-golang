package event

import (
	"log/slog"
	"sync"
)

type hub[T any] struct {
	mu          sync.RWMutex
	subscribers map[chan T]struct{}
}

func NewHub[T any]() *hub[T] {
	return &hub[T]{
		subscribers: make(map[chan T]struct{}),
	}
}

func (h *hub[T]) Subscribe() chan T {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan T, 16) // Buffered channel to avoid blocking
	h.subscribers[ch] = struct{}{}
	return ch
}

func (h *hub[T]) Unsubscribe(ch chan T) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.subscribers, ch)
	close(ch)
}

func (h *hub[T]) Broadcast(event T) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.subscribers {
		select {
		case ch <- event:
		default:
			// If the channel is full, we drop the event to avoid blocking
			slog.Warn("subscriber channel is full, dropping event", "event", event)
		}
	}
}
