# In-memory cache with TTL and restricted size

[![Go Reference][godoc-badge]][godoc-link]

## LRU

This package provides a thread-safe []byte TTL with maximum size control
based on our [simplelru][simplelru-link].

This LRU does **NOT** implement the Cache interface.

## See also

* [Cache][cache-link]
* [Groupcache][groupcache-link]
* [simplelru][simplelru-link]

[godoc-link]: https://pkg.go.dev/darvaza.org/cache/x/memcache
[godoc-badge]: https://pkg.go.dev/badge/darvaza.org/cache/x/memcache.svg
[cache-link]: https://pkg.go.dev/darvaza.org/cache
[groupcache-link]: https://pkg.go.dev/darvaza.org/cache/x/groupcache
[simplelru-link]: https://pkg.go.dev/darvaza.org/cache/x/simplelru
