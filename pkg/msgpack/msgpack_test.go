package msgpack

import (
	"reflect"
	"testing"
)

func TestS2Compress(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{"normal_test", args{[]byte("aaa,len<64")}},
		{"large_body_test", args{[]byte("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")}},
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
	}{
		{"normal_test", args{S2Compress([]byte("aaa,len<64"))}, false},
		{"large_body_test", args{S2Compress([]byte("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := S2Decompress(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("S2Decompress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.args.b) {
				t.Errorf("S2Decompress() = %v, want %v", got, tt.args.b)
			}
		})
	}
}
