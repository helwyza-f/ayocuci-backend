package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Connect(dbPath string) {
	log.Println("DB PATH:", dbPath)

	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbPath,
	}, &gorm.Config{})

	if err != nil {
		log.Fatal("failed connect database:", err)
	}

	DB = db
}
