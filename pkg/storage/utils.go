package storage

import "encoding/binary"

// convert uint64 to 8-byte big endian representation
func uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
