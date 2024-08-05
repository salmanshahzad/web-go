package app

import (
	"io/fs"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"

	"github.com/salmanshahzad/web-go/internal/database"
	"github.com/salmanshahzad/web-go/internal/utils"
)

type Application struct {
	db     database.Querier
	cfg    *utils.Config
	public *fs.FS
	rdb    *redis.Client
	router *chi.Mux
	sm     *scs.SessionManager
}

func NewApplication(db database.Querier, cfg *utils.Config, public *fs.FS, rdb *redis.Client, sm *scs.SessionManager) *Application {
	app := Application{
		db:     db,
		cfg:    cfg,
		public: public,
		rdb:    rdb,
		router: chi.NewRouter(),
		sm:     sm,
	}

	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.GetHead)
	apiRouter.Mount("/health", app.newHealthRouter())
	apiRouter.Mount("/session", app.newSessionRouter())
	apiRouter.Mount("/user", app.newUserRouter())

	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.router.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   cfg.CorsOrigins,
	}).Handler)
	app.router.Use(middleware.GetHead)
	app.router.Use(sm.LoadAndSave)
	app.router.Mount("/api", apiRouter)

	publicFs := http.FileServer(http.FS(*public))
	app.router.Get("/*", publicFs.ServeHTTP)

	return &app
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.ServeHTTP(w, r)
}
