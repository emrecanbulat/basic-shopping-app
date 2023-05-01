package model

import (
	"errors"
	"log"
	"shoppingApp/internal/client"
)

var ErrRecordNotFound = errors.New("record not found")

func Migrate() {
	err := client.PostgreSqlClient.Migrator().AutoMigrate(
		&Product{},
		&User{},
		&Token{},
		&Order{},
	)

	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Migration successfully completed")
}
