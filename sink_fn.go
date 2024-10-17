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

// DefaultClone attempts to make a copy of a cached object
// using known interfaces.
func DefaultClone[T any](x *T) (*T, bool) {
	if x == nil {
		return nil, false
	}

	switch v := any(x).(type) {
	case interface{ Clone() *T }:
		return v.Clone(), true
	case interface{ Copy() *T }:
		return v.Copy(), true
	case interface{ Clone() T }:
		p := v.Clone()
		return &p, true
	case interface{ Copy() T }:
		p := v.Copy()
		return &p, true
	default:
		return nil, false
	}
}

var (
	_ Sink       = (*SinkFn[any])(nil)
	_ TSink[any] = (*SinkFn[any])(nil)
)

// SinkFn is a [TSink] implementation made using a [SinkType] factory.
type SinkFn[T any] struct {
	typ   *SinkType[T]
	val   *T
	exp   time.Time
	bytes []byte
}

// IsZero tells if the [Sink] is valid but empty.
func (p *SinkFn[_]) IsZero() bool {
	return p.Valid() && len(p.bytes) == 0
}

// Valid tells if the [Sink] can be used.
func (p *SinkFn[T]) Valid() bool {
	switch {
	case p == nil,
		p.typ == nil,
		p.typ.Decode == nil,
		p.typ.Encode == nil,
		p.typ.Clone == nil:
		return false
	default:
		return true
	}
}

// Bytes returns the encoded value stored in the [Sink].
func (p *SinkFn[_]) Bytes() []byte {
	if p != nil {
		return p.bytes
	}
	return nil
}

// Len returns the size in bytes of the encoded value stored in the [Sink].
func (p *SinkFn[_]) Len() int { return len(p.Bytes()) }

// Expire tells when the stored value expires.
func (p *SinkFn[_]) Expire() time.Time {
	if p == nil {
		return time.Time{}
	}

	return p.exp
}

// Reset discards any stored value.
func (p *SinkFn[T]) Reset() {
	if p != nil {
		if p.bytes != nil {
			p.bytes = nil
		}
		p.val = nil
		p.exp = time.Time{}
	}
}

// SetBytes ...
func (p *SinkFn[T]) SetBytes(b []byte, e time.Time) error {
	switch {
	case len(b) == 0:
		return core.Wrap(core.ErrInvalid, "empty value")
	case p == nil:
		return core.ErrNilReceiver
	case !p.Valid():
		return ErrInvalidSink
	default:
		v := new(T)
		if err := p.typ.Decode(b, v); err != nil {
			// failed to decode
			return core.Wrap(err, "decode")
		}

		p.exp = e
		p.val = v
		p.bytes = make([]byte, len(b))
		copy(p.bytes, b)

		return nil
	}
}

// SetValue ...
func (p *SinkFn[T]) SetValue(x *T, e time.Time) error {
	switch {
	case x == nil:
		return core.Wrap(core.ErrInvalid, "missing value")
	case p == nil:
		return core.ErrNilReceiver
	case !p.Valid():
		return ErrInvalidSink
	default:
		b, err := p.typ.Encode(x)
		if err != nil {
			// failed to encode
			return core.Wrap(err, "encode")
		}

		if v, ok := p.typ.Clone(x); ok {
			// store copy
			p.val = v
		} else {
			// must be decoded every time
			p.val = nil
		}

		p.exp = e
		p.bytes = b
		return nil
	}
}

// Value ...
func (p *SinkFn[T]) Value() (*T, bool) {
	switch {
	case p == nil:
		// invalid
		return nil, false
	case len(p.bytes) == 0:
		// empty
		return nil, false
	case p.val != nil:
		// clone
		return p.typ.Clone(p.val)
	default:
		// assemble new
		v := new(T)
		if err := p.typ.Decode(p.bytes, v); err != nil {
			// failed to decode
			return nil, false
		}
		return v, true
	}
}
