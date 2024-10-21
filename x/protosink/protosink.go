// Package protosink provides a cache.Sink using Protobuf
// encoding
package protosink

import (
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"darvaza.org/cache"
	"darvaza.org/core"
)

var (
	_ cache.Sink                = NewProtoSink[dummyMessage]()
	_ cache.TSink[dummyMessage] = NewProtoSink[dummyMessage]()
)

type dummyMessage struct{}

func (*dummyMessage) ProtoReflect() protoreflect.Message { return nil }

// ProtoMessage is a generic constraint for a type which pointer implements proto.Message
type ProtoMessage[T any] interface {
	proto.Message
	*T
}

// ProtoSink is a Sink using Proto for encoding
type ProtoSink[T any, M ProtoMessage[T]] struct {
	cache.ByteSink

	val *T
}

// Reset clears everything but the type pointer assigned during creation
func (sink *ProtoSink[T, M]) Reset() {
	if sink != nil {
		sink.ByteSink.Reset()
		sink.val = nil
	}
}

// SetBytes sets the object of the ProtoSink and its expiration time
// from a Protobuf encoded byte array
func (sink *ProtoSink[T, M]) SetBytes(b []byte, e time.Time) error {
	switch {
	case len(b) == 0:
		return cache.ErrNoData
	case sink == nil:
		return core.ErrNilReceiver
	default:
		v := new(T)
		if err := sink.doDecode(b, v); err != nil {
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

// SetValue sets the object of the ProtoSink and its expiration time
func (sink *ProtoSink[T, M]) SetValue(v *T, e time.Time) error {
	switch {
	case v == nil:
		return cache.ErrNoData
	case sink == nil:
		return core.ErrNilReceiver
	default:
		b, err := sink.doEncode(v)
		if err != nil {
			// failed to encode
			return core.Wrap(err, "encode")
		}

		if err := sink.ByteSink.UnsafeSetBytes(b, e); err != nil {
			// failed to store bytes
			sink.Reset()
			return err
		}

		p, _ := cache.DefaultClone(v)
		sink.val = p
		return nil
	}
}

// Value returns the stored object
func (sink *ProtoSink[T, M]) Value() (*T, bool) {
	if sink == nil {
		// invalid
		return nil, false
	}

	// copy
	if sink.val != nil {
		out, ok := cache.DefaultClone(sink.val)
		if ok {
			return out, true
		}
	}

	// decode new
	if b := sink.Bytes(); len(b) > 0 {
		out := new(T)
		if err := sink.doDecode(b, out); err == nil {
			return out, true
		}
	}

	return nil, false
}

func (*ProtoSink[T, M]) doEncode(v *T) ([]byte, error) {
	return proto.Marshal((M)(v))
}

func (*ProtoSink[T, M]) doDecode(b []byte, out *T) error {
	if out == nil {
		// errors only
		out = new(T)
	}

	return proto.Unmarshal(b, (M)(out))
}

// NewProtoSink creates a [ProtoSink] for a particular proto.Message type
func NewProtoSink[T any, M ProtoMessage[T]]() cache.TSink[T] {
	return new(ProtoSink[T, M])
}
