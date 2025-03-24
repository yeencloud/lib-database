package database

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/yeencloud/lib-database/domain"
	"github.com/yeencloud/lib-shared/errors"
)

func GetDatabaseFromContext(ctx context.Context) (*gorm.DB, error) {
	dbFromContext := ctx.Value(domain.DatabaseCtxKey)
	if dbFromContext == nil {
		return nil, &errors.ObjectNotInContextError{Object: domain.DatabaseCtxKey}
	}

	db, ok := dbFromContext.(*gorm.DB)
	if !ok {
		return nil, &errors.WrongObjectTypeInContextError{Object: domain.DatabaseCtxKey, ExpectedType: "*gorm.DB"}
	}

	return db, nil
}

func WithLogger(ctx context.Context, logger *log.Entry) context.Context {
	return context.WithValue(ctx, domain.DatabaseCtxKey, logger)
}
