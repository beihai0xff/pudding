package storage

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
)

// convert uint64 to 8-byte big endian representation
func uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)

	return b
}

func getFilePath(segmentID, interval uint64, dir string) string {
	startAt := segmentID
	endAt := startAt + interval
	fileName := fmt.Sprintf("segment_%d-%d.log", startAt, endAt)

	return filepath.Join(dir, fileName)
}

func getSegmentID(deliverAt, interval uint64) uint64 {
	return (deliverAt / interval) * interval
}
