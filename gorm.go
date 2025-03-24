package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	sharedLogger "github.com/yeencloud/lib-shared/log"
)

type Database struct {
	Gorm *gorm.DB
}

func (d *Database) RegisterModels(models ...interface{}) error {
	ctx := context.Background()

	logger := logrus.NewEntry(logrus.StandardLogger())
	ctx = sharedLogger.WithLogger(ctx, logger)

	return d.Gorm.WithContext(ctx).AutoMigrate(models...)
}
