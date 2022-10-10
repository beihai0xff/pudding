package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(defaultTestConfig)
	assert.NotNil(t, logger)
}

func Test_ZapLog(t *testing.T) {
	testLogger.Trace("hello world")
	testLogger.Debug("hello world")
	testLogger.Info("hello world")
	testLogger.Warn("hello world")
	testLogger.Error("hello world")

	puff := "puff"
	testLogger.Tracef("hello world %s", puff)
	testLogger.Debugf("hello world %s", puff)
	testLogger.Infof("hello world %s", puff)
	testLogger.Warnf("hello world %s", puff)
	testLogger.Errorf("hello world %s", puff)

	testLogger.WithFields("field", "testfield").Debug("testdebug")
}
