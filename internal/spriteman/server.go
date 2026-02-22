package spriteman

// spriteman aims to act as the manager of the current sprite, reporting
// capacity back to the controller so that it can intelligently assign
// jobs to a sprite in the pool.

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type Manager struct {
	maxAgents int
	interval  time.Duration
}

type JobRequest struct {
	JobUUID string `json:"job_uuid"`
}

func NewManager(maxAgents int, interval time.Duration) *Manager {
	return &Manager{
		maxAgents: maxAgents,
		interval:  interval,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	log.Info("Starting spriteman server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/job", jobHandler)
	mux.HandleFunc("/health", healthHandler)

	addr := ":8080"
	log.Info("Server started", "port", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server error", "error", err)
		return err
	}

	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func jobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var payload JobRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Info("Received job request", "jobUUID", payload.JobUUID)

	// We'll do some stuff with the request to fire up an agent here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"`))
}
