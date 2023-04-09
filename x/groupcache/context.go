package groupcache

import (
	"darvaza.org/cache"
	"darvaza.org/core"
)

var (
	ctxSink = core.NewContextKey[cache.Sink]("groupcache.Sink")
)
