package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) newHealthRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", app.handleHealthCheck)
	return r
}

func (app *Application) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
