package scheduler

import (
	"context"
	"log"
	"time"
)

type Task func(context.Context) error

type Runner struct {
	logger   *log.Logger
	interval time.Duration
	timeout  time.Duration
}

func NewRunner(logger *log.Logger, interval, timeout time.Duration) *Runner {
	return &Runner{
		logger:   logger,
		interval: interval,
		timeout:  timeout,
	}
}

func (r *Runner) Start(ctx context.Context, name string, task Task) {
	if r.interval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runCtx := ctx
				cancel := func() {}
				if r.timeout > 0 {
					runCtx, cancel = context.WithTimeout(ctx, r.timeout)
				}

				if err := task(runCtx); err != nil && r.logger != nil {
					r.logger.Printf("%s: %v", name, err)
				}
				cancel()
			}
		}
	}()
}
