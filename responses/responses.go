package responses

import (
	"fmt"
	"net/http"
)

func NotFoundError(w http.ResponseWriter, err error) {
	RespondWithError(w, 404, fmt.Sprintf("Not Found Error: %s", err))
}

func AuthenticationError(w http.ResponseWriter, err error) {
	RespondWithError(w, 403, fmt.Sprintf("Authentization Error: %s", err))
}

func UnauthorizedError(w http.ResponseWriter, err error) {
	RespondWithError(w, 401, fmt.Sprintf("Authorization Error: %s", err))
}

func BadRequestError(w http.ResponseWriter, err error) {
	RespondWithError(w, 400, fmt.Sprintf("Bad Request Error: %s", err))
}

func InternalServerError(w http.ResponseWriter, err error) {
	RespondWithError(w, 500, fmt.Sprintf("Internal Server Error: %s", err))
}

func StatusOkWithContent(w http.ResponseWriter, payload interface{}) {
	RespondWithJSON(w, 200, payload)
}

func StatusOkNoContent(w http.ResponseWriter) {
	RespondWithJSON(w, 201, struct{}{})
}