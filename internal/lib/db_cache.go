package lib

import (
	"context"
	"fmt"

	"golang.org/x/sync/singleflight"
)

type Singleflight[T any] struct {
	sf singleflight.Group
}

func (s *Singleflight[T]) Do(ctx context.Context, key string, fn func() (T, error)) (T, error) {
	var zero T
	result, err, _ := s.sf.Do(key, func() (any, error) {
		return fn()
	})
	if err != nil {
		return zero, err
	}

	res, ok := result.(T)
	if !ok {
		return zero, fmt.Errorf("unexpected type: %T", result)
	}

	return res, nil
}
