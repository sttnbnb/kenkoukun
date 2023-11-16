package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {
	dbc, err := sql.Open("sqlite3", "./db/production.db")
	if err != nil {
		log.Fatal(err)
	}

	return dbc
}

func CreateTables(dbc *sql.DB) {
	cmd := `
		CREATE TABLE IF NOT EXISTS kenkou_settings(
			guild_id STRING NOT NULL PRIMARY KEY,
			channel_id STRING,
			time TIMESTAMP NOT NULL
		)
	`
	_, err := dbc.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}
}
