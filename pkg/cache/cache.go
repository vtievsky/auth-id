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
	mu       sync.Mutex
	lastTime time.Time
}

func New[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		m:        map[K]V{},
		mu:       sync.Mutex{},
		lastTime: time.Time{},
	}
}

func (s *Cache[K, V]) Get(ctx context.Context, key K, cacheSyncFunc func(ctx context.Context) (map[K]V, error)) (V, error) {
	copyValues := func(values map[K]V) {
		for key, value := range values {
			s.m[key] = value
		}

		s.lastTime = time.Now()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		ok    bool
		value V
	)

	// Синхронизация кеша с истекшим сроком годности
	if cacheTTL < time.Since(s.lastTime) {
		values, err := cacheSyncFunc(ctx)
		if err != nil {
			return value, err
		}

		copyValues(values)
	}

	// Поиск значения в актуальном кеше
	if value, ok = s.m[key]; ok {
		return value, nil
	}

	// Синхронизация кеша в случае, когда кеш еще актуален,
	// а значение могло быть добавлено в хранилище другим экземпляром приложения
	values, err := cacheSyncFunc(ctx)
	if err != nil {
		return value, err
	}

	copyValues(values)

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

func (s *Cache[K, V]) Del(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, key)
}
