// Package storage provides the Storage interface and implementation.
// storage.go defines the Storage interface
package storage

import (
	"time"
)

const (
	// defaultInitialMmapSize is the initial size of the mmapped region. Setting this larger than
	// the potential max db size can prevent writer from blocking reader.
	// This only works for linux.
	defaultInitialMmapSize = 256 * 1024 * 1024 // 256 MB

	defaultBatchLimit    = 100
	defaultBatchInterval = 500 * time.Millisecond

	// StartAt is the start log id
	StartAt = 0
)

// Storage defines the interface of storage
type Storage interface {
	// View to manage key/value pairs
	// Read-Only transactions
	View(bucket, key []byte) ([]byte, error)
	// Create or update a key
	// If the key already exist it will return
	Create(bucket, key, value []byte) error
	// Update will create or update a key
	// If the key not exist it will Create the key
	Update(bucket, key, value []byte) error
	// Delete a key from segmentID
	Delete(bucket, key []byte) error
	// Index Generate a index for the url and store it
	Index(value []byte) (uint64, error)
	//
	// Batch(func(tx *bolt.Tx) error)

	// CreateSegment will create a segment
	CreateSegment(segmentID uint64) error
	// DeleteSegment will delete a segment
	DeleteSegment(segmentID uint64) error
}
