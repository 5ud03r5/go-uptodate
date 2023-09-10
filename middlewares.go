package main

import (
	"fmt"
	"net/http"

	"github.com/5ud03r5/uptodate/handlers"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)


func authenticatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			handlers.RespondWithError(w, 401, fmt.Sprintf("Unauthorized: %s", err))
		}

		if token == nil || jwt.Validate(token) != nil {
			handlers.RespondWithError(w, 401, "Unauthorized: Missing token or missing required claims")
		}
		next.ServeHTTP(w, r)
	})
}