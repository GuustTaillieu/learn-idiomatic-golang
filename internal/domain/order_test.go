package domain_test

import (
	"fmt"
	"testing"
	"uuid"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

func BenchmarkOrderSerialization(b *testing.B) {
	order := domain.NewOrder(domain.ItemID(uuid.New()), 10)

	b.ResetTimer()
	b.ReportAllocs() // 💡 Tracks exact bytes and allocation count!

	for b.Loop() { // Go 1.24+ idiomatic loop (or for i := 0; i < b.N; i++)
		_ = serializeOrder(order)
	}
}

func serializeOrder(order *domain.Order) []byte {
	// Simulate serialization (e.g., JSON, Protobuf, etc.)
	// For simplicity, we'll just return a byte slice representation of the order ID and quantity.
	return []byte(fmt.Sprintf("%s:%d", order.ID.String(), order.Amount))
}
