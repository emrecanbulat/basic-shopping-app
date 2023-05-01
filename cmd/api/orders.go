package main

import (
	"fmt"
	"net/http"
	"shoppingApp/internal/model"
	"shoppingApp/internal/validator"
)

// createProductHandler for the "POST /v1/products" endpoint.
func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductID   int64  `json:"product_id"`
		PaymentType string `json:"payment_type"`
		AmountPaid  int    `json:"amount_paid"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product, _ := model.Product{}.Find("id", input.ProductID)
	if product.ID == 0 {
		app.serverErrorResponse(w, r, model.ErrProductNotFound)
		return
	}

	order := &model.Order{
		User:        *app.contextGetUser(r),
		Product:     product,
		PaymentType: input.PaymentType,
		AmountPaid:  input.AmountPaid,
	}

	v := validator.New()

	if model.ValidateOrder(v, order); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	newOrder, createErr := order.Create()
	if createErr != nil {
		app.serverErrorResponse(w, r, createErr)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/orders/%d", newOrder.ID))

	//todo response struct

	err = app.writeJSON(w, http.StatusCreated, envelope{"order": newOrder}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
