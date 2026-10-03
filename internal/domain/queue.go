package domain

import "context"

type OrderQueuer interface {
	Submit(ctx context.Context, order *Order) error
	Start(ctx context.Context, numWorkers int)
	Stop()
}

type Processor[T any] interface {
	Process(ctx context.Context, item T) (cleanup func(context.Context) error, err error)
}
