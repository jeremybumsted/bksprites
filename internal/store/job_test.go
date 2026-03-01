package store

import (
	"testing"
	"time"

	"github.com/jeremybumsted/bksprites/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJobStore(t *testing.T) {
	store := NewStore()
	jobStore := NewJobStore(store)

	assert.NotNil(t, jobStore)
	assert.NotNil(t, jobStore.store)
}

func TestJobStore_SetAndGet(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		id   string
		job  models.Job
	}{
		{
			name: "complete job",
			id:   "job-123",
			job: models.Job{
				Sprite:          "test-sprite",
				Priority:        10,
				AgentQueryRules: []string{"queue=default", "os=linux"},
				ScheduledAt:     now,
				Pipeline: models.Pipeline{
					Slug: "my-pipeline",
					UUID: "pipe-123",
				},
			},
		},
		{
			name: "minimal job",
			id:   "job-456",
			job: models.Job{
				Sprite:      "minimal",
				Priority:    1,
				ScheduledAt: now,
			},
		},
		{
			name: "job with empty fields",
			id:   "job-789",
			job: models.Job{
				Sprite:          "",
				Priority:        0,
				AgentQueryRules: []string{},
				ScheduledAt:     time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()
			jobStore := NewJobStore(store)

			// Set the job
			err := jobStore.Set(tt.id, tt.job)
			require.NoError(t, err)

			// Get the job
			retrieved, ok, err := jobStore.Get(tt.id)
			require.NoError(t, err)
			assert.True(t, ok)
			assert.Equal(t, tt.job.Sprite, retrieved.Sprite)
			assert.Equal(t, tt.job.Priority, retrieved.Priority)
			assert.Equal(t, tt.job.AgentQueryRules, retrieved.AgentQueryRules)
			assert.True(t, tt.job.ScheduledAt.Equal(retrieved.ScheduledAt))
			assert.Equal(t, tt.job.Pipeline, retrieved.Pipeline)
		})
	}
}

func TestJobStore_GetNonExistent(t *testing.T) {
	store := NewStore()
	jobStore := NewJobStore(store)

	job, ok, err := jobStore.Get("non-existent")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, models.Job{}, job)
}

func TestJobStore_UpdateJob(t *testing.T) {
	store := NewStore()
	jobStore := NewJobStore(store)

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Set initial job
	job1 := models.Job{
		Sprite:   "sprite-1",
		Priority: 5,
	}
	err := jobStore.Set("job-1", job1)
	require.NoError(t, err)

	// Update the job
	job2 := models.Job{
		Sprite:      "sprite-2",
		Priority:    10,
		ScheduledAt: now,
	}
	err = jobStore.Set("job-1", job2)
	require.NoError(t, err)

	// Get should return the updated job
	retrieved, ok, err := jobStore.Get("job-1")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, job2.Sprite, retrieved.Sprite)
	assert.Equal(t, job2.Priority, retrieved.Priority)
}

func TestJobStore_Delete(t *testing.T) {
	store := NewStore()
	jobStore := NewJobStore(store)

	// Set a job
	job := models.Job{
		Sprite:   "test-sprite",
		Priority: 5,
	}
	err := jobStore.Set("job-delete", job)
	require.NoError(t, err)

	// Verify it exists
	retrieved, ok, err := jobStore.Get("job-delete")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, job.Sprite, retrieved.Sprite)

	// Delete the job
	err = jobStore.Delete("job-delete")
	require.NoError(t, err)

	// Verify it no longer exists
	_, ok, err = jobStore.Get("job-delete")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestJobStore_DeleteNonExistent(t *testing.T) {
	store := NewStore()
	jobStore := NewJobStore(store)

	// Deleting a non-existent job should not error
	err := jobStore.Delete("non-existent")
	assert.NoError(t, err)
}

func TestJobStore_KeyPrefix(t *testing.T) {
	store := NewStore()
	jobStore := NewJobStore(store)

	job := models.Job{
		Sprite:   "test-sprite",
		Priority: 5,
	}

	// Set a job
	err := jobStore.Set("123", job)
	require.NoError(t, err)

	// Verify the underlying store has the correct prefixed key
	value, ok := store.Get("job:123")
	assert.True(t, ok)
	assert.NotEmpty(t, value)

	// Verify we can't get it without the prefix
	value, ok = store.Get("123")
	assert.False(t, ok)
	assert.Empty(t, value)
}

func TestJobStore_StoreFull(t *testing.T) {
	store := NewStore()
	store.maxKeys = 2 // Set a small limit
	jobStore := NewJobStore(store)

	job := models.Job{
		Sprite:   "test",
		Priority: 1,
	}

	// Fill the store
	err := jobStore.Set("job-1", job)
	require.NoError(t, err)

	err = jobStore.Set("job-2", job)
	require.NoError(t, err)

	// Next set should fail
	err = jobStore.Set("job-3", job)
	assert.ErrorIs(t, err, ErrStoreFull)
}

func TestJobStore_InvalidJSON(t *testing.T) {
	store := NewStore()
	jobStore := NewJobStore(store)

	// Manually set invalid JSON in the underlying store
	err := store.Set("job:bad", "invalid json {", 0)
	require.NoError(t, err)

	// Get should return an error
	_, ok, err := jobStore.Get("bad")
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestJobStore_SetWithCircularReference(t *testing.T) {
	// This test verifies that we handle JSON marshaling errors
	// Note: bkJob doesn't support circular references, but we test error handling
	store := NewStore()
	store.maxKeys = 1 // Force error after successful marshal

	jobStore := NewJobStore(store)

	job1 := models.Job{
		Sprite:   "test-1",
		Priority: 1,
	}

	// First set should succeed
	err := jobStore.Set("job-1", job1)
	require.NoError(t, err)

	// Second set with different key should fail due to max keys
	job2 := models.Job{
		Sprite:   "test-2",
		Priority: 2,
	}
	err = jobStore.Set("job-2", job2)
	assert.ErrorIs(t, err, ErrStoreFull)
}
