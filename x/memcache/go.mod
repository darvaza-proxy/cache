module darvaza.org/cache/x/memcache

go 1.21

require (
	darvaza.org/cache v0.3.2
	darvaza.org/cache/x/simplelru v0.1.11
	darvaza.org/core v0.15.6
	darvaza.org/slog v0.5.15
)

require (
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)
