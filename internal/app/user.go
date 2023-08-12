package app

import (
	"log"
	"net/http"
	"strings"

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

	payload.Username = strings.TrimSpace(payload.Username)
	payload.Password = strings.TrimSpace(payload.Password)
	if len(payload.Username) == 0 || len(payload.Password) == 0 {
		utils.UnprocessableEntity(w, r, "Username and password are required")
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
	log.Println("Created user", payload.Username)
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

	payload.Username = strings.TrimSpace(payload.Username)
	if len(payload.Username) == 0 {
		utils.UnprocessableEntity(w, r, "Username is required")
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

	payload.Password = strings.TrimSpace(payload.Password)
	if len(payload.Password) == 0 {
		utils.UnprocessableEntity(w, r, "Password is required")
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

	log.Println("Deleted user", user.Username)
	w.WriteHeader(http.StatusNoContent)
}
