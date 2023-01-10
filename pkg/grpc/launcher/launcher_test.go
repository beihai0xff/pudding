package launcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getListen(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name     string
		args     args
		wantAddr string
	}{
		{"test_50050", args{50050}, "[::]:50050"},
		{"test_8081", args{8081}, "[::]:8081"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantAddr, getListen(tt.args.port).Addr().String())
		})
	}

	// test_port_inuse
	assert.Panics(t, func() { getListen(8081) })
}
