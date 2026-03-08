package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:chocolate37900@localhost:5432/notaria178_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	content, err := ioutil.ReadFile("../../migrate_acts.sql")
	if err != nil {
		log.Fatal("Error reading script:", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal("Error executing script:", err)
	}

	fmt.Println("Migración exitosa!")
}
