// Package models provides the type definitions
// for models used by the controller
package models

import "time"

type Job struct {
	Sprite          string
	Priority        int
	AgentQueryRules []string
	ScheduledAt     time.Time
	Pipeline        Pipeline
}

type Pipeline struct {
	Slug string
	UUID string
}

type Build struct {
	Number int
	Branch string
	UUID   string
}

type Step struct {
	Key string
}
