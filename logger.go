package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	metrics "github.com/yeencloud/lib-metrics"
	"gorm.io/gorm/logger"

	"github.com/yeencloud/lib-database/domain"
	sharedLogger "github.com/yeencloud/lib-shared/log"
)

type gormLogger struct {
}

func (g gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

type SQLEntryMetric struct {
	Query        string `metric:"query"`
	AffectedRows int64  `metric:"affected_rows"`
	Duration     int64  `metric:"duration"`
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, affectedRows int64), err error) {
	sql, affectedRows := fc()
	end := time.Now()
	duration := time.Duration(end.UnixMilli() - begin.UnixMilli()).Milliseconds()

	if ctx == nil {
		log.WithField(domain.LogFieldSQLQuery.MetricKey(), sql).Warn("gorm logger called without context")
		return
	}

	logger := sharedLogger.GetLoggerFromContext(ctx)
	logger = logger.WithField(domain.LogFieldDuration.MetricKey(), duration)
	logger = logger.WithField(domain.LogFieldSQLRowsAffected.MetricKey(), affectedRows)

	arrayOfRequest := map[string]interface{}{
		domain.LogFieldSQLQuery.MetricKey():    sql,
		domain.LogFieldTimeStarted.MetricKey(): begin,
		domain.LogFieldTimeEnded.MetricKey():   end,
		domain.LogFieldDuration.MetricKey():    duration,
	}

	if affectedRows > 0 {
		arrayOfRequest[domain.LogFieldSQLRowsAffected.MetricKey()] = affectedRows
	}

	if err != nil {
		logger = logger.WithError(err)
	}

	logger.Debug(sql)

	metric := SQLEntryMetric{
		Query:    sql,
		Duration: duration,
	}

	if affectedRows > 0 {
		metric.AffectedRows = affectedRows
	}

	_ = metrics.WritePoint(ctx, domain.SQLMetricPointName, metric)
}

func newGormLogger() *gormLogger {
	return &gormLogger{}
}

// Those fields are required for the implementation of the correct interface. However, they are not used in the current implementation.
func (g gormLogger) Info(context.Context, string, ...interface{})  {}
func (g gormLogger) Warn(context.Context, string, ...interface{})  {}
func (g gormLogger) Error(context.Context, string, ...interface{}) {}
