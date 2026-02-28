package core

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetDBPool() (*sql.DB, error) {

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("la variable de entorno DB_URL está vacía")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error al verificar la conexión (ping): %w", err)
	}

	fmt.Println("Conexión a PostgreSQL exitosa (Pool configurado)")
	return db, nil
}