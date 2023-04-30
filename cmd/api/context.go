package main

import (
	"context"
	"net/http"
	"shoppingApp/internal/model"
)

type contextKey string

// userContextKey is used as a key for getting and setting user information in the request context.
const userContextKey = contextKey("user")
const userRoleContextKey = contextKey("userRole")

// contextSetUser returns a new copy of the request with the provided User struct added to the context.
func (app *application) contextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextSetUserRole(r *http.Request, role string) *http.Request {
	ctx := context.WithValue(r.Context(), userRoleContextKey, role)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the User struct from the request context. The only time that
// this helper should be used is when we logically expect there to be a User struct value
// in the context, and if it doesn't exist it will firmly be an 'unexpected' error, upon we panic.
func (app *application) contextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

func (app *application) contextGetUserRole(r *http.Request) string {
	userRole, ok := r.Context().Value(userRoleContextKey).(string)
	if !ok {
		panic("missing user value in request context")
	}
	return userRole
}
