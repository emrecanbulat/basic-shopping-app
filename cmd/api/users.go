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
		Address  string `json:"address"`
		Phone    string `json:"phone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &model.User{
		FullName: input.FullName,
		Email:    input.Email,
		IsAdmin:  false,
		Password: []byte(input.Password),
		Address:  input.Address,
		Phone:    input.Phone,
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

	activityErr := SetUserActivity(newUser, r.Method, r.URL.Path, "New User")
	if activityErr != nil {
		app.serverErrorResponse(w, r, activityErr)
		return
	}

	token, err := model.Token{}.Create(newUser, app.config.jwt.secret)
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
