package main

import (
	"fmt"
	"net/http"
	"shoppingApp/internal/data"
	"shoppingApp/internal/validator"
	"time"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       float32  `json:"price"`
		Brand       string   `json:"brand"`
		Category    []string `json:"category"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Brand:       input.Brand,
		Category:    input.Category,
	}

	v := validator.New()

	if data.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

// showProductHandler for the "GET /v1/products/:id" endpoint.
func (app *application) showProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	product := data.Product{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "iPhone 14",
		Description: "An Apple phone",
		Price:       16.999,
		Brand:       "Apple",
		Category:    []string{"Smartphone"},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
