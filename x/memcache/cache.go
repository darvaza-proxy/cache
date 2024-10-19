package memcache

import (
	"context"

	"darvaza.org/cache"
)

var (
	_ cache.Cache[string] = (*Cache[string])(nil)
)

// Cache is a LRU with TTL [cache.Cache]
type Cache[K comparable] struct {
	*SingleFlight[K]

	lru *LRU[K]
}

// NewCache creates a new [Cache] with a maximum size and [cache.Getter]
func NewCache[K comparable](name string, cacheBytes int64, getter cache.Getter[K]) *Cache[K] {
	g := &Cache[K]{}
	lru := NewLRU(cacheBytes, nil, g.onEvict)

	g.lru = lru
	g.SingleFlight = NewSingleFlight[K](name, lru, getter)

	return g
}

func (g *Cache[K]) onEvict(key K, _ []byte, size int64) {
	if log, ok := g.withDebug(); ok {
		log.WithField("key", key).
			WithField("size", size).
			Print("removed")
	}
}

// Stats returns statistics about the Cache. This implementation doesn't
// distinguish among types
func (g *Cache[K]) Stats(_ cache.Type) cache.Stats {
	return g.lru.Stats()
}

// Remove evicts an entry from the [Cache]
func (g *Cache[K]) Remove(_ context.Context, key K) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.lru.Evict(key)
}
