module darvaza.org/cache/x/memcache

go 1.21

toolchain go1.22.6

require (
	darvaza.org/cache v0.2.8
	darvaza.org/cache/x/simplelru v0.1.10
	darvaza.org/core v0.15.1
	darvaza.org/slog v0.5.12
)

require (
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace (
	darvaza.org/cache => ../..
	darvaza.org/cache/x/simplelru => ../simplelru
)
