package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 定义 logger 接口
type Logger interface {
	// Trace logs to TRACE log. Arguments are handled in the manner of fmt.Print.
	Trace(args ...interface{})
	// Tracef logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
	Tracef(format string, args ...interface{})
	// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Print.
	Debug(args ...interface{})
	// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
	Debugf(format string, args ...interface{})
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// Warn logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warn(args ...interface{})
	// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warnf(format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// all Fatal logs will exit with os.Exit(1).
	Fatal(args ...interface{})
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	// all Fatal logs will exit with os.Exit(1).
	Fatalf(format string, args ...interface{})

	// Sync calls the underlying Core's Sync method, flushing any buffered log entries.
	// Applications should take care to call Sync before exiting
	Sync() error

	// WithFields set some custom key-value fields
	// Do not use log.WithFields("k", "v").WithFields("k1", "v1").Debug("hello")
	WithFields(fields ...string) *zap.Logger
}

// NewLogger new a zap log, default callerSkip is 2
func NewLogger(c *OutputConfig) Logger {
	return newZapLogWithCallerSkip(c, 2)
}

// newZapLogWithCallerSkip new a zap log
func newZapLogWithCallerSkip(c *OutputConfig, callerSkip int) Logger {
	if c.Writer == OutputFile {
		// 	TODO: file output
	}
	core, zapLevel := newConsoleCore(c)

	logger := zap.New(
		zapcore.NewTee(core),
		zap.AddCallerSkip(callerSkip),
		zap.AddCaller(),
	)

	return &zapLog{
		level:  zapLevel,
		logger: logger,
	}
}

func newConsoleCore(c *OutputConfig) (zapcore.Core, zap.AtomicLevel) {
	level := zap.NewAtomicLevelAt(Levels[c.Level])
	return zapcore.NewCore(
		newEncoder(c),
		zapcore.Lock(os.Stdout),
		level), level
}

func newEncoder(c *OutputConfig) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "name",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		// TODO: custom EncodeTime
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	if c.Formatter == EncoderTypeJson {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	return encoder
}
