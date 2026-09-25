package healthchecker

import "context"

type HealthChecker interface {
	Ping(context context.Context) error
}

type MultiHealthChecker struct {
	checkers []HealthChecker
}

func NewMultiHealthChecker(checkers ...HealthChecker) *MultiHealthChecker {
	return &MultiHealthChecker{checkers: checkers}
}

func (m *MultiHealthChecker) Ping(ctx context.Context) error {
	for _, c := range m.checkers {
		if err := c.Ping(ctx); err != nil {
			return err
		}
	}
	return nil
}
