# Implemented Channel Enqueuing with RWMutex and Context

The learner independently implemented `Submit(ctx context.Context, task Task) error`, demonstrating correct use of `q.mu.RLock()` / `defer q.mu.RUnlock()`, checking `q.closed` before interacting with the channel, selecting on `q.tasks <- task` vs `<-ctx.Done()`, and wrapping errors with `%w`. This confirms the learner understands shared state protection and non-blocking channel backpressure.
