package domain

import (
	"github.com/yeencloud/lib-shared/config"
)

type PostgresConfig struct {
	Address string `config:"DB_ADDRESS" default:"localhost"`
	Port    int    `config:"DB_PORT" default:"5432"`

	Username string        `config:"DB_USERNAME" default:"postgres"`
	Password config.Secret `config:"DB_PASSWORD" default:"postgres"`
	Database string        `config:"DB_DATABASE" default:"postgres"`
}
