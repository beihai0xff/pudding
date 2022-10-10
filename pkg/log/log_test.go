package log

import (
	"flag"
	"os"
	"testing"
)

var testLogger Logger

func TestMain(m *testing.M) {
	testLogger = NewLogger(defaultTestConfig)

	flag.Parse()
	exitCode := m.Run()

	os.Exit(exitCode)
}
