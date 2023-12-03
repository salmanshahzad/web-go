package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"

	"github.com/salmanshahzad/web-go/internal/app"
	"github.com/salmanshahzad/web-go/internal/database"
	"github.com/salmanshahzad/web-go/internal/utils"
)

//go:embed internal/database/migrations
var migrations embed.FS

//go:embed public
var public embed.FS

func main() {
	env := loadEnvVars()
	db := connectToPostgres(env)
	rdb := connectToRedis(env)
	pub := getPublicFs()

	sm := scs.New()
	sm.Lifetime = 7 * 24 * time.Hour
	sm.Store = goredisstore.New(rdb)

	app := app.NewApplication(db, env, pub, rdb, sm)
	setupGracefulShutdown(app)

	log.Printf("Server starting on port %d", env.Port)
	addr := net.JoinHostPort("0.0.0.0", fmt.Sprint(env.Port))
	if err := http.ListenAndServe(addr, app); err != nil {
		log.Fatalln("Error starting server:", err)
	}
}

func loadEnvVars() *utils.Environment {
	env, err := utils.InitEnv()
	if err != nil {
		log.Fatalln("Could not load environment variables:", err)
	}
	log.Println("Loaded environment variables")
	return env
}

func connectToPostgres(env *utils.Environment) *database.Queries {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", env.DbHost, env.DbPort, env.DbUser, env.DbPassword, env.DbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}
	log.Println("Connected to database")

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalln("Error setting goose dialect:", err)
	}
	goose.SetBaseFS(migrations)
	if err := goose.Up(db, "internal/database/migrations"); err != nil && !errors.Is(err, goose.ErrNoMigrationFiles) {
		log.Fatalln("Error performing database migrations:", err)
	}
	log.Println("Completed database migrations")

	return database.New(db)
}

func connectToRedis(env *utils.Environment) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(env.RedisHost, fmt.Sprint(env.RedisPort)),
		Password: env.RedisPassword,
	})
	log.Println("Connected to Redis")

	return rdb
}

func getPublicFs() *fs.FS {
	pub, err := fs.Sub(public, "public")
	if err != nil {
		log.Fatalln("Could not find public directory:", err)
	}
	return &pub
}

func setupGracefulShutdown(app *app.Application) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting down server")
		app.GracefulShutdown()
		os.Exit(0)
	}()
}
