// Package simplelru provices a non thread-safe LRU cache with maximum size
// and TTL
package simplelru

import (
	"container/list"
	"time"

	"darvaza.org/core"
)

// LRU Implements a least-recently-used cache with a maximum size and
// optional expiration date
type LRU[K comparable, T any] struct {
	maxSize  int
	size     int
	count    int
	items    map[K]*list.Element
	eviction *list.List
	onAdd    func(K, T, int, time.Time)
	onEvict  func(K, T, int)
}

// NewLRU creates a new LRU with maximum size and eviction callback
func NewLRU[K comparable, T any](size int,
	onAdd func(K, T, int, time.Time),
	onEvict func(K, T, int)) *LRU[K, T] {
	//
	lru := &LRU[K, T]{
		maxSize:  size,
		items:    make(map[K]*list.Element),
		eviction: list.New(),
		onAdd:    onAdd,
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

func (m *LRU[K, T]) getEntry(key K) (*list.Element, *entry[K, T], bool) {
	le, ok := m.items[key]
	if !ok {
		return nil, nil, false
	}

	p, ok := m.getListEntry(le)
	if !ok {
		// not possible
		delete(m.items, key)
		return nil, nil, false
	}

	return le, p, true
}

func (*LRU[K, T]) getListEntry(le *list.Element) (*entry[K, T], bool) {
	p, ok := le.Value.(*entry[K, T])
	return p, ok
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

	if le, p, ok := m.getEntry(key); ok {
		// update entry
		m.eviction.MoveToBack(le)

		// old
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
	evicted := m.prune()

	if m.onAdd != nil {
		// notify the user
		m.onAdd(key, value, size, expire)
	}

	return evicted
}

// Evict removes an entry if present
func (m *LRU[K, T]) Evict(key K) {
	if le, ok := m.items[key]; ok {
		m.evictElement(le)
	}
}

// Get tries to find an entry, and returns its value, expiration date,
// and if it was found
func (m *LRU[K, T]) Get(key K) (T, time.Time, bool) {
	var zero T
	if le, p, ok := m.getEntry(key); ok {
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

// prune removes entries if space is needed. It tries
// the oldest expired first, and then just the oldest.
func (m *LRU[K, T]) prune() bool {
	var evicted bool

	if m.needsPruning() {
		// evict expired first
		if m.pruneLoop(true) {
			evicted = true
		}
	}

	if m.needsPruning() {
		// evict oldest
		if m.pruneLoop(false) {
			evicted = true
		}
	}

	return evicted
}

// revive:disable:flag-parameter
func (m *LRU[K, T]) pruneLoop(onlyExpired bool) bool {
	// revive:enable:flag-parameter
	var evicted bool

	core.ListForEachElement(m.eviction,
		func(le *list.Element) bool {
			if m.pruneEvict(le, onlyExpired) {
				evicted = true
			}
			return !m.needsPruning()
		})

	return evicted
}

func (m *LRU[K, T]) pruneEvict(le *list.Element, onlyExpired bool) bool {
	var evict bool
	p, ok := m.getListEntry(le)
	switch {
	case !ok:
		evict = true
	case !onlyExpired:
		evict = true
	default:
		evict = p.Expired()
	}

	if evict {
		m.evictElement(le)
	}

	return evict
}

// EvictExpired scans the whole cache and evicts all expired entries
func (m *LRU[K, T]) EvictExpired() bool {
	var evicted bool
	core.ListForEachElement(m.eviction,
		func(le *list.Element) bool {
			p, ok := m.getListEntry(le)
			if !ok || p.Expired() {
				evicted = true
				m.evictElement(le)
			}

			return false
		})
	return evicted
}

func (m *LRU[K, T]) evictElement(le *list.Element) {
	p, ok := m.getListEntry(le)

	// remove from eviction list
	m.eviction.Remove(le)

	if !ok {
		return
	}

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

// ForEach allows you to iterate over all non-expired entries in the Cache.
// ForEach will automatically evict any expired entry it finds along the way,
// and it will stop if the callback returns true.
func (m *LRU[K, T]) ForEach(fn func(K, T, int, time.Time) bool) {
	if fn != nil {
		core.ListForEachElement(m.eviction, func(le *list.Element) bool {
			return m.forEachIter(le, fn)
		})
	}
}

func (m *LRU[K, T]) forEachIter(le *list.Element, fn func(K, T, int, time.Time) bool) bool {
	var ex time.Time

	p, ok := m.getListEntry(le)
	if !ok || p.Expired() {
		m.evictElement(le)
		return false
	}

	if p.expire != nil {
		ex = *p.expire
	}

	return fn(p.key, p.value, p.size, ex)
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
