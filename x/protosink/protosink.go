// Package protosink provides a cache.Sink using Protobuf
// encoding
package protosink

import (
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/darvaza-proxy/cache"
)

var (
	_ cache.Sink = (*ProtoSink)(nil)
)

// ProtoSink is a Sink using Proto for encoding
type ProtoSink struct {
	out proto.Message
	b   []byte
	e   time.Time
}

// Bytes returns the Protobuf encoded representation of the
// object in the Sink.
func (sink *ProtoSink) Bytes() []byte {
	return sink.b
}

// Len returns the length of the Protobuf encoded representation of the
// object in the Sink. 0 if empty
func (sink *ProtoSink) Len() int {
	return len(sink.b)
}

// Expire tells when this object will be evicted from the Cache
func (sink *ProtoSink) Expire() time.Time {
	return sink.e
}

// Reset clears everything but the type pointer assigned during creation
func (sink *ProtoSink) Reset() {
	sink.b = []byte{}
	sink.e = time.Time{}
	// TODO: how do we clear *sink.out ?
}

// SetBytes sets the object of the ProtoSink and its expiration time
// from a Protobuf encoded byte array
func (sink *ProtoSink) SetBytes(v []byte, e time.Time) error {
	// decode
	err := proto.Unmarshal(v, sink.out)
	if err != nil {
		sink.Reset()
		return cache.ErrInvalid
	}

	// store
	sink.b = make([]byte, len(v))
	sink.e = e
	copy(sink.b, v)
	return nil
}

// SetString sets the object of the ProtoSink and its expiration time
// from a Protobuf encoded string
func (sink *ProtoSink) SetString(s string, e time.Time) error {
	return sink.SetBytes([]byte(s), e)
}

// SetValue sets the object of the ProtoSink and its expiration time
func (sink *ProtoSink) SetValue(v any, e time.Time) error {
	p, ok := v.(proto.Message)
	if !ok {
		sink.Reset()
		return cache.ErrInvalid
	}

	out, err := proto.Marshal(p)
	if err != nil {
		// failed to encode
		sink.Reset()
		return cache.ErrInvalid
	}

	sink.b = out
	sink.e = e
	if sink.out != nil {
		sink.out = p
	}
	return nil
}

// Value returns the stored object
func (sink *ProtoSink) Value() proto.Message {
	return sink.out
}
