package store

import (
	"encoding/json"

	"github.com/jeremybumsted/bksprites/internal/models"
)

type bkJob models.Job

type Job interface {
	Set(id string, j bkJob) error
	Get(id string) (bkJob, bool, error)
	Delete(id string) error
}

type jobStore struct {
	store *Store
}

func NewJobStore(store *Store) *jobStore {
	return &jobStore{
		store: store,
	}
}

func (js *jobStore) Set(id string, j bkJob) error {
	b, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return js.store.Set("job:"+id, string(b), 0)
}

func (js *jobStore) Get(id string) (bkJob, bool, error) {
	raw, ok := js.store.Get("job:" + id)
	if !ok {
		return bkJob{}, false, nil
	}

	var j bkJob

	if err := json.Unmarshal([]byte(raw), &j); err != nil {
		return bkJob{}, false, err
	}

	return j, true, nil
}

func (js *jobStore) Delete(id string) error {
	return js.store.Delete("job:" + id)
}
