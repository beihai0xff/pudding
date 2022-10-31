package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/beihai0xff/pudding/pkg/configs"
)

// NewLogger new a zap logger, default callerSkip is 1
func NewLogger(c *configs.OutputConfig) Logger {
	return &logger{newZapLogWithCallerSkip(c).Sugar()}
}

// newLog new a zap log, default callerSkip is 1
func newLog(c *configs.OutputConfig) *logger {
	return &logger{newZapLogWithCallerSkip(c).Sugar()}
}

// newZapLogWithCallerSkip new a zap log
func newZapLogWithCallerSkip(c *configs.OutputConfig) *zap.Logger {
	return zap.New(
		zapcore.NewTee(newCore(c)),
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.WarnLevel),
	)

}

func newCore(c *configs.OutputConfig) zapcore.Core {
	level := zap.NewAtomicLevelAt(configs.Levels[c.Level])

	var writes []zapcore.WriteSyncer
	for _, writer := range c.Writers {
		if writer == configs.OutputConsole {
			writes = append(writes, getConsoleWriter())
		}
		if writer == configs.OutputFile {
			writes = append(writes, getFileWriter(c.FileConfig))
		}
	}

	writes = append(writes, getConsoleWriter())
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
func getFileWriter(c *configs.FileConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.Filename,   // 日志文件路径
		MaxSize:    c.MaxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     c.MaxAge,     // 文件最多保存多少天
		MaxBackups: c.MaxBackups, // 日志文件最多保存多少个备份
		Compress:   c.Compress,   // 是否压缩
	}

	return zapcore.AddSync(lumberJackLogger)
}

func newEncoder(c *configs.OutputConfig) zapcore.Encoder {
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
	if c.Formatter == configs.EncoderTypeJSON {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	return encoder
}
