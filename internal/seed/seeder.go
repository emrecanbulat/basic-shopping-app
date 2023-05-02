package seed

import (
	"log"
)

func Seed() {
	UserSeed()
	ProductSeed()
	log.Println("Seeding successfully completed")
}
