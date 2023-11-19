package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Setup initializes the database instance.
func Setup() {
	var err error
	db, err = gorm.Open(sqlite.Open("data/db/gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

// GetReadDB returns read database connection.
func GetWriteDB() *gorm.DB {
	return db
}
