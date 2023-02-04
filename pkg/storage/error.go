// Package storage provides the Storage interface and implementation.
// error.go defines the error of storage
package storage

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

var (
	// ErrBucketNotFound is the error when the bucket not found
	ErrBucketNotFound = bolt.ErrBucketNotFound
	// ErrKeyExist is the error when the key already exist
	ErrKeyExist = errors.New("the key already exist")
	// ErrBucketExist is the error when the Bucket already exist
	ErrBucketExist = errors.New("the segmentID already exist")
)