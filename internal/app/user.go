package app

import (
	"net/http"

	"github.com/cohesivestack/valgo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/salmanshahzad/web-go/internal/database"
	"github.com/salmanshahzad/web-go/internal/utils"
)

func (app *Application) newUserRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", app.handleCreateUser)
	r.Group(func(r chi.Router) {
		r.Use(app.verifyAuth)
		r.Get("/", app.handleGetUser)
		r.Put("/username", app.handleEditUsername)
		r.Put("/password", app.handleEditPassword)
		r.Delete("/", app.handleDeleteUser)
	})
	return r
}

func (app *Application) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)
	render.JSON(w, r, map[string]any{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (app *Application) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Username string
		Password string
	}
	payload := new(Payload)
	if err := render.DecodeJSON(r.Body, payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	userCount, err := app.db.CountUsersWithUsername(r.Context(), payload.Username)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	if userCount > 0 {
		utils.UnprocessableEntity(w, r, "Username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	user := database.CreateUserParams{
		Username: payload.Username,
		Password: hashedPassword,
	}
	userId, err := app.db.CreateUser(r.Context(), user)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	app.sm.Put(r.Context(), "userId", userId)
	w.WriteHeader(http.StatusCreated)
}

func (app *Application) handleEditUsername(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Username string
	}
	payload := new(Payload)
	if err := render.DecodeJSON(r.Body, payload); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	val := valgo.Is(
		valgo.String(payload.Username, "username").Not().Blank(),
	)
	if !val.Valid() {
		utils.ValidationError(w, r, val)
		return
	}

	userCount, err := app.db.CountUsersWithUsername(r.Context(), payload.Username)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	if userCount > 0 {
		utils.UnprocessableEntity(w, r, "Username already exists")
		return
	}

	user := r.Context().Value("user").(database.User)
	params := database.UpdateUsernameParams{
		ID:       user.ID,
		Username: payload.Username,
	}
	if err := app.db.UpdateUsername(r.Context(), params); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) handleEditPassword(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Password string
	}
	payload := new(Payload)
	if err := render.DecodeJSON(r.Body, payload); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	val := valgo.Is(
		valgo.String(payload.Password, "password").Not().Empty(),
	)
	if !val.Valid() {
		utils.ValidationError(w, r, val)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	user := r.Context().Value("user").(database.User)
	params := database.UpdatePasswordParams{
		ID:       user.ID,
		Password: hashedPassword,
	}
	if err := app.db.UpdatePassword(r.Context(), params); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := app.sm.Destroy(r.Context()); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	user := r.Context().Value("user").(database.User)
	if err := app.db.DeleteUser(r.Context(), user.ID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
