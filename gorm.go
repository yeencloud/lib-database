package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	shared "github.com/yeencloud/lib-shared/log"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) RegisterModels(models ...interface{}) error {
	ctx := context.Background()

	ctx = context.WithValue(ctx, shared.ContextLoggerKey, logrus.NewEntry(logrus.StandardLogger())) //nolint:staticcheck

	return d.DB.WithContext(ctx).AutoMigrate(models...)
}
