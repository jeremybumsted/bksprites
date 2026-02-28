// Package store provides the interface for storing
// Job and SpriteData for the controller
package store

import (
	"errors"
	"sync"
	"time"
)

var ErrStoreFull = errors.New("storage full")

type Store struct {
	mu      sync.RWMutex
	data    map[string]entry
	maxKeys int
	ttl     time.Duration
}

type entry struct {
	value     string
	expiresAt time.Time
}

func NewStore() *Store {
	maxAge, _ := time.ParseDuration("10m")

	return &Store{
		data:    make(map[string]entry),
		maxKeys: 1000,
		ttl:     maxAge,
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.maxKeys > 0 && len(s.data) >= s.maxKeys {
		if _, exists := s.data[key]; !exists {
			return ErrStoreFull
		}
	}

	e := entry{value: value}
	if ttl > 0 {
		e.expiresAt = time.Now().Add(ttl)
	}

	s.data[key] = e
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var e entry

	e, ok := s.data[key]
	if !ok {
		return "", false
	}

	if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
		return "", false
	}

	return e.value, true
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return nil
}
