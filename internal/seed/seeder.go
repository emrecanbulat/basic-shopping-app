package seed

import "fmt"

func Seed() {
	UserSeed()
	ProductSeed()
	fmt.Println("Seeding finished")
}
