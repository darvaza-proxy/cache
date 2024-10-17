package cache

import (
	"bytes"
	"encoding/gob"
)

// NewGobSink creates a new Sink using Gob as encoding
// and stores the object in the provided pointer
func NewGobSink[T any]() TSink[T] {
	sink, err := NewGobSinkFactory[T]().New()
	if err != nil {
		panic(err)
	}
	return sink
}

// NewGobSinkFactory ...
func NewGobSinkFactory[T any]() *SinkType[T] {
	return &SinkType[T]{
		Decode: DefaultGobDecode[T],
		Encode: DefaultGobEncode[T],
		Clone:  DefaultClone[T],
	}
}

// DefaultGobDecode ...
func DefaultGobDecode[T any](b []byte, out *T) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	return dec.Decode(out)
}

// DefaultGobEncode ...
func DefaultGobEncode[T any](v *T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
