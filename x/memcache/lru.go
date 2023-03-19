package memcache

import (
	"context"
	"sync"
	"time"

	"github.com/darvaza-proxy/cache"
	"github.com/darvaza-proxy/cache/x/simplelru"
)

const (
	KiB = 1024      // KiB is 2^10 (kilobyte)
	MiB = KiB * KiB // MiB is 2^20 (megabyte)
	GiB = KiB * MiB // GiB is 2^30 (gigabyte)
)

// LRU is a least-recently-used cache of bytes with TTL and maximum size
type LRU[K comparable] struct {
	mu      sync.Mutex
	lru     *simplelru.LRU[K, []byte]
	unit    uint
	onEvict func(K, []byte, int64)
	stats   cache.Stats
}

// NewLRU creates a new []byte [LRU] with maximum size and eviction
func NewLRU[K comparable](cacheBytes int64, onEvict func(K, []byte, int64)) *LRU[K] {
	unit := calculateUnit(cacheBytes)
	size := bytesToSize(unit, cacheBytes)

	m := &LRU[K]{
		unit:    unit,
		onEvict: onEvict,
	}

	lru := simplelru.NewLRU(size, nil, m.evictionCallback)
	m.lru = lru

	return m
}

func (m *LRU[K]) evictionCallback(key K, value []byte, size int) {
	// increment evictions count
	m.stats.Evictions++

	if m.onEvict != nil {
		// and inform the user
		m.onEvict(key, value, m.fromUnit(size))
	}
}

// Items returns the number of entries in the cache
func (m *LRU[K]) Items() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.lru.Len()
}

// Size returns the added size if bytes of all entries in the cache
func (m *LRU[K]) Size() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.fromUnit(m.lru.Size())
}

// Stats returns statistics about the [LRU]
func (m *LRU[K]) Stats() cache.Stats {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := m.stats
	stats.Bytes = m.fromUnit(m.lru.Size())
	stats.Items = int64(m.lru.Len())
	return stats
}

// Add adds an entry and cache duration, and returns true if entries were removed
// to free capacity. if expire is 0, it never expires.
func (m *LRU[K]) Add(key K, value []byte, expire time.Time) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := m.toUnit(len64(value))
	return m.lru.Add(key, value, size, expire)
}

// Evict removes an entry if present
func (m *LRU[K]) Evict(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lru.Evict(key)
}

// Get attempts to find an entry in the cache, and returns its value,
// expiration date if any, and if it was found or not
func (m *LRU[K]) Get(key K) ([]byte, *time.Time, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats.Gets++
	v, ex, ok := m.lru.Get(key)
	if ok {
		m.stats.Hits++
	}

	if ex.IsZero() {
		return v, nil, ok
	}
	return v, &ex, ok
}

// EvictExpired periodically scans for expired entries and evicts them from the cache.
// It runs until the provided context is cancelled.
func (m *LRU[K]) EvictExpired(ctx context.Context, period time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(period):
			m.mu.Lock()
			m.lru.EvictExpired()
			m.mu.Unlock()
		}
	}
}

func (m *LRU[K]) fromUnit(size int) int64 {
	return sizeToBytes(m.unit, size)
}

func (m *LRU[K]) toUnit(bytes int64) int {
	return bytesToSize(m.unit, bytes)
}

func calculateUnit(max int64) (unit uint) {
	unit = 1
	for max > GiB {
		unit *= KiB
		max /= KiB
	}
	return unit
}

func bytesToSize(unit uint, bytesCount int64) (size int) {
	if bytesCount < 1 {
		return 0
	}

	n := int64(unit)
	size = int(bytesCount / n)
	if (bytesCount % n) > 0 {
		size++
	}

	return size
}

func sizeToBytes(unit uint, size int) int64 {
	if size == 0 {
		return 0
	}

	return int64(unit) * int64(size)
}

func len64(b []byte) int64 {
	return int64(len(b))
}
