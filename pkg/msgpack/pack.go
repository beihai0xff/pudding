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

var p *MsgPack

type (
	MarshalFunc   func(interface{}) ([]byte, error)
	UnmarshalFunc func([]byte, interface{}) error
)

type MsgPack struct {
	marshal   MarshalFunc
	unmarshal UnmarshalFunc
}

func init() {
	p = &MsgPack{
		marshal:   p._marshal,
		unmarshal: p._unmarshal,
	}
}

func Encode(item interface{}) ([]byte, error) {
	return p.marshal(item)
}

func Decode(b []byte, value interface{}) error {
	return p.unmarshal(b, value)
}

func (p *MsgPack) _marshal(item interface{}) ([]byte, error) {
	switch value := item.(type) {
	case nil:
		return nil, nil
	case []byte:
		return value, nil
	case string:
		return []byte(value), nil
	}

	b, err := msgpack.Marshal(item)
	if err != nil {
		return nil, err
	}

	return compress(b), nil
}

func (p *MsgPack) _unmarshal(b []byte, value interface{}) error {
	if len(b) == 0 {
		return nil
	}

	if err := p.decompress(b, value); err != nil {
		log.Errorf("decompress failed: %v", err)
		return err
	}

	return msgpack.Unmarshal(b, value)
}

func compress(data []byte) []byte {
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
	b = append(b, s2Compression)
	return b
}

func (p *MsgPack) decompress(b []byte, value interface{}) error {
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
