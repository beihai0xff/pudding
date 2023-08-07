package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/pkg/log"
)

func TestGetLogConfig(t *testing.T) {
	assert.NoError(t, Parse("../test/data/config.format.yaml", ConfigFormatYAML, ReadFromFile))

	// happy_path
	logConf := GetLogConfig(log.DefaultLoggerName)
	assert.Equal(t, &log.Config{
		LogName:    "default",
		Writers:    []string{log.OutputConsole},
		Level:      "debug",
		Format:     log.EncoderTypeJSON,
		CallerSkip: 1,
		FileConfig: log.FileConfig{},
	},
		logConf)
}
