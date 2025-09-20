package domain

import (
	"github.com/yeencloud/lib-shared/config"
)

type PostgresConfig struct {
	Address string `config:"PG_ADDRESS" default:"localhost"`
	Port    int    `config:"PG_PORT" default:"5432"`

	Username string        `config:"PG_USERNAME" default:"postgres"`
	Password config.Secret `config:"PG_PASSWORD" default:"postgres"`
	Database string        `config:"PG_DATABASE" default:"postgres"`

	UseTLS bool `config:"PG_TLS" default:"true"`
}
