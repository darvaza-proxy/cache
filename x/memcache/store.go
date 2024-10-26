// Package memcache provides an in-memory LRU cache
package memcache

import (
	"sync"

	"darvaza.org/cache"
	"darvaza.org/core"
	"darvaza.org/slog"
)

var (
	_ cache.Store[string]   = (*Store[string])(nil)
	_ cache.Store[uint32]   = (*Store[uint32])(nil)
	_ cache.Store[[32]byte] = (*Store[[32]byte])(nil)
)

// Store manages in-memory [cache.Cache]s
type Store[K comparable] struct {
	mu  sync.Mutex
	log slog.Logger
	m   map[string]*Cache[K]
}

// New creates a new [Store]
func New[K comparable]() *Store[K] {
	s := &Store[K]{
		m: make(map[string]*Cache[K]),
	}
	return s
}

// DeregisterCache disconnects a [Cache] from the [Store]
func (s *Store[K]) DeregisterCache(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.m[name]
	if ok {
		delete(s.m, name)
	}
}

// GetCache finds a [Cache] by name in the [Store]
func (s *Store[K]) GetCache(name string) cache.Cache[K] {
	s.mu.Lock()
	defer s.mu.Unlock()

	g, ok := s.m[name]
	if ok {
		return g
	}
	return nil
}

// NewCache creates a new in-memory [Cache] attached to the [Store]
func (s *Store[K]) NewCache(name string, cacheBytes int64, getter cache.Getter[K]) cache.Cache[K] {
	if name == "" || getter == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[name]; ok {
		core.Panicf("%s: %s", name, "cache already registered")
	}

	g := NewCache[K](name, cacheBytes, getter)
	g.SetLogger(s.log)
	s.m[name] = g

	if log, ok := s.withDebug(); ok {
		log.Printf("cache:%q created", name)
	}
	return g
}

// SetLogger attaches a [slog.Logger] to the store and any new Cache created through it
func (s *Store[K]) SetLogger(log slog.Logger) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.log = log
}

func (s *Store[K]) withLogger(level slog.LogLevel) (slog.Logger, bool) {
	if s.log != nil {
		return s.log.WithLevel(level).WithEnabled()
	}
	return nil, false
}

func (s *Store[K]) withDebug() (slog.Logger, bool) {
	return s.withLogger(slog.Debug)
}
