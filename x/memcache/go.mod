module darvaza.org/cache/x/memcache

go 1.20

require (
	darvaza.org/cache v0.2.8
	darvaza.org/cache/x/simplelru v0.1.10
	darvaza.org/core v0.14.10
	darvaza.org/slog v0.5.11
)

require (
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)
