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
	MarshalFunc    func(interface{}) ([]byte, error)
	UnmarshalFunc  func([]byte, interface{}) error
	CompressFunc   func(data []byte) []byte
	DecompressFunc func(b []byte, value interface{}) error

	OptionFunc func(*MsgPack)
)

type MsgPack struct {
	marshal    MarshalFunc
	unmarshal  UnmarshalFunc
	compress   CompressFunc
	decompress DecompressFunc
}

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

func WithMarshalFunc(fn MarshalFunc) OptionFunc {
	return func(p *MsgPack) {
		p.marshal = fn
	}
}

func WithUnmarshalFunc(fn UnmarshalFunc) OptionFunc {
	return func(p *MsgPack) {
		p.unmarshal = fn
	}
}

func WithCompressFunc(fn CompressFunc) OptionFunc {
	return func(p *MsgPack) {
		p.compress = fn
	}
}

func WithDecompressFunc(fn DecompressFunc) OptionFunc {
	return func(p *MsgPack) {
		p.decompress = fn
	}
}

/*
	Encode and decode like below:
*/

func Encode(item interface{}) ([]byte, error) {
	switch value := item.(type) {
	case nil:
		return nil, nil
	case []byte:
		return value, nil
	case string:
		return []byte(value), nil
	}

	b, err := defaultPack.marshal(item)
	if err != nil {
		return nil, err
	}

	return defaultPack.compress(b), nil
}

func Decode(b []byte, value interface{}) error {
	if len(b) == 0 {
		return nil
	}

	if err := defaultPack.decompress(b, value); err != nil {
		log.Errorf("Decompress failed: %v", err)
		return err
	}
	return defaultPack.unmarshal(b, value)
}

/*
	compress and decompress like below:
*/

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

func S2Decompress(b []byte, value interface{}) error {
	switch value := value.(type) {
	case nil:
		return nil
	case *[]byte:
		clone := make([]byte, len(b))
		copy(clone, b)
		*value = clone
		return nil
	case *string:
		*value = string(b)
		return nil
	}

	switch c := b[len(b)-1]; c {
	case noCompression:
		b = b[:len(b)-1]
	case s2Compression:
		b = b[:len(b)-1]

		var err error
		b, err = s2.Decode(nil, b)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown compression method: %x", c)
	}

	return nil
}
