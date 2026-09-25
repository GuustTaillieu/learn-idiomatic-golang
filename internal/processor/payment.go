package processor

import (
	"context"
	"time"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

type PaymentProcessor struct{}

func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{}
}

func (p *PaymentProcessor) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	lib.Logger(ctx).Info("Processing payment for order", "orderID", order.ID, "amount", order.Amount)
	time.Sleep(1 * time.Second) // Simulate payment processing delay
	return func() error {
		// Make refund order
		lib.Logger(ctx).Info("Refunding payment for order", "orderID", order.ID, "amount", order.Amount)
		return nil
	}, nil
}
