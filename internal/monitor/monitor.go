// Package monitor watches the configured Buildkite queue on an interval configured by `interval`
package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/buildkite/stacksapi"
	"github.com/charmbracelet/log"

	"github.com/jeremybumsted/bksprites/internal/sprites"
)

type Monitor struct {
	client   *stacksapi.Client
	stackKey string
	queue    string
	interval time.Duration
}

func NewMonitor(client *stacksapi.Client, stackKey string, queue string, interval time.Duration) *Monitor {
	return &Monitor{
		client:   client,
		stackKey: stackKey,
		queue:    queue,
		interval: interval,
	}
}

func (m *Monitor) Start(ctx context.Context) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	log.Info(fmt.Sprintf("Starting monitor for queue: %s", m.queue))

	for {
		select {
		case <-ctx.Done():
			log.Info("Monitor shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := m.pollQueue(ctx, m.queue); err != nil {
				log.Error("Error polling queue", "error", err)
			}
		}
	}
}

func (m *Monitor) pollQueue(ctx context.Context, queueKey string) error {
	var cursor string
	jobsProcessed := 0

	for {
		resp, _, err := m.client.ListScheduledJobs(ctx, stacksapi.ListScheduledJobsRequest{
			StackKey:        m.stackKey,
			ClusterQueueKey: queueKey,
			PageSize:        50,
			StartCursor:     cursor,
		})
		if err != nil {
			return fmt.Errorf("listing scheduled jobs: %w", err)
		}

		if resp.ClusterQueue.Paused {
			log.Info("Queue is paused, skipping")
			return nil
		}

		if len(resp.Jobs) > 0 {
			if err := m.runJobs(ctx, queueKey, resp.Jobs); err != nil {
				log.Error("Error reserving jobs", "error", err)
			} else {
				jobsProcessed += len(resp.Jobs)
			}
		}

		if !resp.PageInfo.HasNextPage {
			break
		}
		cursor = resp.PageInfo.EndCursor
	}
	if jobsProcessed > 0 {
		log.Info(fmt.Sprintf("Processed %v jobs on queue %v", jobsProcessed, queueKey))
	}
	return nil
}

func (m *Monitor) reserveJobs(ctx context.Context, queueKey string, jobs []stacksapi.ScheduledJob) error {
	if len(jobs) == 0 {
		return nil
	}

	jobUUIDs := make([]string, len(jobs))
	for i, job := range jobs {
		jobUUIDs[i] = job.ID
	}

	log.Info("This is our collection of jobUUIDs", jobUUIDs)
	return nil
}

func (m *Monitor) runJobs(ctx context.Context, queueKey string, jobs []stacksapi.ScheduledJob) error {
	if len(jobs) == 0 {
		return nil
	}

	jobUUIDs := make([]string, len(jobs))
	for i, job := range jobs {
		jobUUIDs[i] = job.ID
		err := sprites.RunJob(job.ID)
		if err != nil {
			return err
		}
	}

	log.Info("ran these jobs: ", "jobUUIDs", jobUUIDs)
	return nil
}
