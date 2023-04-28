package model

import (
	"log"
	"shoppingApp/internal/client"
)

func Migrate() {
	err := client.PostgreSqlClient.Migrator().AutoMigrate(
		&Product{},
	)

	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Migration successfully completed")
}
