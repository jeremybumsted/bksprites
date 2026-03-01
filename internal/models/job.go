// Package models provides the type definitions
// for models used by the controller
package models

import "time"

type Job struct {
	Sprite          string    `json:"sprite,omitempty"`
	Priority        int       `json:"priority"`
	AgentQueryRules []string  `json:"agent_query_rules"`
	ScheduledAt     time.Time `json:"scheduled_at"`
	Pipeline        Pipeline  `json:"pipeline"`
	Build           Build     `json:"build"`
	Step            Step      `json:"step"`
}

type Pipeline struct {
	Slug string `json:"slug"`
	UUID string `json:"uuid"`
}

type Build struct {
	Number int    `json:"number"`
	Branch string `json:"branch"`
	UUID   string `json:"uuid"`
}

type Step struct {
	Key string `json:"key"`
}
