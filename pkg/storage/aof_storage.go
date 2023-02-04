// Package storage provides the Storage interface and implementation.
// aof_storage.go is the implementation of aofStorage.
// aofStorage is a append-only file storage.
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/beihai0xff/pudding/pkg/log"
)

// DefaultConfig will return a default config
var DefaultConfig = &Config{
	Path:          "./database/aof.log",
	BatchInterval: defaultBatchInterval,
	BatchLimit:    defaultBatchLimit,
	MmapSize:      defaultInitialMmapSize,
}

// Config is the config of aofStorage
type Config struct {
	// Path is the file path to the aofStorage file.
	Path string
	// BatchInterval is the maximum time before flushing the BatchTx.
	// default is 100ms
	BatchInterval time.Duration
	// BatchLimit is the maximum puts before flushing the BatchTx.
	// if puts >= BatchLimit, the BatchTx will be flushed.
	BatchLimit int
	// MmapSize is the initial size of the mmapped region. Setting this larger than
	// the potential max db size can prevent writer from blocking reader.
	MmapSize int
	// MustBeNewBucket if is true, will return error when create an exist segmentID
	MustBeNewBucket bool
}

// aofStorage is a append-only file storage.
// it will flush the BatchTx to disk every batchInterval or batchLimit.
// it use bolt as the backend storage.
type aofStorage struct {
	db *bolt.DB

	// batchInterval is the maximum time before flushing the BatchTx.
	batchInterval time.Duration
	// batchLimit is the maximum puts before flushing the BatchTx.
	// if puts >= batchLimit, the BatchTx will be flushed.
	batchLimit int

	buckets map[*bolt.Bucket]time.Duration

	stopChan chan struct{}
	doneChan chan struct{}
}

// NewAOFStorage will create a new aofStorage
func NewAOFStorage(c *Config) (Storage, error) {
	return newStorage(c)
}

func newStorage(c *Config) (*aofStorage, error) {
	log.Infof("create storage dir: %s", filepath.Dir(c.Path))
	err := os.MkdirAll(filepath.Dir(c.Path), 0777)
	// Open the ./aofStorage.db data file in your current directory.
	// It will be created if it doesn't exist.
	if err != nil {
		return nil, fmt.Errorf("create dir error: %v", err)
	}
	db, err := bolt.Open(c.Path, 0600, &bolt.Options{Timeout: 3 * time.Second, InitialMmapSize: c.MmapSize})
	if err != nil {
		return nil, err
	}
	s := &aofStorage{
		db: db,

		batchInterval: c.BatchInterval,
		batchLimit:    c.BatchLimit,

		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
	exist, err := s.tryCreateBucket([]byte("index"), true)
	if exist && c.MustBeNewBucket {
		return nil, ErrBucketExist
	}
	return s, err
}

// View a k/v pairs in Read-Only transactions.
func (s *aofStorage) View(bucket, key []byte) ([]byte, error) {
	var v []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		v = b.Get(key)
		return nil
	})
	// if the key not exist will return nil
	return v, err
}

// Create will insert a key/value pair with the given segmentID.
func (s *aofStorage) Create(bucket, key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		if v := b.Get(key); v != nil {
			return ErrKeyExist
		}
		return b.Put(key, value)
	})
}

// Update will create or update a key
// If the key not exist it will Create the key
func (s *aofStorage) Update(bucket, key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

// Delete a key from given segmentID.
func (s *aofStorage) Delete(bucket, key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}
		return b.Delete(key)
	})
}

// Index will Generate a index for the given value and store it.
func (s *aofStorage) Index(value []byte) (uint64, error) {
	var index uint64
	return index, s.db.Update(func(tx *bolt.Tx) error {
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("index"))
		if b == nil {
			return ErrBucketNotFound
		}

		// Generate index for the url.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		index, _ = b.NextSequence()

		// Persist bytes to url segmentID.
		return b.Put(uint64ToBytes(index), value)
	})
}

// CreateSegment create a segmentID
func (s *aofStorage) CreateSegment(segmentID uint64) error {
	bucketName := []byte(fmt.Sprintf("segment_%d.log", segmentID))
	_, err := s.tryCreateBucket(bucketName, false)
	return err
}

// tryCreateBucket will create a Bucket if it not exists
// the field exist tells the caller whether the Bucket already exists.
func (s *aofStorage) tryCreateBucket(bucketName []byte, start bool) (bool, error) {
	var exist bool
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			// Bucket not exist, create a new Bucket
			exist = false
			b, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return fmt.Errorf("create bucketName %b test_aof failed: %w", bucketName, err)
			}
			if start {
				err = b.SetSequence(StartAt)
				if err != nil {
					log.Errorf("failed to set sequence: %v", err)
				}
			}
		} else {
			exist = true // Bucket exists
		}
		return nil
	})
	return exist, err
}

// DeleteSegment the given segmentID
func (s *aofStorage) DeleteSegment(segmentID uint64) error {
	bucketName := []byte(fmt.Sprintf("segment_%d.log", segmentID))
	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucketName)
	})
	if err != nil {
		if err == ErrBucketNotFound {
			err = nil
		}
	}

	return err
}
