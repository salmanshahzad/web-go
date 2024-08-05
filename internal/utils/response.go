package utils

import (
	"net/http"

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

func sendMessage(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.WriteHeader(code)
	render.JSON(w, r, map[string]string{
		"message": msg,
	})
}
