# Channel Signaling: Broadcast via close() vs Single Send

The learner introduced `stopChan` to decouple shutdown from mutex contention. A key concurrency principle was illuminated:
1. Sending a value `ch <- struct{}{}` to an unbuffered channel is a 1-to-1 rendezvous. If no goroutine is actively receiving, it blocks forever (causing deadlocks during `Stop()`). Even with receivers, it only unblocks a single goroutine.
2. Closing a channel `close(ch)` is Go's universal broadcast mechanism. It immediately and permanently unblocks all current and future receivers without blocking the sender.
