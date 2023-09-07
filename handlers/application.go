package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/5ud03r5/uptodate/db"
	"github.com/go-chi/chi"
)

func HandlerGetApplicationByName(w http.ResponseWriter, r *http.Request) {
	applicationName := chi.URLParam(r, "applicationName")
	applications, err := db.GetApplicationByName(applicationName)
	if err != nil {
		fmt.Printf("Error getting the applications: %s", err)
		return
	}
	respondWithJSON(w, 200, applications)
}

func HandlerUsertApplication(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	application := db.Application{
		Name: params.Name,
		Version: params.Version,
		Source: params.Source,
		Vulnerable: params.Vulnerable,
		CreatedAt: time.Now().UTC(),
	}

	err = db.UpsertApplication(application)
	if err != nil {
		fmt.Printf("Error adding application: %s", err)
		return
	}
}