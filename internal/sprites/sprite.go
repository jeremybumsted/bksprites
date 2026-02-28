// Package sprites provides an interface for invoking the
// buildkite agent on a Fly.io Sprite
package sprites

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	sprites "github.com/superfly/sprites-go"
)

const (
	spriteCommandTimeout = 5 * time.Minute
	spriteRunMaxAttempts = 3
	spriteRetryDelay     = 2 * time.Second
)

type AgentSprite struct {
	Name    string // This is the name of the sprite the agent will be run on.
	Address string // This is the ip address of the sprite
	// command sprites.Command  <- Don't know if this is useful yet.
}

func NewAgentSprite(name string) *AgentSprite {
	spriteAuthToken := os.Getenv("SPRITE_API_TOKEN")
	client := sprites.New(spriteAuthToken)
	sprite := client.Sprite(name)

	addr := sprite.URL

	return &AgentSprite{
		Name:    name,
		Address: addr,
	}
}

func (a *AgentSprite) RunJob(jobUUID string) error {
	log.Info("We'll run this job", "uuid", jobUUID)

	spriteAuthToken := os.Getenv("SPRITE_API_TOKEN")
	client := sprites.New(spriteAuthToken)
	sprite := client.Sprite(a.Name)

	var err error
	for attempt := 1; attempt <= spriteRunMaxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), spriteCommandTimeout)
		cmd := sprite.CommandContext(ctx, ".buildkite-agent/bin/buildkite-agent", "start", "--acquire-job", jobUUID)
		err = cmd.Run()
		cancel()

		if err == nil {
			return nil
		}

		if !isRetryableRunError(err) || attempt == spriteRunMaxAttempts {
			return fmt.Errorf("failed to start sprite command after %d attempt(s): %w", attempt, err)
		}

		delay := spriteRetryDelay * time.Duration(1<<(attempt-1))
		log.Warn("Sprite run attempt failed, retrying",
			"sprite", a.Name,
			"jobUUID", jobUUID,
			"attempt", attempt,
			"maxAttempts", spriteRunMaxAttempts,
			"retryIn", delay,
			"error", err,
		)
		time.Sleep(delay)
	}

	return fmt.Errorf("failed to start sprite command: %w", err)
}

func isRetryableRunError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "i/o timeout") ||
		strings.Contains(msg, "failed to connect") ||
		strings.Contains(msg, "connection reset by peer")
}
