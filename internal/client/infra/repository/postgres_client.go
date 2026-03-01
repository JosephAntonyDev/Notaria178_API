package repository

import "database/sql"

type PostgresClientRepository struct {
	db *sql.DB
}

func NewPostgresClientRepository(db *sql.DB) *PostgresClientRepository {
	return &PostgresClientRepository{db: db}
}
