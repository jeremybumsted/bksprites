package store

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/jeremybumsted/bksprites/internal/models"
)

type Job interface {
	Set(id string, j models.Job) error
	Get(id string) (models.Job, bool, error)
	Delete(id string) error
}

type JobStore struct {
	store *Store
}

func NewJobStore(store *Store) *JobStore {
	return &JobStore{
		store: store,
	}
}

func (js *JobStore) Set(id string, j models.Job) error {
	b, err := json.Marshal(j)
	if err != nil {
		return err
	}
	log.Info("stored job", "uuid", id)
	return js.store.Set("job:"+id, string(b), 0)
}

func (js *JobStore) Get(id string) (models.Job, bool, error) {
	raw, ok := js.store.Get("job:" + id)
	if !ok {
		return models.Job{}, false, nil
	}

	var j models.Job

	if err := json.Unmarshal([]byte(raw), &j); err != nil {
		return models.Job{}, false, err
	}

	return j, true, nil
}

func (js *JobStore) Delete(id string) error {
	log.Info("Deleted job", "uuid", id)
	return js.store.Delete("job:" + id)
}
