// Package cache provides a generic cache interface
package cache

import (
	"context"
	"time"
)

var (
	_ Getter[string] = Cache[string](nil)
	_ Setter[string] = Cache[string](nil)
)

// A Cache namespace
type Cache[K comparable] interface {
	// Stats returns stats about the Cache namespace
	Stats(Type) Stats
	// Name returns the name of the cache
	Name() string

	// Set adds an entry to the Cache
	Set(ctx context.Context, key K, value []byte, expire time.Time, cacheType Type) error
	// Get reads an entry into a Sink
	Get(ctx context.Context, key K, dest Sink) error
	// Remove removes an entry from the Cache
	Remove(ctx context.Context, key K)
}

// A Getter loads data for a key.
type Getter[K comparable] interface {
	// Get returns the value identified by key, populating dest.
	//
	// The returned data must be unversioned. That is, key must
	// uniquely describe the loaded data, without an implicit
	// current time, and without relying on cache expiration
	// mechanisms.
	Get(ctx context.Context, key K, dest Sink) error
}

// A GetterFunc implements Getter with a function.
type GetterFunc[K comparable] func(ctx context.Context, key K, dest Sink) error

// Get allows a GetterFunc to implement the Getter interface
func (f GetterFunc[K]) Get(ctx context.Context, key K, dest Sink) error {
	return f(ctx, key, dest)
}

// A Setter stores data for a key
type Setter[K comparable] interface {
	Set(ctx context.Context, key K, value []byte,
		expire time.Time, cacheType Type) error
}

// A SetterFunc implements Setter with a function
type SetterFunc[K comparable] func(ctx context.Context, key K, value []byte,
	expire time.Time, cacheType Type) error

// Set allows a SetterFunc to implement the Setter interface
func (f SetterFunc[K]) Set(ctx context.Context, key K, value []byte,
	expire time.Time, cacheType Type) error {
	//
	return f(ctx, key, value, expire, cacheType)
}
