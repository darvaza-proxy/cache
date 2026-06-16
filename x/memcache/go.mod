module darvaza.org/cache/x/memcache

go 1.24.0

require (
	darvaza.org/cache v0.5.0
	darvaza.org/cache/x/simplelru v0.3.0
	darvaza.org/core v0.19.1
	darvaza.org/slog v0.9.1
)

require (
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/text v0.34.0 // indirect
)

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)
