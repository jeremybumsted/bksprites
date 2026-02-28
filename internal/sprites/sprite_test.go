package sprites

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAgentSprite(t *testing.T) {
	// Set environment variable for test
	t.Setenv("SPRITE_API_TOKEN", "test-token")

	sprite := NewAgentSprite("test-sprite")

	assert.NotNil(t, sprite)
	assert.Equal(t, "test-sprite", sprite.Name)
	// Address may be empty if the sprite client can't connect during initialization
	// We just verify the sprite object is created properly
}

func TestNewAgentSprite_EmptyToken(t *testing.T) {
	// Unset the environment variable
	t.Setenv("SPRITE_API_TOKEN", "")

	sprite := NewAgentSprite("test-sprite")

	assert.NotNil(t, sprite)
	assert.Equal(t, "test-sprite", sprite.Name)
}

func TestIsRetryableRunError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "timeout error",
			err:      &testTimeoutError{timeout: true},
			expected: true,
		},
		{
			name:     "i/o timeout string",
			err:      errors.New("i/o timeout"),
			expected: true,
		},
		{
			name:     "I/O TIMEOUT uppercase",
			err:      errors.New("I/O TIMEOUT occurred"),
			expected: true,
		},
		{
			name:     "failed to connect",
			err:      errors.New("failed to connect to server"),
			expected: true,
		},
		{
			name:     "FAILED TO CONNECT uppercase",
			err:      errors.New("FAILED TO CONNECT"),
			expected: true,
		},
		{
			name:     "connection reset by peer",
			err:      errors.New("connection reset by peer"),
			expected: true,
		},
		{
			name:     "CONNECTION RESET BY PEER uppercase",
			err:      errors.New("CONNECTION RESET BY PEER"),
			expected: true,
		},
		{
			name:     "non-retryable error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "authentication error",
			err:      errors.New("authentication failed"),
			expected: false,
		},
		{
			name:     "permission denied",
			err:      errors.New("permission denied"),
			expected: false,
		},
		{
			name:     "network timeout from net package",
			err:      &testNetError{timeout: true, temporary: true},
			expected: true,
		},
		{
			name:     "non-timeout network error",
			err:      &testNetError{timeout: false, temporary: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableRunError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsRetryableRunError_WrappedErrors(t *testing.T) {
	// Test wrapped timeout error
	baseErr := &testTimeoutError{timeout: true}
	wrappedErr := errors.Join(errors.New("wrapper"), baseErr)

	result := isRetryableRunError(wrappedErr)
	assert.True(t, result)
}

func TestConstants(t *testing.T) {
	// Verify the constants are set to expected values
	assert.Equal(t, 5*time.Minute, spriteCommandTimeout)
	assert.Equal(t, 3, spriteRunMaxAttempts)
	assert.Equal(t, 2*time.Second, spriteRetryDelay)
}

// Mock types for testing

// testTimeoutError implements net.Error with configurable timeout
type testTimeoutError struct {
	timeout   bool
	temporary bool
}

func (e *testTimeoutError) Error() string {
	return "test timeout error"
}

func (e *testTimeoutError) Timeout() bool {
	return e.timeout
}

func (e *testTimeoutError) Temporary() bool {
	return e.temporary
}

// testNetError implements net.Error
type testNetError struct {
	timeout   bool
	temporary bool
}

func (e *testNetError) Error() string {
	return "test network error"
}

func (e *testNetError) Timeout() bool {
	return e.timeout
}

func (e *testNetError) Temporary() bool {
	return e.temporary
}

// Ensure our test types implement net.Error
var _ net.Error = (*testTimeoutError)(nil)
var _ net.Error = (*testNetError)(nil)
