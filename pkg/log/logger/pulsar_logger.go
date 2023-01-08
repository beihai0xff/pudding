package logger

import (
	plog "github.com/apache/pulsar-client-go/pulsar/log"

	"github.com/beihai0xff/pudding/pkg/log"
)

type PulsarLogger struct {
	log.Logger
	with func(fields ...interface{}) log.Logger
}

func (p *PulsarLogger) SubLogger(fields plog.Fields) plog.Logger {
	f := make([]interface{}, 0, 2*len(fields))

	for K, v := range fields {
		f = append(f, K, v)
	}
	return &PulsarLogger{Logger: p.with(f...), with: p.with}
}

func (p *PulsarLogger) WithFields(fields plog.Fields) plog.Entry {
	f := make([]interface{}, 0, 2*len(fields))

	for K, v := range fields {
		f = append(f, K, v)
	}

	return &PulsarLogger{Logger: p.with(f...), with: p.with}
}

func (p *PulsarLogger) WithField(name string, value interface{}) plog.Entry {
	return &PulsarLogger{Logger: p.with(name, value), with: p.with}
}

func (p *PulsarLogger) WithError(err error) plog.Entry {
	return &PulsarLogger{Logger: p.with("error", err), with: p.with}
}

func GetPulsarLogger() plog.Logger {
	l := log.GetLoggerByName(PulsarLoggerName).WithFields("module", "pulsar")
	return &PulsarLogger{Logger: l, with: l.WithFields}
}
