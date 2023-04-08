module darvaza.org/cache/x/memcache

go 1.19

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)

require (
	darvaza.org/cache v0.1.0
	darvaza.org/cache/x/simplelru v0.1.1
	darvaza.org/core v0.9.2
	darvaza.org/slog v0.5.0
)

require (
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
