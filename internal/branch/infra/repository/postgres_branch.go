package repository

import "database/sql"

type PostgresBranchRepository struct {
	db *sql.DB
}

func NewPostgresBranchRepository(db *sql.DB) *PostgresBranchRepository {
	return &PostgresBranchRepository{db: db}
}
