// Package logger defines logging for components in the pudding project.
// pulsar_logger.go defines logging for pulsar.
package logger

import (
	plog "github.com/apache/pulsar-client-go/pulsar/log"

	"github.com/beihai0xff/pudding/pkg/log"
)

// PulsarLogger is a wrapper of log.Logger to implement pulsar.Logger.
type PulsarLogger struct {
	log.Logger
	with func(fields ...interface{}) log.Logger
}

// SubLogger returns a sub logger with the given Fields.
func (p *PulsarLogger) SubLogger(fields plog.Fields) plog.Logger {
	f := make([]interface{}, 0, 2*len(fields))

	for K, v := range fields {
		f = append(f, K, v)
	}

	return &PulsarLogger{Logger: p.with(f...), with: p.with}
}

// WithFields returns a new Entry with the fields added to it.
func (p *PulsarLogger) WithFields(fields plog.Fields) plog.Entry {
	f := make([]interface{}, 0, 2*len(fields))

	for K, v := range fields {
		f = append(f, K, v)
	}

	return &PulsarLogger{Logger: p.with(f...), with: p.with}
}

// WithField returns a new Entry with the field added to it.
func (p *PulsarLogger) WithField(name string, value interface{}) plog.Entry {
	return &PulsarLogger{Logger: p.with(name, value), with: p.with}
}

// WithError returns a new Entry with the field added to it.
func (p *PulsarLogger) WithError(err error) plog.Entry {
	return &PulsarLogger{Logger: p.with("error", err), with: p.with}
}

// GetPulsarLogger returns a pulsar.Logger that uses the given pudding logger.
func GetPulsarLogger() plog.Logger {
	l := log.GetLoggerByName(PulsarLoggerName).WithFields("module", "pulsar")
	return &PulsarLogger{Logger: l, with: l.WithFields}
}
