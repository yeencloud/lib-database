package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yeencloud/lib-database/domain"
)

func NewPostgresDatabase(config *domain.PostgresConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.Address, config.Username, config.Password.Value, config.Database, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newGormLogger( /*nil*/ ),
	})

	if err != nil {
		return nil, err
	}

	return &Database{
		DB: db,
	}, nil
}
