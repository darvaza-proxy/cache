package cache

import "darvaza.org/slog"

// A Store allows us to create or access Cache namespaces
type Store interface {
	// GetCache returns the named cache previously created with NewCache,
	// or nil if there's no such namespace.
	GetCache(name string) Cache
	// NewCache creates a new Cache namespace
	NewCache(name string, cacheBytes int64, getter Getter) Cache
	// DeregisterCache removes a Cache namespace
	DeregisterCache(name string)

	// SetLogger binds the Store and its Cache namespaces to a logger
	SetLogger(log slog.Logger)
}

// Stats provides a snapshot on the state of a Cache namespace
type Stats struct {
	Bytes     int64
	Items     int64
	Gets      int64
	Hits      int64
	Evictions int64
}

// Type represents a type of cache
type Type int

const (
	// MainCache represents the cache for items we owner by the application.
	MainCache Type = iota + 1
	// HotCache represents the cache for frequently accessed items, even if
	// not owned by the application.
	HotCache
)
