package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/beihai0xff/pudding/pkg/log"
)

const (
	traceStr     = "[%dms] [rows:%v] %s"
	traceWarnStr = "%s\n[%dms] [rows:%v] %s"
	traceErrStr  = "%s\n[%dms] [rows:%v] %s"
)

type GORMLogger struct {
	l                                   *zap.SugaredLogger
	level                               logger.LogLevel
	IgnoreRecordNotFoundError           bool
	SlowThreshold                       time.Duration
	traceStr, traceErrStr, traceWarnStr string
}

func (l *GORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

func (l *GORMLogger) Info(ctx context.Context, s string, i ...interface{}) {
	l.l.Infof(s, i...)
}

func (l *GORMLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.l.Warnf(s, i...)
}

func (l *GORMLogger) Error(ctx context.Context, s string, i ...interface{}) {
	l.l.Errorf(s, i...)
}

func (l *GORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	if rows == -1 {
		rows = 0
	}
	switch {
	case err != nil && l.level >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.l.Errorf(l.traceErrStr, err, elapsed.Milliseconds(), rows, sql)

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.level >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.l.Warnf(l.traceWarnStr, slowLog, elapsed.Milliseconds(), rows, sql)

	case l.level == logger.Info:
		l.l.Infof(l.traceStr, elapsed.Milliseconds(), rows, sql)

	}
}

func GetGORMLogger() *GORMLogger {
	return &GORMLogger{
		l: log.NewLogger(&log.OutputConfig{
			Writer:     "",
			Formatter:  log.OutputConsole,
			Level:      "debug",
			CallerSkip: 2,
		}).Sugar().With("module", "backend"),
		level:                     logger.Info,
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             1 * time.Second,
		traceStr:                  traceStr,
		traceErrStr:               traceErrStr,
		traceWarnStr:              traceWarnStr,
	}
}
