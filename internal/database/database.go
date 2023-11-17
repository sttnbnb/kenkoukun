package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open(sqlite.Open("./db/production.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to establish database connection.")
	}

	log.Println("Database connection established.")
}
