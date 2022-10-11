package log

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	flag.Parse()
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger(defaultConfig)
	assert.NotNil(t, logger)
}

func Test_ZapLog(t *testing.T) {
	Debugf("hello world")
	Infof("hello world")
	Warnf("hello world")
	Errorf("hello world")

	pudding := "pudding"
	Debugf("hello world %s", pudding)
	Infof("hello world %s", pudding)
	Warnf("hello world %s", pudding)
	Errorf("hello world %s", pudding)

	WithFields("field", "testfield").Debug("testdebug")
}
