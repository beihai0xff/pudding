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
