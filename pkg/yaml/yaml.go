package yaml

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Users struct {
	Name    string   `yaml:"name"`
	Age     int8     `yaml:"age"`
	Address string   `yaml:"address"`
	Hobby   []string `yaml:"hobby"`
}

func Parse(filePath string, data interface{}) error {
	file, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, data)
}
