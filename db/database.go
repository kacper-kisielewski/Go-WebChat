package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB struct
var DB *gorm.DB

func init() {
	var err error

	dsn := "host=127.0.0.1 user=postgres password=password dbname=website port=5432"
	DB, err = gorm.Open(postgres.Open(dsn), new(gorm.Config))
	if err != nil {
		log.Panic(err)
	}
	DB.AutoMigrate(new(User))
}
