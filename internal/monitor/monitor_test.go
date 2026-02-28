package monitor

import (
	"context"
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
	err := monitor.reserveJobs(ctx, "default", []stacksapi.ScheduledJob{})

	// Should return nil for empty jobs slice
	assert.NoError(t, err)
}

func TestReserveJobs_NilJobs(t *testing.T) {
	client := &stacksapi.Client{}
	monitor := NewMonitor(client, "test-stack", "default", 30*time.Second)

	ctx := context.Background()
	err := monitor.reserveJobs(ctx, "default", nil)

	// Should return nil for nil jobs slice
	assert.NoError(t, err)
}

// Note: Testing pollQueue, reserveJobs with actual jobs, and runJob would require
// mocking the stacksapi.Client and sprites, which would be more appropriate
// as integration tests or would require refactoring to inject dependencies
// via interfaces. For unit tests, we've covered the structural and lifecycle aspects.
