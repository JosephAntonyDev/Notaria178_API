package repository

import "database/sql"

type PostgresDocumentRepository struct {
	db *sql.DB
}

func NewPostgresDocumentRepository(db *sql.DB) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{db: db}
}
