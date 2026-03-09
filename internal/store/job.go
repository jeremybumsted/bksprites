package store

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/jeremybumsted/bksprites/internal/types"
)

type Job interface {
	Set(id string, j types.Job) error
	Get(id string) (types.Job, bool, error)
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

func (js *JobStore) Set(id string, j types.Job) error {
	b, err := json.Marshal(j)
	if err != nil {
		return err
	}
	log.Info("stored job", "uuid", id)
	return js.store.Set("job:"+id, string(b), 0)
}

func (js *JobStore) Get(id string) (types.Job, bool, error) {
	raw, ok := js.store.Get("job:" + id)
	if !ok {
		return types.Job{}, false, nil
	}

	var j types.Job

	if err := json.Unmarshal([]byte(raw), &j); err != nil {
		return types.Job{}, false, err
	}

	return j, true, nil
}

func (js *JobStore) Delete(id string) error {
	log.Info("Deleted job", "uuid", id)
	return js.store.Delete("job:" + id)
}
