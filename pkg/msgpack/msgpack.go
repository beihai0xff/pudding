// Package msgpack provides a msgpack codec for encoding and decoding
package msgpack

import (
	"fmt"

	"github.com/klauspost/compress/s2"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/beihai0xff/pudding/pkg/log"
)

const (
	compressionThreshold = 64

	noCompression = 0x0
	s2Compression = 0x1
)

var defaultPack = &MsgPack{
	marshal:    msgpack.Marshal,
	unmarshal:  msgpack.Unmarshal,
	compress:   S2Compress,
	decompress: S2Decompress,
}

type (
	// MarshalFunc is a function that marshals a value into a byte array.
	MarshalFunc func(interface{}) ([]byte, error)
	// UnmarshalFunc is a function that unmarshals a byte array into a value.
	UnmarshalFunc func([]byte, interface{}) error
	// CompressFunc is a function that compresses a byte array.
	CompressFunc func(data []byte) []byte
	// DecompressFunc is a function that decompresses a byte array.
	DecompressFunc func(b []byte) ([]byte, error)

	// OptionFunc is a function that configures a MsgPack.
	OptionFunc func(*MsgPack)
)

// MsgPack is a msgpack codec.
type MsgPack struct {
	marshal    MarshalFunc
	unmarshal  UnmarshalFunc
	compress   CompressFunc
	decompress DecompressFunc
}

// New creates a new MsgPack.
func New(opts ...OptionFunc) *MsgPack {
	pack := *defaultPack
	for _, opt := range opts {
		opt(defaultPack)
	}
	return &pack
}

/*
	Functional Options Pattern
*/

// WithMarshalFunc sets the marshal function.
func WithMarshalFunc(fnm MarshalFunc, fnu UnmarshalFunc) OptionFunc {
	return func(p *MsgPack) {
		p.marshal = fnm
		p.unmarshal = fnu
	}
}

// WithCompressFunc sets the compress function.
func WithCompressFunc(fnc CompressFunc, fnd DecompressFunc) OptionFunc {
	return func(p *MsgPack) {
		p.compress = fnc
		p.decompress = fnd
	}
}

/*
	Encode and decode like below:
*/

// Encode wrap for msgpack.Encode
func Encode(item interface{}) ([]byte, error) {
	return defaultPack.Encode(item)
}

// Decode wrap for msgpack.Decode
func Decode(b []byte, value interface{}) error {
	return defaultPack.Decode(b, value)
}

// Encode wraps msgpack.Marshal and compresses the result.
func (p *MsgPack) Encode(item interface{}) ([]byte, error) {
	switch value := item.(type) {
	case nil:
		return nil, nil
	case []byte:
		return value, nil
	case string:
		return []byte(value), nil
	}

	b, err := p.marshal(item)
	if err != nil {
		return nil, err
	}

	return p.compress(b), nil
}

// Decode a msgpack encoded byte array
func (p *MsgPack) Decode(b []byte, value interface{}) error {
	if len(b) == 0 {
		return nil
	}

	var err error

	if b, err = p.decompress(b); err != nil {
		log.Errorf("Decompress failed: %v", err)
		return err
	}
	return p.unmarshal(b, value)
}

/*
	compress and decompress like below:
*/

// S2Compress compresses a byte array using s2.
func S2Compress(data []byte) []byte {
	// if data length is less than compressionThreshold, skip compress.
	if len(data) < compressionThreshold {
		n := len(data) + 1
		b := make([]byte, n)
		copy(b, data)
		b[len(b)-1] = noCompression
		return b
	}

	n := s2.MaxEncodedLen(len(data)) + 1
	b := make([]byte, n)
	b = s2.Encode(b, data)
	// use the last byte to store positive compression method
	b = append(b, s2Compression)
	return b
}

// S2Decompress decompresses a byte array using s2.
func S2Decompress(b []byte) ([]byte, error) {
	switch c := b[len(b)-1]; c {
	case noCompression:
		b = b[:len(b)-1]
	case s2Compression:
		b = b[:len(b)-1]

		var err error
		b, err = s2.Decode(nil, b)
		if err != nil {
			return b, err
		}
	default:
		return nil, fmt.Errorf("unknown compression method: %x", c)
	}

	return b, nil
}
