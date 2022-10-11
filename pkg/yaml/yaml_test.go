package yaml

import (
	"reflect"
	"testing"
)

type users struct {
	Name    string   `yaml:"name"`
	Time    int      `yaml:"time"`
	Address string   `yaml:"address"`
	Hobby   []string `yaml:"hobby"`
}

func TestParse(t *testing.T) {
	var u, u2 []users
	type args struct {
		filePath string
		data     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []users
	}{
		{"normal_test", args{"test.yaml", &u}, false, []users{
			{"pudding",
				2022,
				"beijing",
				[]string{"redis://default:default@192.168.10.117:6379/11", "debug"},
			},
			{"pudding2",
				2022,
				"beijing",
				[]string{"redis://default:default@192.168.10.117:6379/11", "debug"},
			},
		}},
		{"wrong_yaml_file_path", args{"wrong.yaml", &u2}, true, []users{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Parse(tt.args.filePath, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !reflect.DeepEqual(tt.args.data, &tt.want) {
				t.Errorf("Parse() got = %v, want %v", tt.args.data, tt.want)
			}
		})
	}
}
