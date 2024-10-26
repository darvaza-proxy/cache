package cache

import (
	"darvaza.org/core"
)

var (
	// ErrInvalidSink tells the [Sink] isn't in a usable state.
	ErrInvalidSink = core.QuietWrap(core.ErrInvalid, "%s", "invalid sink")

	// ErrNoData indicates to data was provided.
	ErrNoData = core.Wrap(core.ErrInvalid, "no data")
)
