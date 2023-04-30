package main

import (
	"errors"
	"net/http"
	"shoppingApp/internal/model"
	"shoppingApp/internal/validator"
	"time"
)

// registerUserHandler handles the registration of new users and generates a JWT token
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &model.User{
		FullName:  input.FullName,
		Email:     input.Email,
		Activated: true,
		Password:  []byte(input.Password),
	}

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	userPass := password{}
	err = userPass.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.Password = userPass.hash

	newUser, createErr := user.Create()
	if createErr != nil {
		switch {
		case errors.Is(createErr, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, createErr)
		}
		return
	}

	role := "user"
	if newUser.ID == 1 {
		role = "admin"
	}

	token, err := model.Token{}.Create(newUser, role, app.config.jwt.secret)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	authentication := map[string]string{
		"token":  string(token.Hash),
		"expiry": token.Expiry.Format(time.DateTime),
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": newUser, "access_token": authentication}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
