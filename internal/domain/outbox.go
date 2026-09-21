package domain

import (
	"context"
	"time"
)

type OutboxDispatcher struct {
	OrderStore OrderStore
	Queue      *Queue
}

func NewOutboxDispatcher(orderStore OrderStore, queue *Queue) *OutboxDispatcher {
	return &OutboxDispatcher{
		OrderStore: orderStore,
		Queue:      queue,
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
	pendingOrders, err := d.OrderStore.GetPendingOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range pendingOrders {
		if err := d.Queue.Submit(ctx, order); err != nil {
			return err
		}
	}
	return nil
}
