package main

import (
	"errors"
	"net/http"

	"github.com/5ud03r5/uptodate/custom"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)


func authenticatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			custom.UnauthorizedError(w, err)
		}

		if token == nil || jwt.Validate(token) != nil {
			custom.UnauthorizedError(w, errors.New("missing or invalid token"))
		}
		
		next.ServeHTTP(w, r)
	})
}