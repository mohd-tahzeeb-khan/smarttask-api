package models

import (
	"log"

	"github.com/smarttask/api/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

package models

import (
"log"

"github.com/smarttask/api/internal/config"
"gorm.io/driver/sqlite"
"gorm.io/gorm"
"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	DB, err = gorm.Open(sqlite.Open(config.App.DBPath), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB.Exec("PRAGMA journal_mode=WAL;")
	DB.Exec("PRAGMA foreign_keys=ON;")

	if err := DB.AutoMigrate(&User{}, &Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("✅ Database connected and migrated")
}