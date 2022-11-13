package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/beihai0xff/pudding/pkg/log"
)

const (
	traceStr     = "[%dms] [rows:%v] %s"
	traceWarnStr = "%s\n[%dms] [rows:%v] %s"
	traceErrStr  = "%s\n[%dms] [rows:%v] %s"
)

// levels gorm logger level
var levels = map[string]logger.LogLevel{
	"":       logger.Info,
	"debug":  logger.Info,
	"info":   logger.Info,
	"warn":   logger.Warn,
	"error":  logger.Error,
	"silent": logger.Silent,
}

type GORMLogger struct {
	l                                   log.Logger
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

func (l *GORMLogger) Trace(c context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	if rows == -1 {
		rows = 0
	}
	switch {
	case err != nil && l.level >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError): //nolint:lll
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
		l:                         log.GerLoggerByName(BackendLoggerName).WithFields("module", "backend"),
		level:                     levels["debug"],
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             1 * time.Second,
		traceStr:                  traceStr,
		traceErrStr:               traceErrStr,
		traceWarnStr:              traceWarnStr,
	}
}
