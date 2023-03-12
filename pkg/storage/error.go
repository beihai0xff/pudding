// Package storage provides the Storage interface and implementation.
// error.go defines the error of storage
package storage

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

var (
	// ErrBucketCreateFailed is the error when create bucket failed
	ErrBucketCreateFailed = errors.New("create bucket failed")
	// ErrBucketNotFound is the error when the bucket not found
	ErrBucketNotFound = bolt.ErrBucketNotFound
	// ErrKeyExist is the error when the key already exist
	ErrKeyExist = errors.New("the key already exist")
	// ErrSegmentExist is the error when the Bucket already exist
	ErrSegmentExist = errors.New("the segmentID already exist")
)
