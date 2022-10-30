package logger

import (
	plog "github.com/apache/pulsar-client-go/pulsar/log"
	"go.uber.org/zap"

	"github.com/beihai0xff/pudding/pkg/log"
)

type PulsarLogger struct {
	*zap.SugaredLogger
}

func (p *PulsarLogger) SubLogger(fields plog.Fields) plog.Logger {
	var f []interface{}
	for K, v := range fields {
		f = append(f, K, v)
	}
	return &PulsarLogger{p.With(f...)}
}

func (p *PulsarLogger) WithFields(fields plog.Fields) plog.Entry {

	var f []interface{}
	for K, v := range fields {
		f = append(f, K, v)
	}
	return &PulsarLogger{p.With(f...)}

}

func (p *PulsarLogger) WithField(name string, value interface{}) plog.Entry {
	return &PulsarLogger{p.With(name, value)}
}

func (p *PulsarLogger) WithError(err error) plog.Entry {
	return &PulsarLogger{p.With("error", err)}
}

func GetPulsarLogger() *PulsarLogger {
	return &PulsarLogger{log.NewLogger(&log.OutputConfig{
		Writer:     "",
		Formatter:  log.OutputConsole,
		Level:      "debug",
		CallerSkip: 1,
	}).Sugar().With("pulsar")}
}
