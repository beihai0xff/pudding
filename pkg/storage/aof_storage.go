// Package storage provides the Storage interface and implementation.
// aof_storage.go is the implementation of aofStorage.
// aofStorage is a append-only file storage.
package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/pkg/log"
)

const (
	// DefaultAOFDir is the default dir of aofStorage
	DefaultAOFDir = "./database/"

	// defaultInitialMmapSize is the initial size of the mmapped region. Setting this larger than
	// the potential max db size can prevent writer from blocking reader.
	// This only works for linux.
	defaultInitialMmapSize = 256 * 1024 * 1024 // 256 MB

	defaultBatchLimit      = 100
	defaultBatchInterval   = 500 * time.Millisecond
	defaultSegmentInterval = 1
)

var (
	defaultBucketName      = byte(0x30)
	defaultIndexBucketName = []byte("index")
)

// DefaultConfig will return a default config
var DefaultConfig = &Config{
	Dir:             DefaultAOFDir,
	SegmentInterval: defaultSegmentInterval,
	BatchInterval:   defaultBatchInterval,
	BatchLimit:      defaultBatchLimit,
	MMapSize:        defaultInitialMmapSize,
}

// Config is the config of aofStorage
type Config struct {
	// Dir is the file dir to the aofStorage file.
	Dir string
	// SegmentInterval is the interval of the bucket
	SegmentInterval uint64
	// BatchInterval is the maximum time before flushing the BatchTx.
	// default is 100ms
	BatchInterval time.Duration
	// BatchLimit is the maximum puts before flushing the BatchTx.
	// if puts >= BatchLimit, the BatchTx will be flushed.
	BatchLimit int
	// MMapSize is the initial size of the mmapped region. Setting this larger than
	// the potential max db size can prevent writer from blocking reader.
	MMapSize int
	// MustBeNewBucket if is true, will return error when create an exist segmentID
	MustBeNewBucket bool
}

// aofStorage is a append-only file storage.
// it will flush the BatchTx to disk every batchInterval or batchLimit.
// it use bolt as the backend storage.
type aofStorage struct {
	// dir is the file dir to the aofStorage file.
	dir string
	// mmapSize is the initial size of the mmapped region. Setting this larger than
	// the potential max db size can prevent writer from blocking reader.
	mmapSize        int
	db              map[uint64]*bolt.DB
	segmentInterval uint64
	// batchInterval is the maximum time before flushing the BatchTx.
	batchInterval time.Duration
	// batchLimit is the maximum puts before flushing the BatchTx.
	// if puts >= batchLimit, the BatchTx will be flushed.
	batchLimit int

	stopChan chan struct{}
	doneChan chan struct{}
}

// NewAOFStorage will create a new aofStorage
func NewAOFStorage(c *Config) (Storage, error) {
	return newStorage(c)
}

func newStorage(c *Config) (*aofStorage, error) {
	log.Infof("create storage dir: %s", filepath.Dir(c.Dir))
	err := os.MkdirAll(c.Dir, 0777)
	// Open the ./aofStorage.db data file in your current directory.
	// It will be created if it doesn't exist.
	if err != nil {
		return nil, fmt.Errorf("create dir error: %v", err)
	}

	s := &aofStorage{
		dir: c.Dir,
		db:  map[uint64]*bolt.DB{},

		segmentInterval: c.SegmentInterval,
		batchInterval:   c.BatchInterval,
		batchLimit:      c.BatchLimit,

		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}

	return s, err
}

// View a k/v pairs in Read-Only transactions.
func (s *aofStorage) View(segmentID, timestamp, sequence uint64) (*types.Message, error) {
	msg := types.Message{}

	return &msg, s.db[segmentID].View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.getBucketName(segmentID, timestamp))
		if b == nil {
			return ErrBucketNotFound
		}
		v := b.Get(uint64ToBytes(sequence))

		return proto.Unmarshal(v, &msg)
	})
}

// Insert will insert a key/value pair with the given segmentID.
func (s *aofStorage) Insert(msg *types.Message) (uint64, error) {
	var sequence uint64

	segmentID := getSegmentID(msg.DeliverAt, s.segmentInterval)

	db, ok := s.db[segmentID]
	if !ok {
		var err error
		if db, err = s.createSegmentDB(segmentID); err != nil {
			log.Errorf("create segment db [%d] error: %v", segmentID, err)
			return sequence, err
		}
	}

	bucketName := s.getBucketName(segmentID, msg.DeliverAt)
	if err := s.tryCreateDataBucket(db, bucketName); err != nil {
		log.Errorf("failed to create segment [%d] bucket [%b]", segmentID, defaultBucketName)
		return sequence, ErrBucketCreateFailed
	}

	return sequence, db.Update(func(tx *bolt.Tx) error {
		var err error

		b := tx.Bucket(bucketName)
		if b == nil {
			log.Errorf("failed to get segment [%d] bucket [%b]", segmentID, bucketName)
			return ErrBucketNotFound
		}

		// Generate sequence for the value.
		// This returns an error only if the Tx is closed or not writeable.
		if sequence, err = b.NextSequence(); err != nil {
			log.Errorf("failed to get segment [%d] bucket sequence: %v", segmentID, err)
			return err
		}
		sq := uint64ToBytes(sequence)
		if err = s.setIndex(tx, []byte(msg.Key), sq); err != nil {
			log.Errorf("failed to set index for segment [%d] bucket [%b] key [%s] value: %v", segmentID, bucketName, msg.Key, err)
			_ = tx.Rollback()
			return err
		}

		entry, err := proto.Marshal(msg)
		if err != nil {
			log.Errorf("failed to marshal message: %v", err)
			_ = tx.Rollback()
			return err
		}

		if err = b.Put(sq, entry); err != nil {
			log.Errorf("failed to put segment [%d] bucket [%b] key [%d] value: %v", segmentID, bucketName, msg.Key, err)
			_ = tx.Rollback()
			return err
		}

		// will auto commit in Update()
		return nil
	})
}

// tryCreateSegmentDB will create a db if it not exists
func (s *aofStorage) tryCreateDataBucket(db *bolt.DB, bucketName []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(bucketName) != nil {
			return nil
		}
		var err error
		var b *bolt.Bucket
		if b, err = tx.CreateBucket(bucketName); err != nil {
			log.Errorf("create [%s] data bucket [%s] error: %v", db.String(), bucketName, err)
			_ = tx.Rollback()
			return err
		}
		log.Infof("create [%s] data bucket [%s] success", db.String(), bucketName)

		if err = b.SetSequence(StartID); err != nil {
			log.Errorf("set [%s] data bucket sequence: %v", db.String(), err)
			_ = tx.Rollback()
			return err
		}
		log.Infof("set [%s] data bucket sequence from [%d] success", db.String(), StartID)

		return nil
	})
}

// tryCreateSegmentDB will create a db if it not exists
func (s *aofStorage) tryCreateIndexBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		var err error
		if _, err = tx.CreateBucket(defaultIndexBucketName); err != nil {
			log.Errorf("create [%s] index bucket [%s] error: %v", db.String(), defaultIndexBucketName, err)
			return err
		}
		log.Infof("create [%s] index bucket [%s] success", db.String(), defaultIndexBucketName)
		return nil
	})
}

// setIndex will Generate an index for the given value and store it.
func (s *aofStorage) setIndex(tx *bolt.Tx, key, index []byte) error {
	b := tx.Bucket(defaultIndexBucketName)
	if b == nil {
		return ErrBucketNotFound
	}

	// set index
	return b.Put(key, index)
}

// Update will update a key will given value
// If the key not exist, it will return an error.
func (s *aofStorage) Update(bucket, key, value []byte) error {
	return errors.New("not implemented")
}

// Delete a key from given segmentID.
func (s *aofStorage) Delete(bucket, key []byte) error {
	return errors.New("not implemented")
}

// CreateSegment create a segmentID
func (s *aofStorage) CreateSegment(segmentID uint64) error {
	_, err := s.createSegmentDB(segmentID)
	return err
}

// CreateSegment create a segmentID
func (s *aofStorage) createSegmentDB(segmentID uint64) (*bolt.DB, error) {
	db, err := s.tryCreateSegmentDB(segmentID)
	if err != nil {
		log.Errorf("create segment [%d] db error: %v", segmentID, err)
		return nil, err
	}

	if err := s.tryCreateIndexBucket(db); err != nil {
		log.Errorf("create segment [%d] index bucket error: %v", segmentID, err)
	}

	s.db[segmentID] = db

	return db, nil
}

// tryCreateSegmentDB will open an exist db file, or create a db if it not exists
func (s *aofStorage) tryCreateSegmentDB(segmentID uint64) (*bolt.DB, error) {
	path := getFilePath(segmentID, s.segmentInterval, s.dir)

	db, err := bolt.Open(path, 0666, &bolt.Options{Timeout: 3 * time.Second, InitialMmapSize: s.mmapSize})
	if err != nil {
		log.Errorf("create segment file: %s failed", path)
		return nil, err
	}

	log.Infof("create segment file: %s success", path)

	return db, nil
}

// DeleteSegment the given segmentID
func (s *aofStorage) DeleteSegment(segmentID uint64) error {
	path := getFilePath(segmentID, s.segmentInterval, s.dir)
	delete(s.db, segmentID)

	return os.Remove(path)
}

// getBucketName get bucket name by segmentID and deliverAt
func (s *aofStorage) getBucketName(segmentID, deliverAt uint64) []byte {
	return []byte{defaultBucketName + byte(deliverAt-segmentID)}
}
