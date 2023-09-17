package main

import (
	"errors"
	"net/http"

	"github.com/5ud03r5/uptodate/responses"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)


func authenticatorUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, tokenClaims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			responses.UnauthorizedError(w, err)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			responses.UnauthorizedError(w, errors.New("missing or invalid token"))
			return
		}

		if tokenClaims["type"] != "user" {
			responses.UnauthorizedError(w, errors.New("only user can access this endpoint"))
			return
		} 
		
		next.ServeHTTP(w, r)
	})
}

func authenticatorSAMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, tokenClaims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			responses.UnauthorizedError(w, err)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			responses.UnauthorizedError(w, errors.New("missing or invalid token"))
			return
		}

		if tokenClaims["type"] != "service" {
			responses.UnauthorizedError(w, errors.New("only service account can access this endpoint"))
			return
		} 
		
		next.ServeHTTP(w, r)
	})
}