package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/5ud03r5/uptodate/auth"
	"github.com/lestrrat-go/jwx/jwt"
)


func HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
        RespondWithError(w, 400, "Invalid Authorization header format")
    }

	token, err := auth.RefreshTokenAuth.Decode(parts[1])
	
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error during decoding: %s", err))
	}

	errValidate := jwt.Validate(token)
	if errValidate != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error during validation: %s", errValidate))
	}

	extractedClaims, err := token.AsMap(r.Context())
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error during conversion to map: %s", err))
	}

	accessToken, err := auth.GenerateJWTAccessToken(extractedClaims["sub"].(string), extractedClaims["type"].(string))
	if err != nil {
		RespondWithError(w, 500, err.Error())
	}

	RespondWithJSON(w, 200, accessToken)
}