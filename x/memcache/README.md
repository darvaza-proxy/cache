# In-memory cache with TTL and restricted size

[![Go Reference](https://pkg.go.dev/badge/github.com/darvaza-proxy/cache/x/memcache.svg)](https://pkg.go.dev/github.com/darvaza-proxy/cache/x/memcache)

## LRU

This package provides a thread-safe []byte TTL with maximum size control based on our
[simplelru](https://pkg.go.dev/github.com/darvaza-proxy/cache/x/simplelru).

This LRU does **NOT** implement the Cache interface.

# See also

* [Cache](https://pkg.go.dev/github.com/darvaza-proxy/cache)
* [Groupcache](https://pkg.go.dev/github.com/darvaza-proxy/cache/x/groupcache)
* [simplelru](https://pkg.go.dev/github.com/darvaza-proxy/cache/x/simplelru)
