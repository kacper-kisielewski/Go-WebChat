package db

import (
	"Website/settings"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB struct
var DB *gorm.DB

func init() {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		settings.PostgresHost,
		settings.PostgresUser,
		settings.PostgresPassword,
		settings.PostgresDatabase,
		settings.PostgresPort,
	)
	DB, err = gorm.Open(postgres.Open(dsn), new(gorm.Config))
	if err != nil {
		log.Panic(err)
	}
	DB.AutoMigrate(new(User))
}
