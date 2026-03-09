package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobMarshaling(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		job  Job
	}{
		{
			name: "complete job with all fields",
			job: Job{
				Sprite:          "test-sprite",
				Priority:        10,
				AgentQueryRules: []string{"queue=default", "os=linux"},
				ScheduledAt:     now,
				Pipeline: Pipeline{
					Slug: "my-pipeline",
					UUID: "pipe-123",
				},
			},
		},
		{
			name: "minimal job",
			job: Job{
				Sprite:      "minimal",
				Priority:    1,
				ScheduledAt: now,
			},
		},
		{
			name: "job with empty agent query rules",
			job: Job{
				Sprite:          "test",
				Priority:        5,
				AgentQueryRules: []string{},
				ScheduledAt:     now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.job)
			require.NoError(t, err)

			// Unmarshal back to Job
			var unmarshaled Job
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			// Verify the unmarshaled job matches the original
			assert.Equal(t, tt.job.Sprite, unmarshaled.Sprite)
			assert.Equal(t, tt.job.Priority, unmarshaled.Priority)
			assert.Equal(t, tt.job.AgentQueryRules, unmarshaled.AgentQueryRules)
			assert.True(t, tt.job.ScheduledAt.Equal(unmarshaled.ScheduledAt))
			assert.Equal(t, tt.job.Pipeline, unmarshaled.Pipeline)
		})
	}
}

func TestPipelineMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		pipeline Pipeline
	}{
		{
			name: "complete pipeline",
			pipeline: Pipeline{
				Slug: "my-pipeline",
				UUID: "uuid-123",
			},
		},
		{
			name:     "empty pipeline",
			pipeline: Pipeline{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.pipeline)
			require.NoError(t, err)

			var unmarshaled Pipeline
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, tt.pipeline, unmarshaled)
		})
	}
}

func TestBuildMarshaling(t *testing.T) {
	tests := []struct {
		name  string
		build Build
	}{
		{
			name: "complete build",
			build: Build{
				Number: 42,
				Branch: "main",
				UUID:   "build-uuid-123",
			},
		},
		{
			name:  "empty build",
			build: Build{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.build)
			require.NoError(t, err)

			var unmarshaled Build
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, tt.build, unmarshaled)
		})
	}
}

func TestStepMarshaling(t *testing.T) {
	tests := []struct {
		name string
		step Step
	}{
		{
			name: "step with key",
			step: Step{
				Key: "test-step",
			},
		},
		{
			name: "empty step",
			step: Step{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.step)
			require.NoError(t, err)

			var unmarshaled Step
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, tt.step, unmarshaled)
		})
	}
}
