package models

import "time"

// AgentSprite represents the configuration and other metadata
// of the sprite instance that will be used by the controller
// it is primarily used by the create command to provision a new sprite.
type AgentSprite struct {
	Name         string
	Organization string
	ConfigFile   string
	MaxAgents    int
	MinAgents    int

	Agent BuildkiteAgent

	CreatedAt time.Time
	UpdatedAt time.Time
}

type BuildkiteAgent struct {
	Version string
	Name    string
}
