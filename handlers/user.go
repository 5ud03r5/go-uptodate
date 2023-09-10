package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/5ud03r5/uptodate/auth"
	"github.com/5ud03r5/uptodate/db"
)

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
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
	}

	user, err := db.RegisterUser(r.Context(), params.Username, params.Email, params.Endpoint)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Could not register user: %s", err))
		
	}
	RespondWithJSON(w, 201, user)
}

func HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type jwtResponse struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
	}
	user, err := db.LoginUser(r.Context(), params.Username, params.Password)
	if err != nil {
		RespondWithError(w, 403, fmt.Sprintf("Login error: %s", err))
	}

	claims := make(map[string]interface{})
	claims["sub"] = user.ID
	claims["type"] = "user"

	accessToken, refreshToken, err := auth.GenerateJWTTokens(claims)

	if err != nil {
		RespondWithError(w, 500, err.Error())
	}

	tokensPair := jwtResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	RespondWithJSON(w, 200, tokensPair)	
}