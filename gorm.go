package database

import (
	"github.com/fatih/structs"
	"gorm.io/gorm"

	"github.com/yeencloud/lib-shared"

	"github.com/yeencloud/lib-metrics"
)

type Database struct {
	DB *gorm.DB

	metrics metrics.MetricsInterface
}

func (d *Database) RegisterModels(models ...interface{}) error {
	ctx := shared.Context{}

	err := d.DB.WithContext(&ctx).AutoMigrate(models...)
	if err != nil {
		return err
	}

	if d.metrics != nil {
		for _, model := range models {
			_ = model
			structName := structs.Name(model)
			tableName := d.DB.NamingStrategy.TableName(structName)
			d.metrics.LogPoint(metrics.MetricPoint{
				Name: "DB",
				Tags: map[string]string{
					"table": tableName,
				},
			}, metrics.MetricValues{
				"automigrate": 1,
			})
		}
	}
	return nil
}
