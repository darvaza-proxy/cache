module darvaza.org/cache/x/memcache

go 1.19

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)

require (
	darvaza.org/cache v0.2.1
	darvaza.org/cache/x/simplelru v0.1.2
	darvaza.org/core v0.9.4
	darvaza.org/slog v0.5.2
)

require (
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/text v0.11.0 // indirect
)
