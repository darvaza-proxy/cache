module darvaza.org/cache/x/memcache

go 1.19

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)

require (
	darvaza.org/cache v0.2.6
	darvaza.org/cache/x/simplelru v0.1.8
	darvaza.org/core v0.14.6
	darvaza.org/slog v0.5.8
)

require (
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
