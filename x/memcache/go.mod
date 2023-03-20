module github.com/darvaza-proxy/cache/x/memcache

go 1.19

replace (
	github.com/darvaza-proxy/cache => ../..
	github.com/darvaza-proxy/cache/x/simplelru => ../simplelru
)

require (
	github.com/darvaza-proxy/cache v0.0.5
	github.com/darvaza-proxy/cache/x/simplelru v0.0.3
	github.com/darvaza-proxy/core v0.6.5
	github.com/darvaza-proxy/slog v0.4.6
)

require (
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)
