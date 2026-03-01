// Package monitor watches the configured Buildkite queue on an interval configured by `interval`
package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/buildkite/stacksapi"
	"github.com/charmbracelet/log"

	"github.com/jeremybumsted/bksprites/internal/models"
	"github.com/jeremybumsted/bksprites/internal/sprites"
	"github.com/jeremybumsted/bksprites/internal/store"
)

type Monitor struct {
	client   *stacksapi.Client
	stackKey string
	queue    string
	interval time.Duration
	jobStore *store.JobStore
}

func NewMonitor(client *stacksapi.Client, stackKey string, queue string, interval time.Duration) *Monitor {
	s := store.NewStore()
	js := store.NewJobStore(s)

	return &Monitor{
		client:   client,
		stackKey: stackKey,
		queue:    queue,
		interval: interval,
		jobStore: js,
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
			jobList, err := m.pollQueue(ctx, m.queue)
			if err != nil {
				log.Error("Error polling queue", "error", err)
			}
			if err = m.reserveJobs(ctx, jobList); err != nil {
				log.Error("Error reserving jobs", "error", err)
			}

		}
	}
}

func (m *Monitor) pollQueue(ctx context.Context, queueKey string) ([]stacksapi.ScheduledJob, error) {
	var cursor string
	var jobList []stacksapi.ScheduledJob
	jobsProcessed := 0

	for {
		resp, _, err := m.client.ListScheduledJobs(ctx, stacksapi.ListScheduledJobsRequest{
			StackKey:        m.stackKey,
			ClusterQueueKey: queueKey,
			PageSize:        50,
			StartCursor:     cursor,
		})
		if err != nil {
			return nil, fmt.Errorf("listing scheduled jobs: %w", err)
		}

		if resp.ClusterQueue.Paused {
			log.Info("Queue is paused, skipping")
			return nil, nil
		}

		if len(resp.Jobs) > 0 {
			jobList = append(jobList, resp.Jobs...)
			jobsProcessed += len(resp.Jobs)
		}

		if !resp.PageInfo.HasNextPage {
			break
		}
		cursor = resp.PageInfo.EndCursor
	}
	if jobsProcessed > 0 {
		log.Info(fmt.Sprintf("Processed %v jobs on queue %v", jobsProcessed, queueKey))
	}
	return jobList, nil
}

func (m *Monitor) reserveJobs(ctx context.Context, jobs []stacksapi.ScheduledJob) error {
	if len(jobs) == 0 {
		return nil
	}

	log.Info("we're in reserveJobs now", "job slice length", len(jobs))

	jobUUIDs := make([]string, len(jobs))
	for i, job := range jobs {
		bkJob := models.Job{
			Priority:        job.Priority,
			AgentQueryRules: job.AgentQueryRules,
			ScheduledAt:     job.ScheduledAt,
			Pipeline: models.Pipeline{
				Slug: job.Pipeline.Slug,
				UUID: job.Pipeline.UUID,
			},
			Build: models.Build{
				Number: job.Build.Number,
				Branch: job.Build.Branch,
				UUID:   job.Build.UUID,
			},
			Step: models.Step{
				Key: job.Step.Key,
			},
		}
		if err := m.jobStore.Set(job.ID, bkJob); err != nil {
			return err
		}
		jobUUIDs[i] = job.ID
	}

	reserveRequest := stacksapi.BatchReserveJobsRequest{
		StackKey:                 m.stackKey,
		JobUUIDs:                 jobUUIDs,
		ReservationExpirySeconds: 30, // Let's default to 30, but this can be a config value later, realistically it shouldn't take more than 30 seconds to start a job.
	}

	resp, _, err := m.client.BatchReserveJobs(ctx, reserveRequest)
	if err != nil {
		log.Error("failed to reserve jobs", "error", err)
	}
	if len(resp.NotReserved) > 0 {
		for i := 0; i < len(resp.NotReserved); i++ {
			job := resp.NotReserved[i]
			if err = m.jobStore.Delete(job); err != nil {
				log.Error("failed to delete job from the job store, but the job is not reserved on Buildkite", "error", err)
			}
		}
		log.Warn("Some jobs were not reserved", "Not Reserved", resp.NotReserved)
	}
	if len(resp.Reserved) > 0 {
		log.Info("We should be running jobs now?")
		for i := 0; i < len(resp.Reserved); i++ {
			job := resp.Reserved[i]
			log.Info("Running this job: ", "uuid", job)
			err = m.runJob(ctx, job)
			if err != nil {
				log.Error("error running jobs", "error", err)
			}
			if err = m.jobStore.Delete(job); err != nil {
				log.Error("failed to delete job from the job store, but the job finished running", "error", err)
			}
		}
	}
	return nil
}

func (m *Monitor) runJob(ctx context.Context, jobUUID string) error {
	// This will eventually get farmed out to a sprite registry
	spr := sprites.NewAgentSprite("bk-test-1")

	go func() {
		if err := spr.RunJob(jobUUID); err != nil {
			if err = m.finishJob(ctx, jobUUID, fmt.Sprintf("failed to run job: %s", jobUUID)); err != nil {
				log.Error("failed to finish job after run error", "error", err)
			}
		}
	}()
	return nil
}

// finishJob returns a status back to Buildkite to surface failures starting an agent
func (m *Monitor) finishJob(ctx context.Context, job string, msg string) error {
	req := stacksapi.FinishJobRequest{
		StackKey:   m.stackKey,
		JobUUID:    job,
		ExitStatus: -1,
		Detail:     msg,
	}
	_, err := m.client.FinishJob(ctx, req)
	if err != nil {
		log.Error("failed to finish the job", "error", err)
		return err
	}
	return nil
}
