package logger

import (
	"google.golang.org/grpc/grpclog"

	"github.com/beihai0xff/pudding/pkg/log"
)

type GRPCLogger struct {
	log.Logger
	verbosity int
}

func (l *GRPCLogger) Warning(args ...interface{}) {
	l.Warn(args...)
}

func (l *GRPCLogger) Warningln(args ...interface{}) {
	l.Warnln(args...)
}

func (l *GRPCLogger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func (l *GRPCLogger) V(level int) bool {
	return level < l.verbosity
}

func GetGRPCLogger() grpclog.LoggerV2 {
	l := log.GerLoggerByName(GRPCLoggerName).WithFields("module", "grpc")

	logger := &GRPCLogger{Logger: l}
	grpclog.SetLoggerV2(logger)

	return logger
}
