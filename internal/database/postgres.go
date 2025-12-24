package database

import (
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Connect(dbPath string) {
	log.Println("DB PATH:", dbPath)

	var dialector gorm.Dialector

	// 1. Cek apakah path mengandung format postgres
	if strings.Contains(dbPath, "postgres://") {
		log.Println("Connecting to PostgreSQL...")
		dialector = postgres.Open(dbPath)
	} else {
		// 2. Jika bukan postgres, gunakan setup asli kamu (SQLite)
		log.Println("Connecting to SQLite...")
		dialector = sqlite.Dialector{
			DriverName: "sqlite",
			DSN:        dbPath,
		}
	}

	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Fatal("failed connect database:", err)
	}

	DB = db
}