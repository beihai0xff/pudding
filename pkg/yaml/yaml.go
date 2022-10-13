package yaml

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func Parse(filePath string, data interface{}) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, data)
}
