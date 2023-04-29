package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"shoppingApp/internal/client"
	"shoppingApp/internal/data"
	"shoppingApp/internal/validator"
	"time"
	"unicode"
)

type Product struct {
	ID          int64          `json:"id" gorm:"primarykey"` // Unique integer ID for the product
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Title       string         `json:"title" gorm:"type:text;size:100;not null"`
	Description string         `json:"description" gorm:"type:text;size:500;not null"`
	Price       int32          `json:"price" gorm:"type:integer;not null"`
	Brand       string         `json:"brand" gorm:"type:text;size:60;not null"`
	Category    pq.StringArray `json:"category" gorm:"type:text[];not null"`
}

func (product Product) Create() (Product, error) {
	err := client.PostgreSqlClient.Create(&product)
	return product, err.Error
}

func (product Product) Find(query ...interface{}) (Product, error) {
	err := client.PostgreSqlClient.First(&product, query...)
	return product, err.Error
}

func (product Product) Get(filter data.Filters, query ...interface{}) ([]Product, data.Metadata) {
	var products []Product
	postClient := client.PostgreSqlClient
	var totalCount int64

	postClient = postClient.Order(filter.SortColumn() + " " + filter.SortDirection())
	postClient = postClient.Limit(filter.Limit())
	postClient = postClient.Offset(filter.Offset())

	postClient.Find(&products, query...)
	postClient.Model(&Product{}).Count(&totalCount)

	metadata := data.CalculateMetadata(totalCount, filter.Page, filter.PageSize)
	return products, metadata
}

func (product Product) Update(column string, value interface{}) error {
	err := client.PostgreSqlClient.Model(&product).Update(column, value)
	return err.Error
}

func (product Product) Updates(data Product) error {
	err := client.PostgreSqlClient.Model(&product).Updates(data)
	return err.Error
}

func (product Product) Delete(column string, value interface{}) error {
	err := client.PostgreSqlClient.Model(&product).Delete(column, value)
	return err.Error
}

func (product Product) Count(column string, value interface{}) int64 {
	var counter int64
	client.PostgreSqlClient.Model(&product).Where(column, value).Count(&counter)
	return counter
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Title != "", "title", "must be provided")
	v.Check(len(product.Title) <= 100, "title", "must not be more than 100 characters long")

	v.Check(product.Description != "", "description", "must be provided")
	v.Check(len(product.Description) <= 500, "description", "must not be more than 500 characters long")

	v.Check(product.Price != 0 && !unicode.IsDigit(product.Price), "price", "must be provided")
	v.Check(product.Price >= 0, "price", "must be a positive number")

	v.Check(product.Brand != "", "brand", "must be provided")
	v.Check(len(product.Brand) <= 60, "brand", "must not be more than 60 characters long")

	v.Check(product.Category != nil, "category", "must be provided")
	v.Check(len(product.Category) >= 1, "category", "must contain at least 1 genre")
	v.Check(len(product.Category) <= 5, "category", "must not contain more than 5 genres")
	v.Check(validator.Unique(product.Category), "category", "must not contain duplicate values")
}
