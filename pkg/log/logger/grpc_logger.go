package logger

import (
	"google.golang.org/grpc/grpclog"

	"github.com/beihai0xff/pudding/configs"
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
	l := log.GetLoggerByName(GRPCLoggerName).WithFields("module", "grpc")

	level := configs.Levels[configs.GetLogConfig(BackendLoggerName).Level]
	logger := &GRPCLogger{Logger: l, verbosity: int(level)}

	return logger
}
