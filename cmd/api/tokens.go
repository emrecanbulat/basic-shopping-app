package main

import (
	"errors"
	"net/http"
	"shoppingApp/internal/model"
	"shoppingApp/internal/validator"
	"time"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	model.ValidateEmail(v, input.Email)
	model.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := model.User{}.Find("email", input.Email)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	password := password{hash: user.Password}

	match, err := password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token := model.Token{}

	tokenCheck := model.Token{}.Count("user_id", user.ID)
	if tokenCheck > 0 {
		token, err = model.Token{}.Find("user_id", user.ID)
		newToken, _ := model.GenerateToken(user, app.config.jwt.secret)
		err = token.Update("hash", newToken.Hash)
		token.Hash = newToken.Hash
	} else {
		token, err = model.Token{}.Create(user, app.config.jwt.secret)
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	authentication := map[string]string{
		"token":  string(token.Hash),
		"expiry": token.Expiry.Format(time.DateTime),
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication": authentication}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
