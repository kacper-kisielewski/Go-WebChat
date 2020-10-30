package db

import "gorm.io/gorm"

//User struct
type User struct {
	gorm.Model
	Username       string `gorm:"unique; not null"`
	Email          string `gorm:"unique; not null"`
	HashedPassword []byte `gorm:"not null"`
}
