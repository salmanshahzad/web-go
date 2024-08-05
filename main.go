package main

import (
	"context"
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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
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
	cfg := loadConfig()
	db := connectToPostgres(cfg)
	rdb := connectToRedis(cfg)
	pub := getPublicFs()

	sm := scs.New()
	sm.Lifetime = 7 * 24 * time.Hour
	sm.Store = goredisstore.New(rdb)

	app := app.NewApplication(db, cfg, pub, rdb, sm)
	setupGracefulShutdown(app)

	log.Printf("Server starting on port %d", cfg.Port)
	addr := net.JoinHostPort("0.0.0.0", fmt.Sprint(cfg.Port))
	if err := http.ListenAndServe(addr, app); err != nil {
		log.Fatalln("Error starting server:", err)
	}
}

func loadConfig() *utils.Config {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalln("Could not load config:", err)
	}
	log.Println("Loaded config")
	return cfg
}

func connectToPostgres(cfg *utils.Config) database.Querier {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}
	log.Println("Connected to database")

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalln("Error setting goose dialect:", err)
	}
	goose.SetBaseFS(migrations)
	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()
	if err := goose.Up(db, "internal/database/migrations"); err != nil && !errors.Is(err, goose.ErrNoMigrationFiles) {
		log.Fatalln("Error performing database migrations:", err)
	}
	log.Println("Completed database migrations")

	return database.New(conn)
}

func connectToRedis(cfg *utils.Config) *redis.Client {
	opts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		log.Fatalln("Error parsing Redis URL:", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("Error connecting to Redis:", err)
	}
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
