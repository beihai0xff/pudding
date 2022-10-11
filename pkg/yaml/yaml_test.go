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
	var u []users
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Parse(tt.args.filePath, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(u, tt.want) {
				t.Errorf("Parse() got = %v, want %v", u, tt.want)
			}
		})
	}
}
