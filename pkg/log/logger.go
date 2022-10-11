package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger new a zap log, default callerSkip is 1
func NewLogger(c *OutputConfig) *zap.Logger {
	return newZapLogWithCallerSkip(c, 1)
}

// newZapLogWithCallerSkip new a zap log
func newZapLogWithCallerSkip(c *OutputConfig, callerSkip int) *zap.Logger {
	if c.Writer == OutputFile {
		// 	TODO: file output
	}
	core := newConsoleCore(c)

	return zap.New(
		zapcore.NewTee(core),
		zap.AddCallerSkip(callerSkip),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.WarnLevel),
	)

}

func newConsoleCore(c *OutputConfig) zapcore.Core {
	level := zap.NewAtomicLevelAt(Levels[c.Level])
	hook := lumberjack.Logger{
		Filename:   "./log/pudding.log", // 日志文件路径
		MaxSize:    128,                 // 每个日志文件保存的大小 单位:M
		MaxAge:     7,                   // 文件最多保存多少天
		MaxBackups: 30,                  // 日志文件最多保存多少个备份
		Compress:   false,               // 是否压缩
	}
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook), zapcore.AddSync(os.Stdout)}
	return zapcore.NewCore(
		newEncoder(c),
		zapcore.NewMultiWriteSyncer(writes...),
		level)
}

func newEncoder(c *OutputConfig) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "name",
		CallerKey:  "caller",
		MessageKey: "message",
		// StacktraceKey: "stacktrace",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalLevelEncoder,
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
