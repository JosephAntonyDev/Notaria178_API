package repository

import "database/sql"

type PostgresWorkRepository struct {
	db *sql.DB
}

func NewPostgresWorkRepository(db *sql.DB) *PostgresWorkRepository {
	return &PostgresWorkRepository{db: db}
}
