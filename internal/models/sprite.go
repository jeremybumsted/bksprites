package models

import "time"

// Sprite represents the configuration and other metadata
// of the sprite instance that will be used by the controller
// it is primarily used by the create command to provision a new sprite.
type Sprite struct {
	Name         string
	Organization string
	// Config is disabled for now; maybe we want a Config at some point.
	// Config	[]string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type BuildkiteAgent struct {
	Version string
	Name    string
}
