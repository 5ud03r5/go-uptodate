package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/5ud03r5/uptodate/auth"
	"github.com/5ud03r5/uptodate/db"
	"github.com/5ud03r5/uptodate/responses"
	"github.com/lestrrat-go/jwx/jwt"
)


func HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		responses.BadRequestError(w, errors.New("invalid Authorization header format"))
		return
    }

	token, err := auth.RefreshTokenAuth.Decode(parts[1])
	
	if err != nil {
		responses.InternalServerError(w, err)
		return
		
	}
	err = jwt.Validate(token)
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	extractedClaims, err := token.AsMap(r.Context())
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	accessToken, err := auth.GenerateJWTAccessToken(extractedClaims["sub"].(string), extractedClaims["type"].(string))
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	responses.StatusOkWithContent(w, accessToken)
}

func HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email string `json:"email"`
		Endpoint string `json:"endpoint"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responses.BadRequestError(w, err)
		return
	}

	user, err := db.RegisterUser(r.Context(), params.Username, params.Email, params.Endpoint)
	if err != nil {
		responses.InternalServerError(w, err)
		return
		
	}
	responses.StatusOkWithContent(w, user)
}

func HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responses.BadRequestError(w, err)
		return
	}
	user, err := db.LoginUser(r.Context(), params.Username, params.Password)
	if err != nil {
		responses.UnauthorizedError(w, err)
		return
	}

	claims := make(map[string]interface{})
	claims["sub"] = user.ID
	claims["type"] = "user"

	tokensPair, err := auth.GenerateJWTTokens(claims)

	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	responses.StatusOkWithContent(w, tokensPair)	
}