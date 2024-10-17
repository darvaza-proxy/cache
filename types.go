package cache

import (
	"errors"

	"darvaza.org/slog"
)

// ErrInvalidSink tells the [Sink] isn't in a usable state.
var ErrInvalidSink = errors.New("invalid sink")

// A Store allows us to create or access Cache namespaces
type Store interface {
	// GetCache returns the named cache previously created with
	// NewCache, or nil if there's no such namespace.
	GetCache(name string) Cache[string]
	// NewCache creates a new Cache namespace
	NewCache(name string, cacheBytes int64, getter Getter[string]) Cache[string]
	// DeregisterCache removes a Cache namespace
	DeregisterCache(name string)

	// SetLogger binds the Store and its Cache namespaces to a logger
	SetLogger(log slog.Logger)
}

// Stats provides an snapshot on the state of a Cache namespace
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
	// MainCache is the cache for items we own
	MainCache Type = iota + 1
	// HotCache is the cache for items that seem popular
	// even if we don't necessarily own
	HotCache
)
