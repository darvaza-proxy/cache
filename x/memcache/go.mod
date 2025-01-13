module darvaza.org/cache/x/memcache

go 1.22

require (
	darvaza.org/cache v0.3.3
	darvaza.org/cache/x/simplelru v0.1.11
	darvaza.org/core v0.16.0
	darvaza.org/slog v0.6.0
)

require (
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)
