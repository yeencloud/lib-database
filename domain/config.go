package domain

type DatabaseConfig struct {
	// Bind Address
	Engine string `config:"DB_ENGINE" default:"POSTGRES"`
}
