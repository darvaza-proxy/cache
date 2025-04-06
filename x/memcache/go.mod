module darvaza.org/cache/x/memcache

go 1.22

require (
	darvaza.org/cache v0.4.1
	darvaza.org/cache/x/simplelru v0.2.1
	darvaza.org/core v0.16.3
	darvaza.org/slog v0.6.2
)

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)
