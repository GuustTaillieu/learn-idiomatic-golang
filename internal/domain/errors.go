package domain

import "errors"

var (
	ErrTransient     = errors.New("transient error, please retry")
	ErrOrderNotFound = errors.New("order not found")
	ErrQueueClosed   = errors.New("queue is closed")
)

func IsRetryable(err error) bool {
	if errors.Is(err, ErrTransient) {
		return true
	}
	return false
}
