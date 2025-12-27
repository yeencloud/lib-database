package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/yeencloud/lib-database/domain"
	"github.com/yeencloud/lib-shared/apperr"
)

func GetDatabaseFromContext(ctx context.Context) (*gorm.DB, error) {
	dbFromContext := ctx.Value(domain.DatabaseCtxKey)
	if dbFromContext == nil {
		return nil, &apperr.ObjectNotInContextError{Object: domain.DatabaseCtxKey}
	}

	db, ok := dbFromContext.(*gorm.DB)
	if !ok {
		return nil, &apperr.WrongObjectTypeInContextError{Object: domain.DatabaseCtxKey, ExpectedType: "*gorm.DB"}
	}

	return db, nil
}
