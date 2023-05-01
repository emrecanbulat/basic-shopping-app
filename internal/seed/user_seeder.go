package seed

import (
	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
	"shoppingApp/internal/model"
)

func UserSeed() {
	counter := model.User{}.Count("", "")
	admin := model.User{}.Count("email", "admin@admin.com")
	if admin < 1 {
		user := &model.User{
			FullName: "Admin",
			Email:    "admin@admin.com",
			Phone:    gofakeit.Phone(),
			IsAdmin:  true,
			Password: GeneratePasswordHash(),
			Address:  gofakeit.Address().Address,
		}
		_, err := user.Create()
		if err != nil {
			return
		}
	}

	if counter < 2 {
		for i := 0; i < 10; i++ {
			user := &model.User{
				FullName: gofakeit.Name(),
				Email:    gofakeit.Email(),
				Phone:    gofakeit.Phone(),
				IsAdmin:  false,
				Password: GeneratePasswordHash(),
				Address:  gofakeit.Address().Address,
			}
			_, err := user.Create()
			if err != nil {
				return
			}
		}
	}
}

func GeneratePasswordHash() []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)
	return hash
}
