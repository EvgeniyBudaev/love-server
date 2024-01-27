package app

import (
	"database/sql"
	"fmt"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	_ "github.com/lib/pq"

	"github.com/EvgeniyBudaev/love-server/internal/config"
)

type Database struct {
	logger logger.Logger
	psql   *sql.DB
}

func NewDatabase(logger logger.Logger, psql *sql.DB) *Database {
	return &Database{
		logger: logger,
		psql:   psql,
	}
}

func newPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSlMode,
	)

	return sql.Open("postgres", databaseURL)
}
