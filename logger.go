package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	metrics "github.com/yeencloud/lib-metrics"
	MetricsDomain "github.com/yeencloud/lib-metrics/domain"
	"github.com/yeencloud/lib-shared/namespace"
)

var (
	LogFieldSQL        = namespace.Namespace{Identifier: "sql"}
	LogFieldSQLRequest = namespace.Namespace{Parent: &LogFieldSQL, Identifier: "request"}

	LogFieldSQLRowsAffected = namespace.Namespace{Parent: &LogFieldSQLRequest, Identifier: "rows_affected"}
	LogFieldSQLQuery        = namespace.Namespace{Parent: &LogFieldSQLRequest, Identifier: "query"}

	LogScopeTime        = namespace.Namespace{Parent: &LogFieldSQLRequest, Identifier: "time"}
	LogFieldTimeStarted = namespace.Namespace{Parent: &LogScopeTime, Identifier: "started"}
	LogFieldTimeEnded   = namespace.Namespace{Parent: &LogScopeTime, Identifier: "ended"}
	LogFieldDuration    = namespace.Namespace{Parent: &LogScopeTime, Identifier: "duration_ms"}
)

type gormLogger struct {
}

func (g gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, affectedRows int64), err error) {
	sql, affectedRows := fc()
	end := time.Now()
	duration := time.Duration(end.UnixMilli() - begin.UnixMilli()).Milliseconds()

	if ctx == nil {
		log.WithField("request", sql).Warn("gorm logger called without context")
		return
	}
	logger := ctx.Value("logger")
	if logger == nil {
		log.WithField("request", sql).Warn("gorm logger called without logger")
		return
	}
	logMessage, ok := logger.(*log.Entry)
	if !ok {
		logMessage = log.NewEntry(log.StandardLogger())
	}

	logMessage = logMessage.WithField("latency", duration)
	logMessage = logMessage.WithField("affectedRows", affectedRows)

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
		logMessage = logMessage.WithError(err)
	}

	logMessage.Info(sql)

	err = metrics.LogPoint(MetricsDomain.Point{
		Name: "sql",
	}, MetricsDomain.Values{
		"request_query":            sql,
		"request_rows_affected":    affectedRows,
		"request_time_duration_ms": duration,
	})

	if err != nil {
		log.WithError(err).Error("failed to log metrics")
	}
}

func newGormLogger() *gormLogger {
	return &gormLogger{}
}

// Those fields are required for the implementation of the correct interface. However, they are not used in the current implementation.
func (g gormLogger) Info(context.Context, string, ...interface{})  {}
func (g gormLogger) Warn(context.Context, string, ...interface{})  {}
func (g gormLogger) Error(context.Context, string, ...interface{}) {}
