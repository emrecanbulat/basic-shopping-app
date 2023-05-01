package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"shoppingApp/internal/client"
	"shoppingApp/internal/validator"
	"time"
)

var ErrProductNotFound = errors.New("the product you wanted to buy was not found")

type Order struct {
	ID          int64          `json:"id" gorm:"primarykey"` // Unique integer ID for the product
	UserID      int64          `json:"-"`
	User        User           `json:"user" gorm:"references:ID"`
	ProductID   int64          `json:"-"`
	Product     Product        `json:"product" gorm:"references:ID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Status      int            `json:"status" gorm:"default:0"`
	PaymentType string         `json:"payment_type" gorm:"type:text;size:100;not null"`
	AmountPaid  int            `json:"amount_paid"  gorm:"type:int"`
}

func (order Order) Create() (Order, error) {
	result := client.PostgreSqlClient.Create(&order)
	return order, result.Error
}

func validateOrderProduct(v *validator.Validator, productID int64) {
	product, _ := Product{}.Find("id", productID)
	v.Check(product.ID != 0, "product_id", ErrProductNotFound.Error())
}

func validateOrderPayment(v *validator.Validator, productID int64, price int) {
	product, _ := Product{}.Find("id", productID)
	msg := fmt.Sprintf("The amount you want to pay must be equal to the product price (the product price is %d)", product.Price)
	v.Check(int(product.Price) == price, "price", msg)
}

func validateOrderPaymentType(v *validator.Validator, paymentType string) {
	allowedPaymentTypes := []string{"Cash", "Credit_card"}
	check := inArray(allowedPaymentTypes, paymentType)
	v.Check(check != -1, "payment_type", "the payment type you provided is not allowed (only 'Cash' and 'Credit_card' are acceptable)")
}

func ValidateOrder(v *validator.Validator, order *Order) {
	v.Check(order.Product.ID != 0, "product_id", "must be provided")
	validateOrderProduct(v, order.Product.ID)

	v.Check(order.PaymentType != "", "payment_type", "must be provided")
	validateOrderPaymentType(v, order.PaymentType)

	v.Check(order.AmountPaid != 0, "price", "must be provided")
	validateOrderPayment(v, order.Product.ID, order.AmountPaid)
}

func inArray[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}
