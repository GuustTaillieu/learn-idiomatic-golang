# Learning Record 0027: Task Serialization Dilemma & Generic Concurrency

## Date: 2026-09-27

## Milestone & Breakthrough Insight
- The learner attempted to build `Task[T any]` bundling both data (`item T`) and behavior (`Processor[T]`), encountering the fundamental persistence boundary: code/interfaces cannot be serialized to a database column.
- Clarified the two industry solutions:
  1. *Pure in-memory generic concurrency*: `Queue[T]` coordinates arbitrary types in RAM, leaving persistence to domain stores.
  2. *Action/Task Registry pattern* (Celery/Temporal): Persisting task type strings (`"order.fulfill"`) and JSON payloads, resolving handlers from a runtime map.
- Provided a clear path to unblock compilation by restoring domain queue and isolating `lib.Queue[T]` as a pure generic concurrency utility.
