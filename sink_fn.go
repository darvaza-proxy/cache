package cache

import (
	"time"

	"darvaza.org/core"
)

// SinkType describes a Sink factory.
type SinkType[T any] struct {
	Decode func([]byte, *T) error
	Encode func(*T) ([]byte, error)
	Clone  func(*T) (*T, bool)
}

// SetDefaults fills the gaps and identifies errors.
func (typ *SinkType[T]) SetDefaults() error {
	switch {
	case typ == nil:
		return core.ErrNilReceiver
	case typ.Decode == nil:
		return core.Wrap(core.ErrInvalid, "missing decoder")
	case typ.Encode == nil:
		return core.Wrap(core.ErrInvalid, "missing encoder")
	case typ.Clone == nil:
		typ.Clone = DefaultClone[T]
		return nil
	default:
		return nil
	}
}

// New creates a new Sink using the [SinkType] factory.
func (typ *SinkType[T]) New() (*SinkFn[T], error) {
	if err := typ.SetDefaults(); err != nil {
		return nil, err
	}

	return &SinkFn[T]{typ: typ}, nil
}

var (
	_ Sink       = (*SinkFn[any])(nil)
	_ TSink[any] = (*SinkFn[any])(nil)
)

// SinkFn is a [TSink] implementation made using a [SinkType] factory.
type SinkFn[T any] struct {
	ByteSink

	typ *SinkType[T]
	val *T
}

// IsZero tells if the [Sink] is valid but empty.
func (sink *SinkFn[_]) IsZero() bool {
	return sink.Valid() && sink.Len() == 0
}

// Valid tells if the [Sink] can be used.
func (sink *SinkFn[T]) Valid() bool {
	switch {
	case sink == nil,
		sink.typ == nil,
		sink.typ.Decode == nil,
		sink.typ.Encode == nil,
		sink.typ.Clone == nil:
		return false
	default:
		return true
	}
}

// SetBytes sets the object of the Sink and its expiration time
// from an encoded byte array.
func (sink *SinkFn[T]) SetBytes(b []byte, e time.Time) error {
	switch {
	case len(b) == 0:
		return ErrNoData
	case sink == nil:
		return core.ErrNilReceiver
	case !sink.Valid():
		return ErrInvalidSink
	default:
		v := new(T)
		if err := sink.typ.Decode(b, v); err != nil {
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

// SetValue sets the object of the Sink and its expiration time.
func (sink *SinkFn[T]) SetValue(v *T, e time.Time) error {
	switch {
	case v == nil:
		return ErrNoData
	case sink == nil:
		return core.ErrNilReceiver
	case !sink.Valid():
		return ErrInvalidSink
	default:
		b, err := sink.typ.Encode(v)
		if err != nil {
			// failed to encode
			return core.Wrap(err, "encode")
		}

		if err := sink.ByteSink.UnsafeSetBytes(b, e); err != nil {
			// failed to store bytes
			sink.Reset()
			return err
		}

		p, _ := sink.typ.Clone(v)
		sink.val = p

		return nil
	}
}

// Value returns a copy of the stored object
func (sink *SinkFn[T]) Value() (*T, bool) {
	if !sink.Valid() {
		return nil, false
	}

	// copy
	if sink.val != nil {
		if out, ok := sink.typ.Clone(sink.val); ok {
			return out, true
		}
	}

	// decode new
	if b := sink.Bytes(); len(b) > 0 {
		out := new(T)
		if err := sink.typ.Decode(b, out); err == nil {
			return out, true
		}
	}

	return nil, false
}
