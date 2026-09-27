package processor

import (
	"context"
	"fmt"
	"sync"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"golang.org/x/sync/errgroup"
)

type Parallel struct {
	processors []domain.Processer
}

func NewParallel(processors ...domain.Processer) *Parallel {
	return &Parallel{
		processors: processors,
	}
}

func (p *Parallel) Process(ctx context.Context, order *domain.Order) (func() error, error) {
	g, gCtx := errgroup.WithContext(ctx)
	rollbackFuncs := make([]func() error, len(p.processors))
	mu := sync.Mutex{}

	for i, processor := range p.processors {
		i, processor := i, processor // capture range variables
		g.Go(func() error {
			rollbackFunc, err := processor.Process(gCtx, order)
			if err != nil {
				return err
			}
			mu.Lock()
			rollbackFuncs[i] = rollbackFunc
			mu.Unlock()
			return nil
		})
	}

	rollbackAll := func() error {
		rollbackErrs := make([]error, 0, len(rollbackFuncs))
		for i := len(rollbackFuncs) - 1; i >= 0; i-- {
			if rollbackFuncs[i] != nil {
				if err := rollbackFuncs[i](); err != nil {
					rollbackErrs = append(rollbackErrs, err)
				}
			}
		}
		if len(rollbackErrs) > 0 {
			return fmt.Errorf("one or more rollbacks failed: %v", rollbackErrs)
		}
		return nil
	}

	return rollbackAll, g.Wait()
}
