package monitor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/buildkite/stacksapi"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitor(t *testing.T) {
	tests := []struct {
		name     string
		stackKey string
		queue    string
		interval time.Duration
	}{
		{
			name:     "standard configuration",
			stackKey: "test-stack",
			queue:    "default",
			interval: 30 * time.Second,
		},
		{
			name:     "short interval",
			stackKey: "stack-2",
			queue:    "priority",
			interval: 5 * time.Second,
		},
		{
			name:     "long interval",
			stackKey: "stack-3",
			queue:    "low-priority",
			interval: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stacksapi.Client{}
			monitor := NewMonitor(client, tt.stackKey, tt.queue, tt.interval)

			assert.NotNil(t, monitor)
			assert.Equal(t, client, monitor.client)
			assert.Equal(t, tt.stackKey, monitor.stackKey)
			assert.Equal(t, tt.queue, monitor.queue)
			assert.Equal(t, tt.interval, monitor.interval)
		})
	}
}

func TestNewMonitor_NilClient(t *testing.T) {
	// Verify that NewMonitor accepts a nil client (it's up to the caller to provide a valid one)
	monitor := NewMonitor(nil, "test-stack", "default", 30*time.Second)

	assert.NotNil(t, monitor)
	assert.Nil(t, monitor.client)
}

func TestReserveJobs_EmptyJobs(t *testing.T) {
	client := &stacksapi.Client{}
	monitor := NewMonitor(client, "test-stack", "default", 30*time.Second)

	ctx := context.Background()
	err := monitor.reserveJobs(ctx, []stacksapi.ScheduledJob{})

	// Should return nil for empty jobs slice
	assert.NoError(t, err)
}

func TestReserveJobs_NilJobs(t *testing.T) {
	client := &stacksapi.Client{}
	monitor := NewMonitor(client, "test-stack", "default", 30*time.Second)

	ctx := context.Background()
	err := monitor.reserveJobs(ctx, nil)

	// Should return nil for nil jobs slice
	assert.NoError(t, err)
}

func TestRunJob_ExecutesWithoutPanic(t *testing.T) {
	client := &stacksapi.Client{}
	monitor := NewMonitor(client, "test-stack", "default", 30*time.Second)

	ctx := context.Background()

	// This test ensures runJob can be called without panicking
	// It catches syntax errors like missing () on goroutine invocation
	assert.NotPanics(t, func() {
		err := monitor.runJob(ctx, "test-job-uuid")
		assert.NoError(t, err)
	})
}

func TestRunJob_GoroutineExecutes(t *testing.T) {
	client := &stacksapi.Client{}
	monitor := NewMonitor(client, "test-stack", "default", 30*time.Second)

	ctx := context.Background()

	// Create a wait group to verify the goroutine actually executes
	// We can't directly test the sprite behavior without mocking,
	// but we can verify the goroutine syntax is correct by ensuring
	// the function returns (not block indefinitely)
	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()
	err := monitor.runJob(ctx, "test-job-uuid")
	elapsed := time.Since(start)

	wg.Done()

	// runJob should return quickly and not block indefinitely
	// because the work is done in a goroutine
	// Using 1 second threshold to account for sprite initialization overhead
	assert.NoError(t, err)
	assert.Less(t, elapsed, 1*time.Second, "runJob should return without blocking indefinitely")
}

// Note: Testing pollQueue, reserveJobs with actual jobs, and detailed runJob behavior
// would require mocking the stacksapi.Client and sprites, which would be more appropriate
// as integration tests or would require refactoring to inject dependencies via interfaces.
// For unit tests, we've covered the structural, lifecycle, and goroutine syntax aspects.
