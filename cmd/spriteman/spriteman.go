// Package spriteman provides the command interface for running the spriteman worker on the spirte instance.
package spriteman

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"

	"github.com/jeremybumsted/bksprites/internal/spriteman"
)

type SpriteManCmd struct {
	UpdateInterval string `help:"How often the manager should update the controller" default:"5s" env:"UPDATE_INTERVAL"`
	AgentLimit     int    `help:"Limit of running agent processes on the sprite" default:"4" env:"AGENT_LIMIT"`
}

func (c *SpriteManCmd) Run() error {
	ctx := context.Background()

	updateInterval, err := time.ParseDuration(c.UpdateInterval)
	if err != nil {
		return err
	}

	manager := spriteman.NewManager(c.AgentLimit, updateInterval)
	// Start the spriteman server
	go func() {
		err := manager.Start(ctx)

		if err != nil && err != context.Canceled {
			log.Error("there was an error running the server", "error", err)
		}
	}()
	// Watch for the cancel signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	return nil
}
