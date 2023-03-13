// Package cache provides a generic cache interface
package cache

import (
	"context"
	"time"

	"github.com/darvaza-proxy/slog"
)

// A Store allows us to create or access Cache namespaces
type Store interface {
	// GetCache returns the named cache previously created with
	// NewCache, or nil if there's no such namespace.
	GetCache(name string) Cache
	// NewCache creates a new Cache namespace
	NewCache(name string, cacheBytes int64, getter Getter) Cache
	// DeregisterCache removes a Cache namespace
	DeregisterCache(name string)

	// SetLogger binds the Store and its Cache namespaces to a logger
	SetLogger(log slog.Logger)
}

// A Cache namespace
type Cache interface {
	// Stats returns stats about the Cache namespace
	Stats(Type) Stats
	// Name returns the name of the cache
	Name() string

	// Set adds an entry to the Cache
	Set(ctx context.Context, key string, value []byte, expire time.Time, cacheType Type) error
	// Get reads an entry into a Sink
	Get(ctx context.Context, key string, dest Sink) error
	// Remove removes an entry from the Cache
	Remove(ctx context.Context, key string)
}

// Stats provides an snapshot on the state of a Cache namespace
type Stats struct {
	Bytes     int64
	Items     int64
	Gets      int64
	Hits      int64
	Evictions int64
}

// Type represetns a type of cache
type Type int

const (
	// MainCache is the cache for items we own
	MainCache Type = iota + 1
	// HotCache is the cache for items that seem popular
	// even if we don't necessarily own
	HotCache
)

// Sink receives data from a Get call
type Sink interface {
	// SetString sets the value to s.
	SetString(s string, e time.Time) error

	// SetBytes sets the value to the contents of v.
	// The caller retains ownership of v.
	SetBytes(v []byte, e time.Time) error

	// SetValue sets the value to the object v.
	// The caller retains ownership of v.
	SetValue(v any, e time.Time) error

	// Bytes returns the value encoded as a slice
	// of bytes
	Bytes() []byte

	// Expire returns the time whe this entry will
	// be evicted from the Cache
	Expire() time.Time
}

// A Getter loads data for a key.
type Getter interface {
	// Get returns the value identified by key, populating dest.
	//
	// The returned data must be unversioned. That is, key must
	// uniquely describe the loaded data, without an implicit
	// current time, and without relying on cache expiration
	// mechanisms.
	Get(ctx context.Context, key string, dest Sink) error
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(ctx context.Context, key string, dest Sink) error

// Get allows a GetterFunc to implement the Getter interface
func (f GetterFunc) Get(ctx context.Context, key string, dest Sink) error {
	return f(ctx, key, dest)
}
