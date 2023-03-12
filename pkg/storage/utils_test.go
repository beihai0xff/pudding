package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_uint64ToBytes(t *testing.T) {
	tests := []struct {
		name string
		args uint64
		want []byte
	}{
		// don't need test
		{"test 0", uint64(0), []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{"test 123456", uint64(123456), []byte{0, 0, 0, 0, 0, 0x1, 0xe2, 0x40}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uint64ToBytes(tt.args)
			assert.Equal(t, tt.want, got)
			fmt.Println(got)
		})
	}
}

func Test_getFilePath(t *testing.T) {
	type args struct {
		segmentID uint64
		interval  uint64
		dir       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test 0~60", args{0, 60, "/tmp"}, "/tmp/segment_0-60.log"},
		{"test 60~120", args{60, 60, "/tmp"}, "/tmp/segment_60-120.log"},
		{"test 0~1", args{0, 1, "/tmp"}, "/tmp/segment_0-1.log"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getFilePath(tt.args.segmentID, tt.args.interval, tt.args.dir), "getFilePath(%v, %v, %v)", tt.args.segmentID, tt.args.interval, tt.args.dir)
		})
	}
}

func Test_getSegmentID(t *testing.T) {
	type args struct {
		deliverAt uint64
		interval  uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"test 0~60", args{1, 60}, 0},
		{"test 60~120", args{60, 60}, 60},
		{"test 0~1", args{0, 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getSegmentID(tt.args.deliverAt, tt.args.interval), "getSegmentID(%v, %v)", tt.args.deliverAt, tt.args.interval)
		})
	}
}
