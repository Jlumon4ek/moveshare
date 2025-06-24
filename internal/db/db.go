package db

import (
	"database/sql"
	"fmt"
	"moveshare/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(cfg *config.DatabaseSettings) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	fmt.Printf("Attempting to connect with: host=%s port=%d user=%s dbname=%s\n",
		cfg.Host, cfg.Port, cfg.User, cfg.DB)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}
