package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/beihai0xff/pudding/configs"
)

// newLogger new a zap log, default callerSkip is 1
func newLogger(c *configs.LogConfig) *logger {
	return &logger{newZapLogWithCallerSkip(c).Sugar()}
}

// newZapLogWithCallerSkip new a zap log
func newZapLogWithCallerSkip(c *configs.LogConfig) *zap.Logger {
	return zap.New(
		zapcore.NewTee(newCore(c)),
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.WarnLevel),
	)
}

func newCore(c *configs.LogConfig) zapcore.Core {
	level := zap.NewAtomicLevelAt(configs.Levels[c.Level])

	// get log output writer
	var writes []zapcore.WriteSyncer
	for _, writer := range c.Writers {
		if writer == configs.OutputConsole {
			writes = append(writes, getConsoleWriter())
		}
		if writer == configs.OutputFile {
			writes = append(writes, getFileWriter(&c.FileConfig))
		}
	}

	return zapcore.NewCore(
		newEncoder(c),
		zapcore.NewMultiWriteSyncer(writes...),
		level)
}

// getConsoleWriter write log to console
func getConsoleWriter() zapcore.WriteSyncer {
	return os.Stdout
}

// getFileWriter write log to file
func getFileWriter(c *configs.LogFileConfig) zapcore.WriteSyncer {
	if c.Filepath == "" {
		Fatalf("log file writer set, but log file path is empty, please check your config")
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.Filepath,   // 日志文件路径
		MaxSize:    c.MaxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     c.MaxAge,     // 文件最多保存多少天
		MaxBackups: c.MaxBackups, // 日志文件最多保存多少个备份
		Compress:   c.Compress,   // 是否压缩
	}

	return zapcore.AddSync(lumberJackLogger)
}

func newEncoder(c *configs.LogConfig) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "name",
		CallerKey:  "caller",
		MessageKey: "message",
		// StacktraceKey: "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	if c.Format == configs.EncoderTypeJSON {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	return encoder
}
