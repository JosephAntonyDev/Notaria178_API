package repository

import "database/sql"

type PostgresActRepository struct {
	db *sql.DB
}

func NewPostgresActRepository(db *sql.DB) *PostgresActRepository {
	return &PostgresActRepository{db: db}
}
