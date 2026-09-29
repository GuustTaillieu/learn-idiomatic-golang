package health

import "context"

type Checker interface {
	Ping(context context.Context) error
}

type MultiChecker struct {
	checkers []Checker
}

func NewMultiChecker(checkers ...Checker) *MultiChecker {
	return &MultiChecker{checkers: checkers}
}

func (m *MultiChecker) Ping(ctx context.Context) error {
	for _, c := range m.checkers {
		if err := c.Ping(ctx); err != nil {
			return err
		}
	}
	return nil
}
