// Package groupcache provides a groupcache backed implementation
// of github.com/darvaza-proxy/cache
package groupcache

import (
	"context"
	"time"

	"github.com/mailgun/groupcache/v2"

	"github.com/darvaza-proxy/cache"
	"github.com/darvaza-proxy/slog"
)

var (
	_ cache.Store = (*Pool)(nil)
	_ cache.Cache = (*Group)(nil)
)

// Pool implements a cache.Store using mailgun's groupcache
type Pool struct{}

// NewCache creates a new Group
func (*Pool) NewCache(string, int64, cache.Getter) cache.Cache {
	return nil
}

// GetCache returns a named Group previously created
func (*Pool) GetCache(string) cache.Cache {
	return nil
}

// DeregisterCache removes a Group from the Pool
func (*Pool) DeregisterCache(string) {}

// SetLogger attaches a slog.Logger to groupcache
func (*Pool) SetLogger(slog.Logger) {}

// Group implements cache.Cache around a groupcache.Group
type Group struct {
	g *groupcache.Group
}

// Name returns the name of the Group
func (g *Group) Name() string {
	return g.g.Name()
}

// Set adds an entry to the Group
func (*Group) Set(context.Context, string, []byte, time.Time, cache.Type) error {
	return nil
}

// Get reads an entry into a Sink
func (*Group) Get(context.Context, string, cache.Sink) error {
	return nil
}

// Remove removes an entry from the Group
func (*Group) Remove(context.Context, string) {}

// Stats returns stats about the Group
func (*Group) Stats(cache.Type) cache.Stats {
	return cache.Stats{}
}
