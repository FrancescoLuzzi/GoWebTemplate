package db

import (
	"fmt"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver for sqlx
	"github.com/jmoiron/sqlx"
)

const (
	driverPgx string = "pgx"
)

// Open establishes a connection to the database based on the configuration.
func Open(cfg config.DbConfig) (*sqlx.DB, error) {
	db, err := newSqlx(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func newSqlx(cfg config.DbConfig) (*sqlx.DB, error) {
	return sqlx.Open(driverPgx, cfg.DSN())
}
