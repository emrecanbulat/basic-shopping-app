package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"shoppingApp/internal/model"
	"strings"
	"time"
)

var jwtParse struct {
	Id    int64
	Email string
	Role  string
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		if strings.HasSuffix(r.URL.Path, "/users") && r.Method == http.MethodPost {
			r = app.contextSetUserRole(r, jwtParse.Role)
			next.ServeHTTP(w, r)
			return
		} else if strings.HasSuffix(r.URL.Path, "/authentication") {
			r = app.contextSetUserRole(r, "user")
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		_ = json.Unmarshal(claims.Raw, &jwtParse)

		_, err = model.Token{}.Find("user_id = ? AND hash = ? AND expiry > ?", jwtParse.Id, token, time.Now())

		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		//user, _ := model.User{}.Find("id = ? AND email = ?", jwtParse.Id, jwtParse.Email)

		// Add the user record to the request context and continue as normal.
		//r = app.contextSetUser(r, &user)
		r = app.contextSetUserRole(r, jwtParse.Role)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
