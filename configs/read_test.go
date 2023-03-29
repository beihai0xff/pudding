package configs

import (
	"testing"

	kyaml "github.com/knadh/koanf/parsers/yaml"
)

func TestInit(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{"../test/data/config.test.yaml"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReadFromFile(tt.args.filePath, kyaml.Parser())
		})
	}
}
