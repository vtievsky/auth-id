package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	cacheTTL = time.Second * 60
)

var (
	ErrValueNotFound = errors.New("value not found")
)

type Cache[K comparable, V any] struct {
	m        map[K]V
	mu       sync.RWMutex
	lastTime time.Time
}

func New[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		m:        map[K]V{},
		mu:       sync.RWMutex{},
		lastTime: time.Time{},
	}
}

func (s *Cache[K, V]) Get(ctx context.Context, id K, syncFunc func(ctx context.Context) error) (V, error) {
	s.mu.RLock()

	// Кеш невалидный
	if cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := syncFunc(ctx); err != nil {
			return *new(V), err
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if val, ok := s.m[id]; ok {
		return val, nil
	}

	return *new(V), ErrValueNotFound
}

func (s *Cache[K, V]) Add(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value

	s.lastTime = time.Now()
}
