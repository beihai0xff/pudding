package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogConfig(t *testing.T) {
	Parse("../test/data/config.test.yaml", ConfigFormatYAML, ReadFromFile)

	// happy_path
	logConf := GetLogConfig("default")
	assert.Equal(t, &LogConfig{
		LogName:    "default",
		Writers:    []string{OutputConsole},
		Level:      "debug",
		Format:     EncoderTypeJSON,
		CallerSkip: 1,
		FileConfig: LogFileConfig{},
	},
		logConf)
}
