package cache

import (
	"time"
)

// Sink receives data from a Get call
type Sink interface {
	// SetBytes sets the value to the contents of v.
	// The caller retains ownership of v.
	SetBytes(v []byte, e time.Time) error

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

// TSink extends [Sink] with type-specific SetValue/Value
// methods.
type TSink[T any] interface {
	Sink

	// SetValue sets the value to the object v.
	// The caller retains ownership of v.
	SetValue(v *T, e time.Time) error

	// Value returns a copy of the object previously
	// stored in the Sink.
	Value() (*T, bool)
}
