package lib

import "errors"

var ErrTransient = errors.New("transient error, please retry")

func IsRetryable(err error) bool {
	if errors.Is(err, ErrTransient) {
		return true
	}
	return false
}
