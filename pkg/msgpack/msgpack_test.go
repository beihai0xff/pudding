package msgpack

import (
	"reflect"
	"testing"
)

var (
	testDataL64 = []byte("aaa,len<64")
	testDataG64 = []byte("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
)

func TestS2Compress(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{"normal_test", args{testDataL64}},
		{"large_body_test", args{testDataG64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := S2Compress(tt.args.data)
			res, _ := S2Decompress(got)
			if !reflect.DeepEqual(tt.args.data, res) {
				t.Errorf("S2Compress() = %v, res = %v", got, res)
			}
		})
	}
}

func TestS2Decompress(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []byte
	}{
		{"normal_test", args{S2Compress(testDataL64)}, false, testDataL64},
		{"large_body_test", args{S2Compress(testDataG64)}, false, testDataG64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := S2Decompress(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("S2Decompress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S2Decompress() = %v, want %v", got, tt.args.b)
			}
		})
	}
}
