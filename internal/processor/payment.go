package processor

import (
	"context"
	"log/slog"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type PaymentProcessor struct{}

func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{}
}

func (p *PaymentProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	slog.Info("Processing payment for order", "orderID", order.ID, "amount", order.Amount)
	time.Sleep(1 * time.Second) // Simulate payment processing delay
	return func() error {
		// Make refund order
		slog.Info("Refunding payment for order", "orderID", order.ID, "amount", order.Amount)
		return nil
	}, nil
}
