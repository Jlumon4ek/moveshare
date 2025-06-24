package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type DatabaseSettings struct {
	User     string `env:"POSTGRES_USER,notEmpty"`
	Password string `env:"POSTGRES_PASSWORD,notEmpty"`
	Host     string `env:"POSTGRES_HOST,notEmpty"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	DB       string `env:"POSTGRES_DB,notEmpty"`
}

func (c *DatabaseSettings) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DB)
}

func LoadDatabaseSettings() (*DatabaseSettings, error) {
	_ = godotenv.Load()
	var cfg DatabaseSettings
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
