// Package groupcache provides a groupcache backed implementation
// of darvaza.org/cache
package groupcache

import (
	"context"
	"net/http"
	"time"

	"github.com/mailgun/groupcache/v2"

	"darvaza.org/cache"
	"darvaza.org/cache/internal"
	"darvaza.org/core"
	"darvaza.org/slog"
)

var (
	_ cache.Store  = (*HTTPPool)(nil)
	_ http.Handler = (*HTTPPool)(nil)
	_ cache.Store  = (*Pool)(nil)
	_ cache.Cache  = (*Group)(nil)
)

// NewPool creates a Store placeholder to be used as entrypoint
// to a previously initialised groupcache
func NewPool() *Pool {
	return &Pool{}
}

// HTTPPool implements a cache.Store using mailgun's groupcache.HTTPPool.
type HTTPPool struct {
	Pool

	pool *groupcache.HTTPPool
}

// ServeHTTP handles the BasePath of the cache. "/_groupcache/" if unspecified.
func (p *HTTPPool) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.pool.ServeHTTP(rw, req)
}

// NewHTTPPoolOpts initializes an HTTP pool of peers with the given options.
// Unlike NewHTTPPool, this function does not register the created pool as
// an HTTP handler.
//
// Due to groupcache's own limitations you have to choose to use either
// HTTPPool or NoPeersPool. There can be only one.
func NewHTTPPoolOpts(self string, opts *groupcache.HTTPPoolOptions) *HTTPPool {
	if pool := groupcache.NewHTTPPoolOpts(self, opts); pool != nil {
		return &HTTPPool{
			pool: pool,
		}
	}
	return nil
}

// NewHTTPPool initialises an HTTP pool of peers, and registers itself as a PeerPicker
// Due to groupcache's own limitations you have to choose to use either
// HTTPPool or NoPeersPool. There can be only one.
func NewHTTPPool(self string) *HTTPPool {
	return NewHTTPPoolOpts(self, nil)
}

// Pool implements a cache.Store using mailgun's groupcache.
// Groupcache is global so the Pool object doesn't contain anything
type Pool struct{}

// NewNoPeersPool initialises groupcache with a PeerPicker that never find a peer
// Due to groupcache's own limitations you have to choose to use either
// HTTPPool or NoPeersPool. There can be only one.
func NewNoPeersPool() *Pool {
	pool := &groupcache.NoPeers{}
	groupcache.RegisterPeerPicker(func() groupcache.PeerPicker {
		return pool
	})
	return &Pool{}
}

// NewCache creates a new Group
func (p *Pool) NewCache(name string, cacheBytes int64, getter cache.Getter) cache.Cache {
	// wrap
	fn := func(ctx context.Context, key string, dst groupcache.Sink) error {
		sink, ok := internal.Sink(ctx)
		if ok {
			return p.getBridged(ctx, key, getter, sink, dst)
		}

		// unreachable
		core.Panicf("groupcache[%q]: %s", name, "where is my sink??")
		return cache.ErrInvalid
	}

	if g := groupcache.NewGroup(name, cacheBytes,
		groupcache.GetterFunc(fn)); g != nil {
		return &Group{g}
	}
	return nil
}

// getBridged calls the getter using the Sink in the context, and then copies the
// binary encoded representation in groupcache's Sink so it gets cached
func (*Pool) getBridged(ctx context.Context, key string,
	getter cache.Getter, sink cache.Sink, dst groupcache.Sink) error {
	//
	err := getter.Get(ctx, key, sink)
	if err == nil {
		// good. store it in dst for the cache
		err = dst.SetBytes(sink.Bytes(), sink.Expire())
	}
	return err
}

// GetCache returns a named Group previously created
func (*Pool) GetCache(name string) cache.Cache {
	if g := groupcache.GetGroup(name); g != nil {
		return &Group{g}
	}
	return nil
}

// DeregisterCache removes a Group from the Pool
func (*Pool) DeregisterCache(name string) {
	groupcache.DeregisterGroup(name)
}

// SetLogger attaches a slog.Logger to groupcache
func (*Pool) SetLogger(log slog.Logger) {
	SetLogger(log)
}

// Group implements cache.Cache around a groupcache.Group
type Group struct {
	g *groupcache.Group
}

// Name returns the name of the Group
func (g *Group) Name() string {
	return g.g.Name()
}

// Set adds an entry to the Group
func (g *Group) Set(ctx context.Context, key string, value []byte,
	expire time.Time, t cache.Type) error {
	//
	var hot bool
	if t == cache.HotCache {
		hot = true
	}

	return g.g.Set(ctx, key, value, expire, hot)
}

// Get reads an entry into a Sink
func (g *Group) Get(ctx context.Context, key string, sink cache.Sink) error {
	var b groupcache.ByteView

	if sink == nil {
		return cache.ErrInvalid
	}

	// get a clean slate and attach sink to the context
	sink.Reset()
	ctx = internal.WithSink(ctx, sink)

	s := groupcache.ByteViewSink(&b)
	err := g.g.Get(ctx, key, s)
	if err != nil {
		return err
	}

	if sink.Len() > 0 {
		// the sink we passed via the context was used.
		// we are ready.
		return nil
	}

	// it was a cache hit, so we read the encoded data
	// for the ByteView we created
	e := b.Expire()
	bytes := b.ByteSlice()
	return sink.SetBytes(bytes, e)
}

// Remove removes an entry from the Group
func (g *Group) Remove(ctx context.Context, key string) {
	g.g.Remove(ctx, key)
}

// Stats returns stats about the Group
func (g *Group) Stats(t cache.Type) cache.Stats {
	var t1 groupcache.CacheType
	switch t {
	case cache.MainCache:
		t1 = groupcache.MainCache
	case cache.HotCache:
		t1 = groupcache.HotCache
	default:
		// invalid type
		return cache.Stats{}
	}

	s1 := g.g.CacheStats(t1)
	s := cache.Stats{
		Bytes:     s1.Bytes,
		Items:     s1.Items,
		Gets:      s1.Gets,
		Hits:      s1.Hits,
		Evictions: s1.Evictions,
	}
	return s
}
