// Package memcache provides an in-memory LRU cache
package memcache

import (
	"sync"

	"darvaza.org/core"
	"darvaza.org/slog"
	"github.com/darvaza-proxy/cache"
)

var (
	_ cache.Store = (*Store)(nil)
)

// Store manages in-memory [cache.Cache]s
type Store struct {
	mu  sync.Mutex
	log slog.Logger
	m   map[string]*Cache
}

// New creates a new [Store]
func New() *Store {
	s := &Store{
		m: make(map[string]*Cache),
	}
	return s
}

// DeregisterCache disconnects a [Cache] from the [Store]
func (s *Store) DeregisterCache(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.m[name]
	if ok {
		delete(s.m, name)
	}
}

// GetCache finds a [Cache] by name in the [Store]
func (s *Store) GetCache(name string) cache.Cache {
	s.mu.Lock()
	defer s.mu.Unlock()

	g, ok := s.m[name]
	if ok {
		return g
	}
	return nil
}

// NewCache creates a new in-memory [Cache] attached to the [Store]
func (s *Store) NewCache(name string, cacheBytes int64, getter cache.Getter) cache.Cache {
	if name == "" || getter == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[name]; ok {
		core.Panicf("%s: %s", name, "cache already registered")
	}

	g := NewCache(name, cacheBytes, getter)
	g.SetLogger(s.log)
	s.m[name] = g

	if log, ok := s.withDebug(); ok {
		log.Printf("cache:%q created", name)
	}
	return g
}

// SetLogger attaches a [slog.Logger] to the store and any new Cache created through it
func (s *Store) SetLogger(log slog.Logger) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.log = log
}

func (s *Store) withLogger(level slog.LogLevel) (slog.Logger, bool) {
	if s.log != nil {
		return s.log.WithLevel(level).WithEnabled()
	}
	return nil, false
}

func (s *Store) withDebug() (slog.Logger, bool) {
	return s.withLogger(slog.Debug)
}
