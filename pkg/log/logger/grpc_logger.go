// Package logger defines logging for components in the pudding project.
// grpc_logger.go defines logging for grpc.
package logger

import (
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc/grpclog"

	"github.com/beihai0xff/pudding/pkg/log"
)

// GRPCLogger is a wrapper of log.Logger to implement grpclog.LoggerV2.
type GRPCLogger struct {
	log.Logger
	verbosity int
}

// Warning logs to the WARNING log.
func (l *GRPCLogger) Warning(args ...any) {
	l.Warn(args...)
}

// Warningln logs to the WARNING log.
func (l *GRPCLogger) Warningln(args ...any) {
	l.Warnln(args...)
}

// Warningf logs to the WARNING log.
func (l *GRPCLogger) Warningf(format string, args ...any) {
	l.Warnf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (l *GRPCLogger) V(level int) bool {
	return level <= l.verbosity
}

var (
	grpcLogOnce sync.Once
	grpcLogger  *GRPCLogger
	_           grpclog.LoggerV2 = (*GRPCLogger)(nil)
)

// GetGRPCLogger returns a grpclog.LoggerV2 that uses the given pudding logger.
func GetGRPCLogger() grpclog.LoggerV2 {
	grpcLogOnce.Do(func() {
		l := log.GetLoggerByName(GRPCLoggerName).WithFields("module", "grpc")

		// default verbosity is 2.
		v := 2
		// Get verbosity from environment variable.
		vLevel := os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")
		if vl, err := strconv.Atoi(vLevel); err == nil {
			v = vl
		}

		grpcLogger = &GRPCLogger{Logger: l, verbosity: v}

		grpclog.SetLoggerV2(grpcLogger)
	})

	return grpcLogger
}
