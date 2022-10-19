package random

import (
	"fmt"
	"testing"
)

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
				fmt.Println(got)
			}
		})
	}
}
