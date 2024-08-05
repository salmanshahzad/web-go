package app

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/cohesivestack/valgo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/salmanshahzad/web-go/internal/utils"
)

func (app *Application) newSessionRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", app.handleSignIn)
	r.Delete("/", app.handleSignOut)
	return r
}

func (app *Application) handleSignIn(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Username string
		Password string
	}
	payload := new(Payload)
	if err := render.DecodeJSON(r.Body, payload); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	val := valgo.Is(
		valgo.String(payload.Username, "username").Not().Blank(),
		valgo.String(payload.Password, "password").Not().Empty(),
	)
	if !val.Valid() {
		utils.ValidationError(w, r, val)
		return
	}

	user, err := app.db.GetUserByUsername(r.Context(), payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Unauthorized(w, r, "Incorrect username or password")
			return
		}
		utils.InternalServerError(w, r, err)
		return
	}

	validPassword, err := utils.VerifyPassword(payload.Password, user.Password)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	if !validPassword {
		utils.Unauthorized(w, r, "Incorrect username or password")
		return
	}

	app.sm.Put(r.Context(), "userId", user.ID)
	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) handleSignOut(w http.ResponseWriter, r *http.Request) {
	if err := app.sm.Destroy(r.Context()); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
