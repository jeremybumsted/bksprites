package monitor

// We need to make sure that a sprite is warmed, able to accept jobs,
// and otherwise healthy, so we have a few methods to do this.

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

// CheckSpriteHealth checks to ensure the sprite is active and the server is running
// this has the added benefit of making the sprite warm if it is currently
// cold for autoscaling purposes
func (m *Monitor) CheckSpriteHealth(addr string, port string) (bool, error) {
	if port == "" {
		port = "8080"
	}
	url := fmt.Sprintf("http://%v:%v/health", addr, port)

	resp, err := http.Get(url)
	if err != nil {
		log.Error("Error creating the request", "error", err)
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := errors.New("sprite did not return 200, unhealthy")
		return false, err
	}

	return true, nil
}
