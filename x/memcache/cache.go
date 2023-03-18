package memcache

import (
	"context"

	"github.com/darvaza-proxy/cache"
)

var (
	_ cache.Cache = (*Cache)(nil)
)

// Cache is a LRU with TTL [cache.Cache]
type Cache struct {
	*SingleFlight

	lru *LRU[string]
}

// NewCache creates a new [Cache] with a maximum size and [cache.Getter]
func NewCache(name string, cacheBytes int64, getter cache.Getter) *Cache {
	g := &Cache{}
	lru := NewLRU(cacheBytes, nil, g.onEvict)

	g.lru = lru
	g.SingleFlight = NewSingleFlight(name, lru, getter)

	return g
}

func (g *Cache) onEvict(key string, _ []byte, size int64) {
	if log, ok := g.withDebug(); ok {
		log.WithField("key", key).
			WithField("size", size).
			Print("removed")
	}
}

// Stats returns statistics about the Cache. This implementation doesn't
// distinuish among types
func (g *Cache) Stats(_ cache.Type) cache.Stats {
	return g.lru.Stats()
}

// Remove evicts an entry from the [Cache]
func (g *Cache) Remove(_ context.Context, key string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.lru.Evict(key)
}
