// Package fscache implements a cache over fs.FS
package fscache

import (
	"context"
	"sync"
	"time"

	"darvaza.org/cache"
	"darvaza.org/slog"
)

var (
	_ cache.Store = (*Store)(nil)
	_ cache.Cache = (*Cache)(nil)
)

type Store struct {
	mu  sync.Mutex
	log slog.Logger
	m   map[string]*Cache
}

func (s *Store) SetLogger(log slog.Logger) {}

func (s *Store) DeregisterCache(name string) {}

func (s *Store) GetCache(name string) cache.Cache {}

func (s *Store) NewCache(name string, cacheBytes int64, getter cache.Getter) cache.Cache {}

type Cache struct {
	mu    sync.Mutex
	name  string
	store *Store
}

func (g *Cache) Name() string { return g.name }
func (g *Cache) Stats(t cache.Type) cache.Stats

func (g *Cache) Get(ctx context.Context, key string, dest cache.Sink) error {}
func (g *Cache) Set(ctx context.Context, key string, data []byte, expire time.Time, cacheType cache.Type) error {
}
func (g *Cache) Remove(ctx context.Context, key string) {}
