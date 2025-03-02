package database

import (
	"context"
	"time"

	"gorm.io/gorm/logger"

	"github.com/yeencloud/lib-shared"
	"github.com/yeencloud/lib-shared/log"

	"github.com/yeencloud/lib-logger"
	"github.com/yeencloud/lib-logger/domain"
	"github.com/yeencloud/lib-metrics"
)

var (
	LogFieldSQL        = log.Path{Identifier: "sql"}
	LogFieldSQLRequest = log.Path{Parent: &LogFieldSQL, Identifier: "request"}

	LogFieldSQLRowsAffected = log.Path{Parent: &LogFieldSQLRequest, Identifier: "rows_affected"}
	LogFieldSQLQuery        = log.Path{Parent: &LogFieldSQLRequest, Identifier: "query"}

	LogScopeTime        = log.Path{Parent: &LogFieldSQLRequest, Identifier: "time"}
	LogFieldTimeStarted = log.Path{Parent: &LogScopeTime, Identifier: "started"}
	LogFieldTimeEnded   = log.Path{Parent: &LogScopeTime, Identifier: "ended"}
	LogFieldDuration    = log.Path{Parent: &LogScopeTime, Identifier: "duration_ms"}
)

type gormLogger struct {
	m metrics.MetricsInterface
}

func (g gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g gormLogger) log(ctx *shared.Context, level LoggerDomain.Level, s string) {
	if ctx == nil {
		return
	}

	message := Logger.Log(level)

	fields := log.Fields{}

	message.WithFields(fields).Msg(s)
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, affectedRows int64), err error) {
	sql, affectedRows := fc()
	sharedcontext, ok := ctx.(*shared.Context)
	if !ok {
		Logger.Log(LoggerDomain.LogLevelWarn).Msg("attempted to execute request without context, metrics will be lost.")
		return
	}

	end := time.Now()
	duration := time.Duration(end.UnixMilli() - begin.UnixMilli()).Milliseconds()

	arrayOfRequest := map[string]interface{}{
		LogFieldSQLQuery.MetricKey():    sql,
		LogFieldTimeStarted.MetricKey(): begin,
		LogFieldTimeEnded.MetricKey():   end,
		LogFieldDuration.MetricKey():    duration,
	}

	if affectedRows > 0 {
		arrayOfRequest[LogFieldSQLRowsAffected.MetricKey()] = affectedRows
	}

	if err != nil {
		arrayOfRequest[LoggerDomain.LogFieldError.MetricKey()] = err

	}
	sharedcontext.WithLog(arrayOfRequest)

	g.log(sharedcontext, LoggerDomain.LogLevelSQL, sql)
}

func newGormLogger(m metrics.MetricsInterface) *gormLogger {
	return &gormLogger{
		m: m,
	}
}

// Those fields are required for the implementation of the correct interface. However, they are not used in the current implementation.
func (g gormLogger) Info(context.Context, string, ...interface{})  {}
func (g gormLogger) Warn(context.Context, string, ...interface{})  {}
func (g gormLogger) Error(context.Context, string, ...interface{}) {}
