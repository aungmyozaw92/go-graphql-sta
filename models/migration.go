package models

import (
	"log"

	"github.com/aungmyozaw92/go-graphql/config"
)

func MigrateTable() {
	db := config.GetDB()

	err := db.AutoMigrate(
		&User{},
	)
	if err != nil {
		log.Fatal(err)
	}
}