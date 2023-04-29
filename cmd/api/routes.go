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
	router.HandlerFunc(http.MethodPost, "/v1/products", app.createProductHandler)
	router.HandlerFunc(http.MethodGet, "/v1/products/:id", app.showProductHandler)
	router.HandlerFunc(http.MethodPut, "/v1/products/:id", app.updateProductHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateProductHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/products/:id", app.deleteProductHandler)
	router.HandlerFunc(http.MethodGet, "/v1/products", app.listProductsHandler)

	// Users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(router)
}
