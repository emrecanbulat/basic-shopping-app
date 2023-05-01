package main

import (
	"fmt"
	"net/http"
	"shoppingApp/internal/data"
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

	orderDetails := map[string]interface{}{
		"id":           newOrder.ID,
		"status":       "processing",
		"payment_type": newOrder.PaymentType,
		"amount_paid":  newOrder.AmountPaid,
		"order_date":   order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	orderProduct := map[string]interface{}{
		"title":       newOrder.Product.Title,
		"brand":       newOrder.Product.Brand,
		"description": newOrder.Product.Description,
		"category":    newOrder.Product.Category,
		"price":       newOrder.Product.Price,
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"order_details": orderDetails, "product": orderProduct}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showProductHandler for the "GET /v1/products/:id" endpoint.
func (app *application) showOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// standard users can only view their own purchases
	query := fmt.Sprintf("orders.id = %d AND orders.user_id = %d", id, app.contextGetUser(r).ID)
	if app.contextGetUser(r).IsAdmin {
		query = fmt.Sprintf("orders.id = %d", id)
	}

	order, getErr := model.Order{}.Find(query)
	if getErr != nil {
		app.notFoundResponse(w, r)
		return
	}

	orderDetails := map[string]interface{}{
		"status":       "processing",
		"payment_type": order.PaymentType,
		"amount_paid":  order.AmountPaid,
		"order_date":   order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	orderProduct := map[string]interface{}{
		"title":       order.Product.Title,
		"brand":       order.Product.Brand,
		"description": order.Product.Description,
		"category":    order.Product.Category,
		"price":       order.Product.Price,
	}

	orderUser := map[string]interface{}{
		"id":        order.User.ID,
		"full_name": order.User.FullName,
		"email":     order.User.Email,
		"phone":     order.User.Phone,
		"address":   order.User.Address,
	}

	if app.contextGetUser(r).IsAdmin {
		err = app.writeJSON(w, http.StatusOK, envelope{"order_details": orderDetails, "product": orderProduct, "user": orderUser}, nil)
	} else {
		err = app.writeJSON(w, http.StatusOK, envelope{"order_details": orderDetails, "product": orderProduct}, nil)
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listProductsHandler for the "GET /v1/products" endpoint.
func (app *application) listOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID    int
		OrderDate string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	if app.contextGetUser(r).IsAdmin {
		input.UserID = app.readInt(qs, "user_id", 1, v)
	} else {
		input.UserID = int(app.contextGetUser(r).ID)
	}

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")

	input.SortSafelist = []string{"id", "created_at", "-id", "-created_at"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	query := fmt.Sprintf("orders.id > 0")
	orders, meta := model.Order{}.Get(input.Filters, query)

	err := app.writeJSON(w, http.StatusOK, envelope{"meta": meta, "orders": orders}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
