package scheduler

import (
	"bytes"
	"context"
	"errors"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartDoesNothingWhenIntervalIsNotPositive(t *testing.T) {
	var calls atomic.Int32
	runner := NewRunner(nil, 0, 0)

	runner.Start(context.Background(), "noop", func(context.Context) error {
		calls.Add(1)
		return nil
	})

	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, int32(0), calls.Load())
}

func TestStartRunsTaskWithoutTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var calls atomic.Int32
	runner := NewRunner(nil, 10*time.Millisecond, 0)

	runner.Start(ctx, "tick", func(taskCtx context.Context) error {
		assert.Same(t, ctx, taskCtx)
		calls.Add(1)
		cancel()
		return nil
	})

	require.Eventually(t, func() bool { return calls.Load() == 1 }, 200*time.Millisecond, 10*time.Millisecond)
}

func TestStartLogsTaskErrorsAndAppliesTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		calls  atomic.Int32
		buffer bytes.Buffer
	)
	runner := NewRunner(log.New(&buffer, "", 0), 10*time.Millisecond, 15*time.Millisecond)

	runner.Start(ctx, "refresh", func(taskCtx context.Context) error {
		calls.Add(1)
		<-taskCtx.Done()
		cancel()
		return errors.New("boom")
	})

	require.Eventually(t, func() bool { return calls.Load() == 1 }, 300*time.Millisecond, 10*time.Millisecond)
	require.Eventually(t, func() bool { return buffer.Len() > 0 }, 300*time.Millisecond, 10*time.Millisecond)
	assert.Contains(t, buffer.String(), "refresh: boom")
}
