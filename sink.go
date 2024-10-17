package cache

import (
	"time"
)

// Sink receives data from a Get call
type Sink interface {
	// SetString sets the value to s.
	SetString(s string, e time.Time) error

	// SetBytes sets the value to the contents of v.
	// The caller retains ownership of v.
	SetBytes(v []byte, e time.Time) error

	// SetValue sets the value to the object v.
	// The caller retains ownership of v.
	SetValue(v any, e time.Time) error

	// Bytes returns the value encoded as a slice
	// of bytes
	Bytes() []byte

	// Len tells the length of the internally encoded
	// representation of the value
	Len() int

	// Expire returns the time whe this entry will
	// be evicted from the Cache
	Expire() time.Time

	// Reset empties the content of the Sink
	Reset()
}
