package storage

import (
	"os"
	"reflect"
	"testing"
)

var testAof *aofStorage

func TestMain(m *testing.M) {
	var err error
	DefaultConfig.Path = "/tmp/pudding/database/aof.log"
	testAof, err = newStorage(DefaultConfig)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	_ = os.RemoveAll("/tmp/pudding/")
	os.Exit(exitCode)
}

func Test_aofStorage_Index(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
		want  uint64
	}{
		{"test_short_url", []byte("https://www.google.com/"), 1},
		{"test_long_url", []byte("https://www.google.com/search?sxsrf=ALeKk00rEgE8Gd7-KSZTZUxVkWSzq6exKw%3A1592901527548&ei=l7_xXrmMIYfd9QOOv6CACg&q=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&oq=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&gs_lcp=CgZwc3ktYWIQAzIECCMQJzIECAAQHlDIKVj2K2CdLmgAcAB4AIABiQOIAcMHkgEFMi0yLjGYAQCgAQGqAQdnd3Mtd2l6&sclient=psy-ab&ved=0ahUKEwj5s9jNxJfqAhWHbn0KHY4fCKAQ4dUDCAw&uact=5"),
			2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testAof.Index(tt.value)
			if err != nil {
				t.Errorf("Index() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Index() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aofStorage_Update(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		key    []byte
		value  []byte
	}{
		{"test_short_url", []byte("index"), uint64ToBytes(uint64(123457)), []byte("https://cn.bing.com/")},
		{"test_long_url", []byte("index"), uint64ToBytes(uint64(123458)), []byte("https://cn.bing.com/search?q=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&qs=n&form=QBLHCN&sp=-1&pq=zhao%27hua%27xi%27shi&sc=3-15&sk=&cvid=4FA5BBC53EA84E6B93A6DEC3F006AA4D")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testAof.Update(tt.bucket, tt.key, tt.value); err != nil {
				t.Errorf("Update() error = %v", err)
			}
		})
	}
}

func Test_aofStorage_View(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		key    []byte
		want   []byte
	}{
		{"test_short_url", []byte("index"), uint64ToBytes(uint64(123457)), []byte("https://cn.bing.com/")},
		{"test_long_url", []byte("index"), uint64ToBytes(uint64(123458)), []byte("https://cn.bing.com/search?q=%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE&qs=n&form=QBLHCN&sp=-1&pq=zhao%27hua%27xi%27shi&sc=3-15&sk=&cvid=4FA5BBC53EA84E6B93A6DEC3F006AA4D")},
	}
	Test_aofStorage_Update(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testAof.View(tt.bucket, tt.key)
			if err != nil {
				t.Errorf("View() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("View() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aofStorage_Delete(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		key    []byte
	}{
		{"test_short_url", []byte("index"), uint64ToBytes(uint64(123457))},
		{"test_long_url", []byte("index"), uint64ToBytes(uint64(123458))},
	}
	Test_aofStorage_Update(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testAof.Delete(tt.bucket, tt.key); err != nil {
				t.Errorf("Delete() error = %v", err)
			}
		})
	}
}

func Test_aofStorage_tryCreateBucket(t *testing.T) {
	tests := []struct {
		name   string
		bucket []byte
		exits  bool
	}{
		{"test_CreateBucket", []byte("testCreateBucket"), false},
		{"test_CreateBucket", []byte("testCreateBucket"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testAof.tryCreateBucket(tt.bucket, false)
			if err != nil {
				t.Errorf("tryCreateBucket() error = %v", err)
				return
			}
			if got != tt.exits {
				t.Errorf("tryCreateBucket() got = %v, exits %v", got, tt.exits)
			}
		})
	}
}

func Test_aofStorage_CreateSegment(t *testing.T) {
	tests := []struct {
		name      string
		segmentID uint64
	}{
		{"test_CreateBucket", 123456},
		{"test_CreateBucket", 654321},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testAof.CreateSegment(tt.segmentID)
			if err != nil {
				t.Errorf("CreateBucket() error = %v, ", err)
				return
			}
		})
	}
}

func Test_aofStorage_DeleteBucket(t *testing.T) {
	tests := []struct {
		name      string
		segmentID uint64
		wantErr   error
	}{
		{"test_DeleteExistBucket", 123456, nil},
		{"test_DeleteNotExistBucket", 654321, ErrBucketNotFound},
	}
	err := testAof.CreateSegment(tests[0].segmentID)
	if err != nil {
		t.Errorf("CreateBucket() error = %v, ", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testAof.DeleteSegment(tt.segmentID); err != nil {
				if err == tt.wantErr {
					return
				} else {
					t.Errorf("DeleteBucket() error = %v", err)
				}
			}
		})
	}
}
