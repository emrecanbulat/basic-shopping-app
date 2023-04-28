package main

import (
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"shoppingApp/internal/data"
	"shoppingApp/internal/model"
	"shoppingApp/internal/validator"
)

// createProductHandler for the "POST /v1/products" endpoint.
func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       int32    `json:"price"`
		Brand       string   `json:"brand"`
		Category    []string `json:"category"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &model.Product{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Brand:       input.Brand,
		Category:    input.Category,
	}

	v := validator.New()

	if model.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	newProduct, createErr := product.Create()
	if createErr != nil {
		app.serverErrorResponse(w, r, createErr)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/products/%d", newProduct.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"product": newProduct}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showProductHandler for the "GET /v1/products/:id" endpoint.
func (app *application) showProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	product, getErr := model.Product{}.Find("id", id)
	if getErr != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateProductHandler for the "PUT /v1/products/:id" endpoint.
func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	product, getErr := model.Product{}.Find("id", id)
	if getErr != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		Price       *int32   `json:"price"`
		Brand       *string  `json:"brand"`
		Category    []string `json:"category"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		product.Title = *input.Title
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if product.Price != 0 && input.Price != nil {
		product.Price = *input.Price
	}

	if input.Brand != nil {
		product.Brand = *input.Brand
	}

	if input.Category != nil {
		product.Category = input.Category
	}

	v := validator.New()

	if model.ValidateProduct(v, &product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = product.Updates(product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteProductHandler for the "DELETE /v1/products/:id" endpoint.
func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	check := model.Product{}.Count("id", id)

	if check == 0 {
		app.notFoundResponse(w, r)
		return
	}

	deleteErr := model.Product{}.Delete("id", id)
	if deleteErr != nil {
		app.serverErrorResponse(w, r, deleteErr)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "product successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listProductsHandler for the "GET /v1/products" endpoint.
func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string
		Brand    string
		Category []string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Brand = app.readString(qs, "brand", "")
	input.Category = app.readCSV(qs, "category", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "title", "price", "brand"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// todo add sort to query
	products := model.Product{}.Get(input.Limit(), input.Offset(),
		"(to_tsvector('simple', brand) @@ plainto_tsquery('simple', ?) OR ? = '') AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', ?) OR ? = '') AND (category @> ? OR ? = '{}')",
		input.Brand, input.Brand, input.Title, input.Title, pq.Array(input.Category), pq.Array(input.Category),
	)

	err := app.writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
