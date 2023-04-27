package main

import (
	"fmt"
	"net/http"
	"shoppingApp/internal/data"
	"time"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new product")
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
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
