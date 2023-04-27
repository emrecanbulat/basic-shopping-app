package data

import "time"

type Product struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Brand       string    `json:"brand"`
	Category    []string  `json:"category"`
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
