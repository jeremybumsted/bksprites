// Package monitor watches the configured Buildkite queue on an interval configured by `interval`
package monitor

import (
	"context"
	"errors"
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
			if err := m.reserveJobs(ctx, queueKey, resp.Jobs); err != nil {
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

	reserveRequest := stacksapi.BatchReserveJobsRequest{
		StackKey:                 m.stackKey,
		JobUUIDs:                 jobUUIDs,
		ReservationExpirySeconds: 30, // Let's default to 30, but this can be a config value later
	}

	resp, _, err := m.client.BatchReserveJobs(ctx, reserveRequest)
	if err != nil {
		log.Error("failed to reserve jobs", "error", err)
	}
	if len(resp.NotReserved) > 0 {
		log.Warn("Some jobs were not reserved", "Not Reserved", resp.NotReserved)
	}
	if len(resp.Reserved) > 0 {
		log.Info("Reserved jobs, running on agent sprites")

		for i := 0; i < len(resp.Reserved); i++ {
			err = m.runJob(ctx, resp.Reserved[i])
			if err != nil {
				log.Error("error running jobs", "error", err)
			}
		}
	}
	return nil
}

func (m *Monitor) runJob(ctx context.Context, jobUUID string) error {
	// This will eventually get farmed out to a sprite registry
	spr := sprites.NewAgentSprite("bk-test-1")

	healthy, err := m.CheckSpriteHealth(spr.Address, "8080")
	if err != nil || !healthy {
		log.Error("big ol problem I think", "error", err)
		return err
	}
	if healthy {
		err = spr.RunJob(jobUUID)
		if err != nil {
			log.Error("there was an error running the job", "error", err)
			err = m.finishJob(ctx, m.queue, jobUUID, fmt.Sprintf("%v", err))
			if err != nil {
				log.Error("failed to fail job", "error", err)
			}
		}

		log.Info("ran job", "job", "jobUUID")
		return nil
	}
	return errors.New("funchere was an error running the job")
}

// finishJob returns a status back to Buildkite to surface failures starting an agent
func (m *Monitor) finishJob(ctx context.Context, queueKey string, job string, msg string) error {
	req := stacksapi.FinishJobRequest{
		StackKey:   m.stackKey,
		JobUUID:    job,
		ExitStatus: -1,
		Detail:     msg,
	}
	_, err := m.client.FinishJob(ctx, req)
	if err != nil {
		log.Error("failed to finish the job", "error", err)
	}
	return nil
}
