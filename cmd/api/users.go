package main

import (
	"errors"
	"net/http"
	"shoppingApp/internal/model"
	"shoppingApp/internal/validator"
)

// registerUserHandler handles the registration of new users
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
		Activated: false,
		Password:  []byte(input.Password),
	}

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	password := password{}
	err = password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.Password = password.hash

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

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": newUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
