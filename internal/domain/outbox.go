package domain

import (
	"context"
	"time"
)

type OutboxDispatcher struct {
	orderStore OrderStorer
	queue      OrderQueuer
}

func NewOutboxDispatcher(orderStore OrderStorer, queue OrderQueuer) *OutboxDispatcher {
	return &OutboxDispatcher{
		orderStore: orderStore,
		queue:      queue,
	}
}

func (d *OutboxDispatcher) Start(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := d.DispatchOnce(ctx); err != nil {
				return err
			}
		}
	}
}

func (d *OutboxDispatcher) DispatchOnce(ctx context.Context) error {
	pendingOrders, err := d.orderStore.GetPendingOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range pendingOrders {
		order.Status = OrderStatusRunning
		if err := d.orderStore.Save(ctx, order); err != nil {
			return err
		}
		if err := d.queue.Submit(ctx, order); err != nil {
			return err
		}
	}
	return nil
}
