package domain

type DatabaseConfig struct {
	Engine string `config:"DB_ENGINE" default:"POSTGRES"`
}
