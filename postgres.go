package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yeencloud/lib-database/domain"
)

func NewPostgresDatabase(config *domain.PostgresConfig) (*Database, error) {
	sslmode := "disable"
	if config.UseTLS {
		sslmode = "require"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.Address, config.Username, config.Password.Value, config.Database, config.Port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newGormLogger(),
	})

	if err != nil {
		return nil, err
	}

	return &Database{
		Gorm: db,
	}, nil
}
