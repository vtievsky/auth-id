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

func (s *Cache[K, V]) Get(ctx context.Context, key K, syncFunc func(ctx context.Context) error) (V, error) {
	s.mu.RLock()

	var (
		ok    bool
		value V
	)

	if value, ok = s.m[key]; ok {
		s.mu.RLocker().Unlock()

		return value, nil
	}

	// Выполним синхронизацию, если значение отсутствует или кеш стал неактуальным
	// Проверка на отсутствие значения, в этом условии, необходимо в случае, когда
	// значение было создано в другом экземпляре приложения, а запрашивается здесь
	if !ok || cacheTTL < time.Since(s.lastTime) {
		s.mu.RLocker().Unlock()

		if err := syncFunc(ctx); err != nil {
			return value, err
		}

		s.mu.RLock()
	}

	defer s.mu.RLocker().Unlock()

	if value, ok = s.m[key]; ok {
		return value, nil
	}

	return value, ErrValueNotFound
}

func (s *Cache[K, V]) Add(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value

	s.lastTime = time.Now()
}
