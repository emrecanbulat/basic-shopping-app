package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Products
	router.HandlerFunc(http.MethodPost, "/v1/products", app.requirePermission(app.createProductHandler))
	router.HandlerFunc(http.MethodGet, "/v1/products/:id", app.showProductHandler)
	router.HandlerFunc(http.MethodPut, "/v1/products/:id", app.requirePermission(app.updateProductHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/products/:id", app.requirePermission(app.updateProductHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/products/:id", app.requirePermission(app.deleteProductHandler))
	router.HandlerFunc(http.MethodGet, "/v1/products", app.listProductsHandler)

	// Users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	// Orders
	router.HandlerFunc(http.MethodPost, "/v1/orders", app.createOrderHandler)
	router.HandlerFunc(http.MethodGet, "/v1/orders/:id", app.showOrderHandler)
	router.HandlerFunc(http.MethodGet, "/v1/orders", app.requirePermission(app.listOrderHandler))

	// Tokens (Generate a new token)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	//return app.recoverPanic(router)
	return app.recoverPanic(app.authenticate(router))
}
