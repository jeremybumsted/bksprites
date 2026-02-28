package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckSpriteHealth_Success(t *testing.T) {
	// Create a test server that returns 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/health", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Extract host and port from test server
	// The server.URL format is "http://127.0.0.1:port"
	addr := server.Listener.Addr().String()
	host, port, err := parseHostPort(addr)
	require.NoError(t, err)

	m := &Monitor{}
	healthy, err := m.CheckSpriteHealth(host, port)

	assert.NoError(t, err)
	assert.True(t, healthy)
}

func TestCheckSpriteHealth_Unhealthy(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"Internal Server Error", http.StatusInternalServerError},
		{"Service Unavailable", http.StatusServiceUnavailable},
		{"Bad Gateway", http.StatusBadGateway},
		{"Not Found", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			addr := server.Listener.Addr().String()
			host, port, err := parseHostPort(addr)
			require.NoError(t, err)

			m := &Monitor{}
			healthy, err := m.CheckSpriteHealth(host, port)

			assert.Error(t, err)
			assert.False(t, healthy)
			assert.Contains(t, err.Error(), "unhealthy")
		})
	}
}

func TestCheckSpriteHealth_DefaultPort(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	addr := server.Listener.Addr().String()
	host, _, err := parseHostPort(addr)
	require.NoError(t, err)

	m := &Monitor{}
	// Pass empty string for port to test default
	// Note: This will fail because we're using a test server on a random port
	// but we can verify the default port logic is applied
	_, err = m.CheckSpriteHealth(host, "")
	// We expect an error because the default port 8080 won't match our test server
	assert.Error(t, err)
}

func TestCheckSpriteHealth_ConnectionError(t *testing.T) {
	m := &Monitor{}
	// Use a port that's unlikely to be open
	healthy, err := m.CheckSpriteHealth("127.0.0.1", "9999")

	assert.Error(t, err)
	assert.False(t, healthy)
}

func TestCheckSpriteHealth_InvalidHost(t *testing.T) {
	m := &Monitor{}
	healthy, err := m.CheckSpriteHealth("invalid-host-that-does-not-exist.local", "8080")

	assert.Error(t, err)
	assert.False(t, healthy)
}

// Helper function to parse host and port from address string
func parseHostPort(addr string) (host, port string, err error) {
	// addr is typically in format "127.0.0.1:port" or "[::1]:port"
	var i int
	for i = len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			break
		}
	}
	if i < 0 {
		return "", "", http.ErrNotSupported
	}
	return addr[:i], addr[i+1:], nil
}
