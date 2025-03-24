package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	"github.com/yeencloud/lib-database/domain"
	metrics "github.com/yeencloud/lib-metrics"
	MetricsDomain "github.com/yeencloud/lib-metrics/domain"
	sharedLogger "github.com/yeencloud/lib-shared/log"
	sharedMetrics "github.com/yeencloud/lib-shared/metrics"
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

	logger.Info(sql)

	var mPoint MetricsDomain.Point
	var mValues MetricsDomain.Values

	point, ok := ctx.Value(sharedMetrics.MetricsPointKey).(MetricsDomain.Point)
	if !ok {
		mPoint = MetricsDomain.Point{
			Tags: map[string]string{},
		}
	} else {
		mPoint = point
	}

	mPoint.Name = domain.SQLMetricPointName

	values, ok := ctx.Value(sharedMetrics.MetricsValuesKey).(MetricsDomain.Values)
	if !ok {
		mValues = MetricsDomain.Values{}
	} else {
		mValues = values
	}
	mValues[domain.LogFieldSQLQuery.MetricKey()] = sql
	if affectedRows > 0 {
		mValues[domain.LogFieldSQLRowsAffected.MetricKey()] = affectedRows
	}
	mValues[domain.LogFieldDuration.MetricKey()] = duration

	metrics.LogPoint(mPoint, mValues)
}

func newGormLogger() *gormLogger {
	return &gormLogger{}
}

// Those fields are required for the implementation of the correct interface. However, they are not used in the current implementation.
func (g gormLogger) Info(context.Context, string, ...interface{})  {}
func (g gormLogger) Warn(context.Context, string, ...interface{})  {}
func (g gormLogger) Error(context.Context, string, ...interface{}) {}
