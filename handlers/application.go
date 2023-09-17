package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/5ud03r5/uptodate/custom"
	"github.com/5ud03r5/uptodate/db"
	"github.com/5ud03r5/uptodate/responses"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

var AppRegistryCache *custom.Cache

func HandlerRegisterApplication(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responses.BadRequestError(w, err)
		return
	}

	// Application registry existance check
	cachedIsRegistered := AppRegistryCache.Get(params.Name)
	if cachedIsRegistered != nil {
		responses.BadRequestError(w, errors.New("application is already registered"))
		return
	}
	
	if cachedIsRegistered == nil {
		isRegistered := db.GetRegisteredApplication(r.Context(), params.Name)
		if isRegistered {
			responses.BadRequestError(w, errors.New("application is already registered"))
			return
		}
	}

	// Claims unpack from context
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	// Needs to be user type
	if claims["type"] != "user" {
		responses.BadRequestError(w, errors.New("token needs to be user type"))
		return
	}

	username := claims["sub"].(string)
	// User existance check
	_, err = db.GetUserByUsername(r.Context(), username)
	if err != nil {
		responses.NotFoundError(w, err)
		return
	}

	err = db.RegisterApplication(r.Context(), params.Name, username)
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	// Service account for application creation upon registration
	serviceAccount, err := db.CreateServiceAccount(r.Context(), params.Name)
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	responses.StatusOkWithContent(w, responses.ServiceAccount{AccountName: serviceAccount.AccountName, Password: serviceAccount.Password})
}

func HandlerSubscribeToApplication(w http.ResponseWriter, r *http.Request) {

	applicationName := chi.URLParam(r, "applicationName")

	// Application registry existance check
	cachedIsRegistered := AppRegistryCache.Get(applicationName)
	if cachedIsRegistered == nil {
		isRegistered := db.GetRegisteredApplication(r.Context(), applicationName)
		if !isRegistered {
			responses.NotFoundError(w, errors.New("application is not registered"))
			return
		}
		AppRegistryCache.Set(applicationName, true)
	}

	// Claims unpack from context
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	// Needs to be user type
	if claims["type"] != "user" {
		responses.BadRequestError(w, errors.New("token needs to be user type"))
		return
	}

	username := claims["sub"].(string)

	// User existance check
	_, err = db.GetUserByUsername(r.Context(), username)

	if err != nil {
		responses.NotFoundError(w, err)
		return
	}

	err = db.CreateUserApplicationBinding(r.Context(), username, applicationName)
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}
	responses.StatusOkNoContent(w)
}

func HandlerGetApplicationByName(w http.ResponseWriter, r *http.Request) {
	applicationName := chi.URLParam(r, "applicationName")
	applications, err := db.GetApplicationByName(r.Context(), applicationName)
	if err != nil {
		responses.NotFoundError(w, err)
		return
	}
	responses.StatusOkWithContent(w, applications)
}

func HandlerUpsertApplication(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		Version string `json:"version"`
		Source string `json:"source"`
		Vulnerable bool `json:"vulnerable"`
	}

	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responses.BadRequestError(w, err)
		return
	}

	cachedIsRegistered := AppRegistryCache.Get(params.Name)
	if cachedIsRegistered == nil {
		isRegistered := db.GetRegisteredApplication(r.Context(), params.Name)
		if !isRegistered {
			responses.NotFoundError(w, errors.New("application is not registered"))
			return
		}
		AppRegistryCache.Set(params.Name, true)
	}

	// Claims unpack from context
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}

	// Needs to be application type
	if claims["type"] != "service" {
		responses.BadRequestError(w, errors.New("token needs to be service type"))
		return
	}

	application := db.Application{
		Name: params.Name,
		Version: params.Version,
		Source: params.Source,
		Vulnerable: params.Vulnerable,
		CreatedAt: time.Now().UTC(),
	}

	err = db.UpsertApplication(r.Context(), application)
	if err != nil {
		responses.InternalServerError(w, err)
		return
	}
	responses.StatusOkNoContent(w)
}