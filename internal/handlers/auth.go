package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	auth "github.com/5ud03r5/uptodate/internal/auth"
	db "github.com/5ud03r5/uptodate/internal/db"
	responses "github.com/5ud03r5/uptodate/internal/responses"
	"github.com/lestrrat-go/jwx/jwt"
)
func HandlerGetServiceAccountAccessToken(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AccountName string `json:"account_name"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responses.BadRequestError(w, err)
		return
	}
	serviceAccount, err := db.LoginServiceAccount(r.Context(), params.AccountName, params.Password)
	if err != nil {
		responses.UnauthorizedError(w, err)
		return
	}

	sub := serviceAccount.ID
	accessType := "service"

	accessToken, err := auth.GenerateJWTAccessToken(sub, accessType)

	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	responses.StatusOkWithContent(w, accessToken)
}


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
	responses.StatusOkWithContent(w, responses.User{Username: user.Username, Password: user.Password, Email: user.Email, Endpoint: user.Endpoint})
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