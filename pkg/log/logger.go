package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// newLogger new a zap log, default callerSkip is 1
func newLogger(c *Config) *logger {
	return &logger{newZapLogWithCallerSkip(c).Sugar()}
}

// newZapLogWithCallerSkip new a zap log
func newZapLogWithCallerSkip(c *Config) *zap.Logger {
	return zap.New(
		zapcore.NewTee(newCore(c)),
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.WarnLevel),
	)
}

func newCore(c *Config) zapcore.Core {
	level := zap.NewAtomicLevelAt(Levels[c.Level])

	// get log output writer
	var writes []zapcore.WriteSyncer

	for _, writer := range c.Writers {
		if writer == OutputConsole {
			writes = append(writes, getConsoleWriter())
		}

		if writer == OutputFile {
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
func getFileWriter(c *FileConfig) zapcore.WriteSyncer {
	if c.Filepath == "" {
		Fatalf("log file writer set, but log file path is empty, please check your config")
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.Filepath,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		Compress:   c.Compress,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func newEncoder(c *Config) zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	if c.Format == EncoderTypeJSON {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	return encoder
}
