package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"darvaza.org/core"
)

var (
	_ Sink       = (*GobSink[any])(nil)
	_ TSink[any] = (*GobSink[any])(nil)
)

// GobSink is a Sink using generics for type safety and Gob
// for encoding
type GobSink[T any] struct {
	ByteSink

	val *T
}

// SetBytes sets the object of the GobSink and its expiration time
// from a Gob encoded byte array
func (sink *GobSink[T]) SetBytes(b []byte, e time.Time) error {
	switch {
	case len(b) == 0:
		return ErrNoData
	case sink == nil:
		return core.ErrNilReceiver
	default:
		v := new(T)
		if err := DecodeGob(b, v); err != nil {
			// failed to decode
			return core.Wrap(err, "decode")
		}

		if err := sink.ByteSink.SetBytes(b, e); err != nil {
			// failed to store bytes
			sink.Reset()
			return err
		}

		sink.val = v
		return nil
	}
}

// SetValue sets the object of the GobSink and its expiration time
func (sink *GobSink[T]) SetValue(v *T, e time.Time) error {
	switch {
	case v == nil:
		return ErrNoData
	case sink == nil:
		return core.ErrNilReceiver
	default:
		b, err := EncodeGob(v)
		if err != nil {
			// failed to encode
			return core.Wrap(err, "encode")
		}

		if err := sink.ByteSink.UnsafeSetBytes(b, e); err != nil {
			// failed to store bytes
			sink.Reset()
			return err
		}

		p, _ := DefaultClone(v)
		sink.val = p

		return nil
	}
}

// Value gives a copy of the stored object
func (sink *GobSink[T]) Value() (*T, bool) {
	if sink == nil {
		// invalid
		return nil, false
	}

	// copy
	if sink.val != nil {
		out, ok := DefaultClone(sink.val)
		if ok {
			return out, true
		}
	}
	// decode new
	if b := sink.Bytes(); len(b) > 0 {
		out := new(T)
		if err := DecodeGob(b, out); err == nil {
			return out, true
		}
	}

	return nil, false
}

// Reset clears everything but the type pointer
// assigned during creation
func (sink *GobSink[T]) Reset() {
	if sink != nil {
		sink.ByteSink.Reset()
		sink.val = nil
	}
}

// DecodeGob attempts to transform a byte slice into a Go
// object using Gob encoding.
func DecodeGob[T any](b []byte, out *T) error {
	if len(b) == 0 {
		return ErrNoData
	}

	if out == nil {
		// error-only call
		out = new(T)
	}

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	return dec.Decode(out)
}

// EncodeGob attempts to transform a Go object to a
// bytes slice using Gob encoding.
func EncodeGob[T any](p *T) ([]byte, error) {
	var buf bytes.Buffer

	if p == nil {
		return nil, ErrNoData
	}

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
