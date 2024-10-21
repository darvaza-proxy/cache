package cache

import "darvaza.org/core"

var (
	// ErrNoData indicates to data was provided.
	ErrNoData = core.Wrap(core.ErrInvalid, "no data")
)
