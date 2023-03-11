package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/pkg/log"
)

var testAof *aofStorage

func TestMain(m *testing.M) {
	var err error
	DefaultConfig.Dir = "./test/output/database/aof.log"
	testAof, err = newStorage(DefaultConfig)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	if err = os.RemoveAll("./test"); err != nil {
		log.Errorf("remove test database dir failed: %v", err)
	}
	os.Exit(exitCode)
}

func Test_aofStorage_View(t *testing.T) {
	testMsg := &types.Message{Key: "test_key", Payload: []byte("test_payload"), DeliverAt: 120}
	_ = testAof.CreateSegment(120)
	sequence, _ := testAof.Insert(testMsg)
	assert.Equal(t, uint64(1), sequence)
	tests := []struct {
		name                string
		segmentID, sequence uint64
		want                *types.Message
		wantErr             assert.ErrorAssertionFunc
	}{
		{"test_short_url", 120, sequence, testMsg, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testAof.View(tt.segmentID, tt.sequence)
			tt.wantErr(t, err)
			if !proto.Equal(got, tt.want) {
				t.Errorf("View() got = %v, want %v", got, tt.want)
			}
		})
	}

	testAof.DeleteSegment(120)
}

func Test_aofStorage_Insert(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.Message
		want    uint64
		wantErr assert.ErrorAssertionFunc
	}{
		{"test_Insert", &types.Message{Key: "test_key", Payload: []byte("test_payload"), DeliverAt: 60}, 1, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testAof.Insert(tt.msg)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equalf(t, tt.want, got, "Insert() got = %v, want %v", got, tt.want)
		})
	}
	_ = testAof.DeleteSegment(60)
}

func Test_aofStorage_CreateSegment(t *testing.T) {
	tests := []struct {
		name      string
		segmentID uint64
		wantErr   assert.ErrorAssertionFunc
	}{
		{"test_Create_non-exist_Segment", 123456, assert.NoError},
		{"test_Create_exist_Segment", 123456, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr(t, testAof.CreateSegment(tt.segmentID)) {
				t.Errorf("create segment error, segmentID = %v", tt.segmentID)
			}
			msg := &types.Message{
				Topic:        "test_topic",
				Key:          "test_key",
				Payload:      []byte("test_value"),
				DeliverAfter: 0,
				DeliverAt:    123456,
			}
			sequence, _ := testAof.Insert(msg)
			value, _ := testAof.View(tt.segmentID, sequence)
			if !proto.Equal(value, msg) {
				t.Errorf("CreateSegment get value = %v, want %v", value, msg)
			}
		})
	}
}

func Test_aofStorage_tryCreateBucket(t *testing.T) {
	db, _ := testAof.tryCreateSegmentDB(123)
	tests := []struct {
		name           string
		db             *bolt.DB
		errorAssertion assert.ErrorAssertionFunc
	}{
		{"test_CreateBucket_no_exist", db, assert.NoError},
		{"test_CreateBucket_exist", db, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.errorAssertion(t, testAof.tryCreateBucket(tt.db))
		})
	}
	_ = testAof.DeleteSegment(123)
}

func Test_aofStorage_DeleteSegment(t *testing.T) {
	_ = testAof.CreateSegment(123456)
	tests := []struct {
		name           string
		segmentID      uint64
		errorAssertion assert.ErrorAssertionFunc
	}{
		{"test_DeleteExistSegment", 123456, assert.NoError},
		{"test_DeleteNotExistSegment", 654321, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.errorAssertion(t, testAof.DeleteSegment(tt.segmentID)) {
				if _, ok := testAof.db[tt.segmentID]; ok {
					t.Errorf("DeleteSegment [%d] failed", tt.segmentID)
				}
			}
		})
	}
}
