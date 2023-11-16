package db

import (
	"database/sql"
	"log"
)

func Connect() *sql.DB {
	dbc, err := sql.Open("sqlite3", "./db/production.db")
	if err != nil {
		log.Fatal(err)
	}

	return dbc
}

