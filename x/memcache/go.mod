module darvaza.org/cache/x/memcache

go 1.21

toolchain go1.22.6

require (
	darvaza.org/cache v0.3.1
	darvaza.org/cache/x/simplelru v0.1.10
	darvaza.org/core v0.15.3
	darvaza.org/slog v0.5.14
)

require (
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/text v0.20.0 // indirect
)

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)
