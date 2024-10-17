package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"darvaza.org/core"
)

var (
	_ Sink = (*GobSink[any])(nil)
)

// GobSink is a Sink using generics for type safety and Gob
// for encoding
type GobSink[T any] struct {
	out   *T
	e     time.Time
	bytes []byte
}

// Bytes returns the Gob encoded representation of the object
// in the Sink
func (sink *GobSink[T]) Bytes() []byte {
	return sink.bytes
}

// Len returns the length of the Gob encoded representation of the
// object in the Sink. 0 if empty.
func (sink *GobSink[T]) Len() int {
	return len(sink.bytes)
}

// Expire tells when this object will be evicted from the Cache
func (sink *GobSink[T]) Expire() time.Time {
	return sink.e
}

// SetBytes sets the object of the GobSink and its expiration time
// from a Gob encoded byte array
func (sink *GobSink[T]) SetBytes(v []byte, e time.Time) error {
	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(sink.out); err != nil {
		// failed to decode
		sink.Reset()
		return err
	}

	// store
	sink.bytes = make([]byte, len(v))
	sink.e = e
	copy(sink.bytes, v)
	return nil
}

// SetString isn't supported by GobSink, but its needed to satisfy
// the Sink interface
func (sink *GobSink[T]) SetString(string, time.Time) error {
	sink.Reset()
	return core.ErrInvalid
}

// SetValue sets the object of the GobSink and its expiration time
func (sink *GobSink[T]) SetValue(v any, e time.Time) error {
	var buf bytes.Buffer

	p, ok := v.(*T)
	if !ok {
		return core.ErrInvalid
	}

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		sink.Reset()
		return err
	}

	sink.bytes = buf.Bytes()
	sink.e = e
	*sink.out = *p
	return nil
}

// Value gives a copy of the stored object
func (sink *GobSink[T]) Value() T {
	if sink.out == nil {
		var zero T
		return zero
	}
	return *sink.out
}

// Reset clears everything but the type pointer
// assigned during creation
func (sink *GobSink[T]) Reset() {
	sink.bytes = []byte{}
	sink.e = time.Time{}

	if sink.out != nil {
		var zero T
		*sink.out = zero
	}
}

// NewGobSink creates a new Sink using Gob as encoding
// and stores the object in the provided pointer
func NewGobSink[T any](out *T) *GobSink[T] {
	if out == nil {
		var zero T
		core.Panicf("NewGobSink[%T]: output can not be nil", zero)
	}
	return &GobSink[T]{out: out}
}
