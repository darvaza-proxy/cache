package cache

import (
	"errors"

	"darvaza.org/core"
)

var (
	// ErrInvalidSink tells the [Sink] isn't in a usable state.
	ErrInvalidSink = errors.New("invalid sink")

	// ErrNoData indicates to data was provided.
	ErrNoData = core.Wrap(core.ErrInvalid, "no data")
)
