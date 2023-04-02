module github.com/darvaza-proxy/cache/x/memcache

go 1.19

replace (
	github.com/darvaza-proxy/cache => ../..
	github.com/darvaza-proxy/cache/x/simplelru => ../simplelru
)

require (
	darvaza.org/core v0.9.0
	darvaza.org/slog v0.5.0
	github.com/darvaza-proxy/cache v0.0.5
	github.com/darvaza-proxy/cache/x/simplelru v0.0.3
)

require (
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)
