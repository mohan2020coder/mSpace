// internal/db/db.go
package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(dsn string) *gorm.DB {
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	return gdb
}
