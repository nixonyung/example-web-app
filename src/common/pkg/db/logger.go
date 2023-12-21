// (ref.) [go-gorm/gorm - logger/logger.go]https://github.com/go-gorm/gorm/blob/master/logger/logger.go

package db

import (
	"context"
	"errors"
	"io"
	"runtime"
	"strings"
	"time"

	gorm_logger "gorm.io/gorm/logger"
	"source.local/common/pkg/formatter"
	"source.local/common/pkg/logger"
)

// (ref.) [go-gorm/gorm - utils/utils.go - FileWithLineNum](https://github.com/go-gorm/gorm/blob/87decced23be0ce21929fe393fc4fa3a936b1ec8/utils/utils.go#L34)
func locOffsetToQueryCaller() int {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if !strings.HasPrefix(file, "/gomodcache/gorm.io") {
			return i
		}
	}
	return 0
}

type dbLogger struct {
	io.Writer
	gorm_logger.LogLevel
}

var _ gorm_logger.Interface = &dbLogger{} // implements gorm_logger.Interface

func (l *dbLogger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gorm_logger.Info {
		(&logger.Logger{LocOffset: locOffsetToQueryCaller()}).Printf("db [info] "+msg, data...)
	}
}

func (l *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gorm_logger.Warn {
		(&logger.Logger{LocOffset: locOffsetToQueryCaller()}).Printf("db [warn] "+msg, data...)
	}
}

func (l *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gorm_logger.Error {
		(&logger.Logger{LocOffset: locOffsetToQueryCaller()}).Printf("db [error] "+msg, data...)
	}
}

func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gorm_logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	query, numRows := fc()
	if err != nil && l.LogLevel >= gorm_logger.Error && !errors.Is(err, gorm_logger.ErrRecordNotFound) {
		_logger := &logger.Logger{LocOffset: locOffsetToQueryCaller()}
		{
			_logger.Printf("db [%s] [rows:%d] %s",
				formatter.SecondsInEngineeringNotation(elapsed),
				numRows,
				query,
			)
			_logger.Printf("db [error] %v", err)
		}
	} else if l.LogLevel == gorm_logger.Info {
		(&logger.Logger{LocOffset: locOffsetToQueryCaller()}).Printf("db [%s] [rows:%d] %s",
			formatter.SecondsInEngineeringNotation(elapsed),
			numRows,
			query,
		)
	}
}
