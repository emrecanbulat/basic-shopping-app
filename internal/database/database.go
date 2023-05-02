package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// GetPostgreSqlConnection return a PostgreSQL connection as a gorm DB object
func GetPostgreSqlConnection(postgresHost, postgresUser, postgresPass, postgresDb, postgresPort, postgresSsl string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", postgresHost, postgresUser, postgresPass, postgresDb, postgresPort, postgresSsl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("successfully connected to the database")
	return db
}
