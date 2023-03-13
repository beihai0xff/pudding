package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOutBoundIP(t *testing.T) {
	tests := []struct {
		name string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEqual(t, "", GetOutBoundIP())
			println(GetOutBoundIP())
		})
	}
}

func TestGetRand(t *testing.T) {
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
	}{
		{"normal_test", args{100, 200}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.args.end; i++ {
				got := GetRand(tt.args.start, tt.args.end)
				if got < tt.args.start || got > tt.args.end {
					t.Errorf("got should less than %d and greater than %d, but got %d", tt.args.end, tt.args.start, got)
				}
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name   string
		setEnv string
		want   string
	}{
		{"empty_env", "", "dev"},
		{"dev_env", "dev", "dev"},
		{"test_env", "test", "test"},
		{"prod_env", "prod", "prod"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("PUDDING_ENV", tt.setEnv)
			assert.Equal(t, tt.want, GetEnv())
		})
	}
	_ = os.Unsetenv("PUDDING_ENV")
}

func TestGetHealthEndpointPath(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// {"normal_test", args{"/api"}, "/api/healthz"},
		{"empty_prefix", args{""}, "/healthz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetHealthEndpointPath(tt.args.prefix), "GetHealthEndpointPath(%v)", tt.args.prefix)
		})
	}
}

func TestGetSwaggerEndpointPath(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal_test", args{"/api"}, "/api/swagger"},
		{"empty_prefix", args{""}, "/swagger"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetSwaggerEndpointPath(tt.args.prefix), "GetSwaggerEndpointPath(%v)", tt.args.prefix)
		})
	}
}
