package seed

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
	"shoppingApp/internal/model"
)

func ProductSeed() {
	counter := model.Product{}.Count("", "")
	if counter < 1 {
		for i := 0; i < 10; i++ {
			product := &model.Product{
				Title:       gofakeit.Lunch(),
				Price:       int32(gofakeit.Price(100, 10000)),
				Description: gofakeit.Sentence(4),
				Brand:       gofakeit.Company(),
				Category: pq.StringArray{
					"Food",
				},
			}
			_, err := product.Create()
			if err != nil {
				return
			}
		}
	}
}
