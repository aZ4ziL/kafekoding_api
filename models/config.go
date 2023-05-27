package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	db := DB()
	db.AutoMigrate(&User{}, &UserBio{})
}

func DB() *gorm.DB {
	dsn := "host=localhost port=5432 user=fajhri password=root dbname=kafekoding_api sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
