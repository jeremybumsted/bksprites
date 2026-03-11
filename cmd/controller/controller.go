// Package controller provides the kong command interface for running the controlle
package controller

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buildkite/stacksapi"
	"github.com/charmbracelet/log"

	"github.com/jeremybumsted/bksprites/internal/monitor"
)

type ControllerCmd struct {
	AgentToken   string `help:"Buildkite agent token" env:"BUILDKITE_AGENT_TOKEN" required:""`
	SpriteToken  string `help:"Sprites API token" env:"SPRITE_API_TOKEN" required:""`
	StackKey     string `help:"unique stack key" default:"bk-sprites"`
	Queue        string `help:"queue the stack will monitor" default:"default"`
	PollInterval string `help:"Poll interval" default:"1s" env:"POLL_INTERVAL"`
	LogLevel     string `help:"Log level (debug, info, warn, error)" default:"info" env:"LOG_LEVEL"`
}

func (c *ControllerCmd) Run() error {
	// Set log level
	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Warn("Invalid log level, using info", "level", c.LogLevel)
		level = log.InfoLevel
	}
	log.SetLevel(level)

	ctx := context.Background()
	log.Info("Starting controller")
	log.Info(fmt.Sprintf("Stack Key: %v", c.StackKey))
	log.Info(fmt.Sprintf("Queue: %v", c.Queue))

	// Verify sprite token is set
	if c.SpriteToken == "" {
		log.Error("SPRITE_API_TOKEN is empty - sprites authentication will fail")
		os.Exit(1)
	}
	log.Debug("Sprite token configured", "tokenLength", len(c.SpriteToken))

	client, err := stacksapi.NewClient(c.AgentToken)
	if err != nil {
		log.Error("Error creating the API client:", "error", err)
		os.Exit(1)
	}

	stack, _, err := client.RegisterStack(context.Background(), stacksapi.RegisterStackRequest{
		Key:      c.StackKey,
		Type:     stacksapi.StackTypeCustom,
		QueueKey: c.Queue,
		Metadata: map[string]string{
			"test": "true",
		},
	})
	if err != nil {
		log.Error("There was an error registering the stack", "error", err)
		os.Exit(1)
	}

	pollInterval, err := time.ParseDuration(c.PollInterval)
	if err != nil {
		return err
	}

	queueMonitor := monitor.NewMonitor(client, c.StackKey, c.Queue, pollInterval, c.SpriteToken)
	go func() {
		if err := queueMonitor.Start(ctx); err != nil && err != context.Canceled {
			log.Error("There was a monitor error", "error", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	log.Info(fmt.Sprintf("Deregistering stack %v...", stack.Key))
	_, err = client.DeregisterStack(context.Background(), stack.Key)
	if err != nil {
		log.Error("There was an error deregistering the stack", "error", err)
		os.Exit(1)
	}

	log.Info("Shutting down now, buh-bye!")

	//TODO, add the following
	// - Finish a job:
	//
	return nil
}
