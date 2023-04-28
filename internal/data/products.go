package data

import (
	"shoppingApp/internal/validator"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Brand       string    `json:"brand"`
	Category    []string  `json:"category"`
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Title != "", "title", "must be provided")
	v.Check(len(product.Title) <= 100, "title", "must not be more than 100 characters long")

	v.Check(product.Description != "", "description", "must be provided")
	v.Check(len(product.Description) <= 500, "description", "must not be more than 500 characters long")

	// todo: add number validation
	v.Check(product.Price != 0, "price", "must be provided")
	v.Check(product.Price >= 0, "price", "must be a positive number")

	v.Check(product.Brand != "", "brand", "must be provided")
	v.Check(len(product.Brand) <= 60, "brand", "must not be more than 60 characters long")

	v.Check(product.Category != nil, "category", "must be provided")
	v.Check(len(product.Category) >= 1, "category", "must contain at least 1 genre")
	v.Check(len(product.Category) <= 5, "category", "must not contain more than 5 genres")
	v.Check(validator.Unique(product.Category), "category", "must not contain duplicate values")
}

/*
	// Dummy data
	  "id": 1,
      "title": "iPhone 9",
      "description": "An Apple phone",
      "price": 549,
      "brand": "Apple",
      "category": "smartphones",
*/
