package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	store := NewStore()

	assert.NotNil(t, store)
	assert.NotNil(t, store.data)
	assert.Equal(t, 1000, store.maxKeys)
	assert.Equal(t, 10*time.Minute, store.ttl)
}

func TestStore_SetAndGet(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		ttl   time.Duration
	}{
		{
			name:  "basic set and get",
			key:   "test-key",
			value: "test-value",
			ttl:   0,
		},
		{
			name:  "set with TTL",
			key:   "ttl-key",
			value: "ttl-value",
			ttl:   1 * time.Minute,
		},
		{
			name:  "empty value",
			key:   "empty",
			value: "",
			ttl:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()

			// Set the value
			err := store.Set(tt.key, tt.value, tt.ttl)
			require.NoError(t, err)

			// Get the value
			value, ok := store.Get(tt.key)
			assert.True(t, ok)
			assert.Equal(t, tt.value, value)
		})
	}
}

func TestStore_GetNonExistent(t *testing.T) {
	store := NewStore()

	value, ok := store.Get("non-existent")
	assert.False(t, ok)
	assert.Empty(t, value)
}

func TestStore_TTLExpiration(t *testing.T) {
	store := NewStore()

	// Set a value with a short TTL
	err := store.Set("expire-key", "expire-value", 50*time.Millisecond)
	require.NoError(t, err)

	// Value should be retrievable immediately
	value, ok := store.Get("expire-key")
	assert.True(t, ok)
	assert.Equal(t, "expire-value", value)

	// Wait for TTL to expire
	time.Sleep(100 * time.Millisecond)

	// Value should no longer be retrievable
	value, ok = store.Get("expire-key")
	assert.False(t, ok)
	assert.Empty(t, value)
}

func TestStore_UpdateExistingKey(t *testing.T) {
	store := NewStore()

	// Set initial value
	err := store.Set("key", "value1", 0)
	require.NoError(t, err)

	// Update the value
	err = store.Set("key", "value2", 0)
	require.NoError(t, err)

	// Get should return the updated value
	value, ok := store.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "value2", value)
}

func TestStore_Delete(t *testing.T) {
	store := NewStore()

	// Set a value
	err := store.Set("delete-key", "delete-value", 0)
	require.NoError(t, err)

	// Verify it exists
	value, ok := store.Get("delete-key")
	assert.True(t, ok)
	assert.Equal(t, "delete-value", value)

	// Delete the value
	err = store.Delete("delete-key")
	require.NoError(t, err)

	// Verify it no longer exists
	value, ok = store.Get("delete-key")
	assert.False(t, ok)
	assert.Empty(t, value)
}

func TestStore_DeleteNonExistent(t *testing.T) {
	store := NewStore()

	// Deleting a non-existent key should not error
	err := store.Delete("non-existent")
	assert.NoError(t, err)
}

func TestStore_MaxKeysLimit(t *testing.T) {
	store := NewStore()
	store.maxKeys = 3 // Set a small limit for testing

	// Fill the store to capacity
	for i := 0; i < 3; i++ {
		err := store.Set(string(rune('a'+i)), "value", 0)
		require.NoError(t, err)
	}

	// Trying to add a new key should fail
	err := store.Set("d", "value", 0)
	assert.ErrorIs(t, err, ErrStoreFull)

	// Updating an existing key should still work
	err = store.Set("a", "new-value", 0)
	require.NoError(t, err)

	value, ok := store.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "new-value", value)
}

func TestStore_MaxKeysZero(t *testing.T) {
	store := NewStore()
	store.maxKeys = 0 // No limit

	// Should be able to add many keys
	for i := 0; i < 100; i++ {
		err := store.Set(string(rune(i)), "value", 0)
		require.NoError(t, err)
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	store := NewStore()

	// Set initial value
	err := store.Set("concurrent", "value", 0)
	require.NoError(t, err)

	done := make(chan bool)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			value, ok := store.Get("concurrent")
			assert.True(t, ok)
			assert.Equal(t, "value", value)
			done <- true
		}()
	}

	// Concurrent writes
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			err := store.Set("key"+string(rune(i)), "value", 0)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}
}
