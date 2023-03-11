// Package storage provides the Storage interface and implementation.
// storage.go defines the Storage interface
package storage

import "github.com/beihai0xff/pudding/api/gen/pudding/types/v1"

const (
	// StartID is the start log id
	StartID uint64 = 0
)

// Storage defines the interface of storage
type Storage interface {
	// View to get message log
	// Read-Only transactions
	View(segmentID, sequence uint64) (*types.Message, error)
	// Insert Create or update a key
	// If the key already exist it will return
	Insert(msg *types.Message) (uint64, error)
	// Update will update a key
	// If the key not exist it will Create the key
	Update(bucket, key, value []byte) error
	// Delete a key from segmentID
	Delete(bucket, key []byte) error
	// Batch(func(tx *bolt.Tx) error)

	// CreateSegment will create a segment
	CreateSegment(segmentID uint64) error
	// DeleteSegment will delete a segment
	DeleteSegment(segmentID uint64) error
}
