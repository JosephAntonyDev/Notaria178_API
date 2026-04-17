package repository

import (
	"database/sql"
)

type PostgresDeviceTokenRepository struct {
	db *sql.DB
}

func NewPostgresDeviceTokenRepository(db *sql.DB) *PostgresDeviceTokenRepository {
	return &PostgresDeviceTokenRepository{
		db: db,
	}
}
