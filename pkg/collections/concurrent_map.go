package collections

import (
	"github.com/samber/lo"
	"sync"
)

type ConcurrentMap[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, val V)
	Values() []V
	Clone() map[K]V
}

type concurrentMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func NewConcurrentMap[K comparable, V any](size int) ConcurrentMap[K, V] {
	return &concurrentMap[K, V]{
		m: make(map[K]V, size),
	}
}

func (c *concurrentMap[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}

func (c *concurrentMap[K, V]) Set(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = val
}

func (c *concurrentMap[K, V]) Values() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return lo.Values(c.m)
}

func (c *concurrentMap[K, V]) Clone() map[K]V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// create new map with same values
	return lo.PickBy(c.m, func(key K, value V) bool {
		return true
	})
}
