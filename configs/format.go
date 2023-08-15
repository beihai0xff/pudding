package configs

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/knadh/koanf/v2"
)

// UnmarshalToStruct unmarshal config to struct
func UnmarshalToStruct(path string, c any) error {
	if err := k.UnmarshalWithConf(path, c, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// JSONFormat returns the json format of the config.
func JSONFormat(c any) (*bytes.Buffer, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("marshal config %v failed: %w", c, err)
	}

	var buf bytes.Buffer
	if err = json.Indent(&buf, b, "", "    "); err != nil {
		return nil, fmt.Errorf("indent config %v failed: %w", c, err)
	}

	return &buf, nil
}
