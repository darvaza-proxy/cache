module darvaza.org/cache/x/memcache

go 1.19

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)

require (
	darvaza.org/cache v0.2.2
	darvaza.org/cache/x/simplelru v0.1.4
	darvaza.org/core v0.11.2
	darvaza.org/slog v0.5.5
)

require (
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
