// Package internal provides helpers to implement handlers
package internal

import (
	"context"

	"github.com/darvaza-proxy/cache"
)

type contextKey struct {
	name string
}

// WithSink attaches a cache.Sink to the context
func WithSink(ctx context.Context, sink cache.Sink) context.Context {
	return context.WithValue(ctx, contextSinkKey, sink)
}

// Sink extracts a cache.Sink from the context
func Sink(ctx context.Context) (cache.Sink, bool) {
	v, ok := ctx.Value(contextSinkKey).(cache.Sink)
	return v, ok
}

var (
	contextSinkKey = &contextKey{"Sink"}
)
