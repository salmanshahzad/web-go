package utils

import (
	"net/http"

	"github.com/cohesivestack/valgo"
	"github.com/go-chi/render"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, msg string) {
	sendMessage(w, r, http.StatusUnauthorized, msg)
}

func UnprocessableEntity(w http.ResponseWriter, r *http.Request, msg string) {
	sendMessage(w, r, http.StatusUnprocessableEntity, msg)
}

func ValidationError(w http.ResponseWriter, r *http.Request, val *valgo.Validation) {
	errors := map[string][]string{}
	for k, v := range val.Errors() {
		errors[k] = v.Messages()
	}
	render.JSON(w, r, map[string]map[string][]string{
		"errors": errors,
	})
	w.WriteHeader(http.StatusUnprocessableEntity)
}

func sendMessage(w http.ResponseWriter, r *http.Request, code int, msg string) {
	render.JSON(w, r, map[string]string{
		"message": msg,
	})
	w.WriteHeader(code)
}
