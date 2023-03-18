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
type LRU struct {
	maxSize  int
	size     int
	count    int
	items    map[string]*list.Element
	eviction *list.List
	onEvict  func(string, any, int)
}

// NewLRU creates a new LRU with maximum size and eviction callback
func NewLRU(size int, onEvict func(string, any, int)) *LRU {
	lru := &LRU{
		maxSize:  size,
		items:    make(map[string]*list.Element),
		eviction: list.New(),
		onEvict:  onEvict,
	}
	return lru
}

// Len returns the number of entries in the cache
func (m *LRU) Len() int {
	return m.count
}

// Size returns the added size of all entries in the cache
func (m *LRU) Size() int {
	return m.size
}

// Available tells how much space is available without evictions
func (m *LRU) Available() int {
	return m.maxSize - m.size
}

// Add adds an entry of a given size and optional expiration date, and
// returns true if entries were removed
func (m *LRU) Add(key string, value any, size int, expire time.Time) bool {
	var ex *time.Time
	if !expire.IsZero() {
		ex = &expire
	}

	e := entry{
		key:    key,
		value:  value,
		size:   size,
		expire: ex,
	}

	if le, ok := m.items[key]; ok {
		// update entry
		m.eviction.MoveToBack(le)

		// old
		p := le.Value.(*entry)
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
func (m *LRU) Remove(key string) {
	if le, ok := m.items[key]; ok {
		m.evictElement(le)
	}
}

// Get tries to find an entry, and returns its value and if it was found
func (m *LRU) Get(key string) (any, bool) {
	if le, ok := m.items[key]; ok {
		p := le.Value.(*entry)
		if !p.Expired() {
			m.eviction.MoveToBack(le)
			return p.value, true
		}
		m.evictElement(le)
	}
	return nil, false
}

// Prune removes entries if space is needed. It tries
// the oldests expired first, and then just the oldests.
func (m *LRU) Prune() bool {
	evicted := false

	if m.Available() < 0 {
		// evict expired first
		core.ListForEachElement(m.eviction,
			func(le *list.Element) bool {
				p := le.Value.(*entry)
				if p.Expired() {
					evicted = true
					m.evictElement(le)
				}

				return m.Available() >= 0
			})
	}

	if m.Available() < 0 {
		// evict oldest
		core.ListForEachElement(m.eviction,
			func(le *list.Element) bool {
				evicted = true
				m.evictElement(le)

				return m.Available() >= 0
			})
	}

	return evicted
}

func (m *LRU) evictElement(le *list.Element) {
	p := le.Value.(*entry)

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

type entry struct {
	key    string
	value  any
	size   int
	expire *time.Time
}

func (e *entry) Expired() bool {
	if e.expire == nil {
		return false
	}
	return time.Now().After(*e.expire)
}
