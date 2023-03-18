// Package simplelru provices a non thread-safe LRU cache with maximum size
// and TTL
package simplelru

import (
	"container/list"
	"time"

	"github.com/darvaza-proxy/core"
)

// LRU Implements a least-recently-used cache with a maximum size and
// optional expiration date
type LRU[K comparable, T any] struct {
	maxSize  int
	size     int
	count    int
	items    map[K]*list.Element
	eviction *list.List
	onEvict  func(K, T, int)
}

// NewLRU creates a new LRU with maximum size and eviction callback
func NewLRU[K comparable, T any](size int, onEvict func(K, T, int)) *LRU[K, T] {
	lru := &LRU[K, T]{
		maxSize:  size,
		items:    make(map[K]*list.Element),
		eviction: list.New(),
		onEvict:  onEvict,
	}
	return lru
}

// Len returns the number of entries in the cache
func (m *LRU[K, T]) Len() int {
	return m.count
}

// Size returns the added size of all entries in the cache
func (m *LRU[K, T]) Size() int {
	return m.size
}

// Available tells how much space is available without evictions
func (m *LRU[K, T]) Available() int {
	return m.maxSize - m.size
}

func (m *LRU[K, T]) needsPruning() bool {
	return m.size > m.maxSize
}

// Add adds an entry of a given size and optional expiration date, and
// returns true if entries were removed
func (m *LRU[K, T]) Add(key K, value T, size int, expire time.Time) bool {
	var ex *time.Time
	if !expire.IsZero() {
		ex = &expire
	}

	e := entry[K, T]{
		key:    key,
		value:  value,
		size:   size,
		expire: ex,
	}

	if le, ok := m.items[key]; ok {
		// update entry
		m.eviction.MoveToBack(le)

		// old
		p := le.Value.(*entry[K, T])
		m.size -= p.size
		// new
		*p = e
		m.size += e.size
	} else {
		// new entry
		le := m.eviction.PushBack(&e)
		m.items[key] = le
		m.size += e.size
		m.count++
	}

	// evict entries if needed
	return m.Prune()
}

// Remove removes an entry if present
func (m *LRU[K, T]) Remove(key K) {
	if le, ok := m.items[key]; ok {
		m.evictElement(le)
	}
}

// Get tries to find an entry, and returns its value and if it was found
func (m *LRU[K, T]) Get(key K) (T, bool) {
	v, _, ok := m.GetWithExpire(key)
	return v, ok
}

// GetWithExpire tries to find an entry, and returns its value, expiration date,
// and if it was found
func (m *LRU[K, T]) GetWithExpire(key K) (T, time.Time, bool) {
	var zero T
	if le, ok := m.items[key]; ok {
		p := le.Value.(*entry[K, T])
		if !p.Expired() {
			var e time.Time

			m.eviction.MoveToBack(le)
			if ex := p.expire; ex != nil {
				e = *ex
			}

			return p.value, e, true
		}
		m.evictElement(le)
	}
	return zero, time.Time{}, false
}

// Prune removes entries if space is needed. It tries
// the oldests expired first, and then just the oldests.
func (m *LRU[K, T]) Prune() bool {
	evicted := false

	if m.needsPruning() {
		// evict expired first
		core.ListForEachElement(m.eviction,
			func(le *list.Element) bool {
				p := le.Value.(*entry[K, T])
				if p.Expired() {
					evicted = true
					m.evictElement(le)
				}

				return !m.needsPruning()
			})
	}

	if m.needsPruning() {
		// evict oldest
		core.ListForEachElement(m.eviction,
			func(le *list.Element) bool {
				evicted = true
				m.evictElement(le)

				return !m.needsPruning()
			})
	}

	return evicted
}

// EvictExpired scans the whole cache and evicts all expired entries
func (m *LRU[K, T]) EvictExpired() bool {
	var evicted bool
	core.ListForEachElement(m.eviction,
		func(le *list.Element) bool {
			p := le.Value.(*entry[K, T])
			if p.Expired() {
				evicted = true
				m.evictElement(le)
			}

			return false
		})
	return evicted
}

func (m *LRU[K, T]) evictElement(le *list.Element) {
	p := le.Value.(*entry[K, T])

	// remove from eviction list
	m.eviction.Remove(le)
	// remove from items
	delete(m.items, p.key)
	// remove from size
	m.size -= p.size
	// remove from count
	m.count--

	// notify user
	if fn := m.onEvict; fn != nil {
		fn(p.key, p.value, p.size)
	}
}

type entry[K comparable, T any] struct {
	key    K
	value  T
	size   int
	expire *time.Time
}

func (e *entry[K, T]) Expired() bool {
	if e.expire == nil {
		return false
	}
	return time.Now().After(*e.expire)
}
