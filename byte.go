package cache

import (
	"time"

	"darvaza.org/core"
)

var _ Sink = (*ByteSink)(nil)

// ByteSink implements the simplest form of Sink
type ByteSink struct {
	bytes []byte
	exp   time.Time
}

// Bytes returns the stored bytes
func (s *ByteSink) Bytes() []byte {
	if s != nil {
		return s.bytes
	}
	return nil
}

// Len returns the length of the store bytes
func (s *ByteSink) Len() int {
	return len(s.Bytes())
}

// Expire indicates when the stored bytes are due to expire
func (s *ByteSink) Expire() time.Time {
	if s != nil {
		return s.exp
	}
	return time.Time{}
}

// Reset blanks the stored data
func (s *ByteSink) Reset() {
	if s != nil {
		// truncate
		s.bytes = s.bytes[:0]
		s.exp = time.Time{}
	}
}

// SetBytes stores a copy of the given bytes and expiration date.
func (s *ByteSink) SetBytes(b []byte, e time.Time) error {
	switch {
	case s == nil:
		return core.ErrNilReceiver
	case len(b) == 0:
		return ErrNoData
	default:
		// store copy
		s.copyBytes(b)
		s.exp = e
		return nil
	}
}

func (s *ByteSink) copyBytes(b []byte) {
	n := len(b)
	if n <= cap(s.bytes) {
		s.bytes = s.bytes[:n]
	} else {
		s.bytes = make([]byte, n)
	}

	copy(s.bytes, b)
}

// UnsafeSetBytes is like SetBytes but it uses the given byte slice
// directly instead of making a copy.
func (s *ByteSink) UnsafeSetBytes(b []byte, e time.Time) error {
	switch {
	case s == nil:
		return core.ErrNilReceiver
	case len(b) == 0:
		return ErrNoData
	default:
		// store as-is
		s.bytes = b
		s.exp = e
		return nil
	}
}
