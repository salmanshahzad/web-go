package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alexedwards/scs/goredisstore"
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
	ctx := context.Background()
	cfg := loadConfig(ctx)
	dbConn := connectToPostgres(ctx, cfg)
	db := database.New(dbConn)
	rdb := connectToRedis(ctx, cfg)
	pub := getPublicFs()
	sessStore := goredisstore.New(rdb)

	app := app.NewApplication(cfg, db, rdb, pub, sessStore)
	go func() {
		log.Printf("Server starting on port %d", cfg.Port)
		addr := net.JoinHostPort(net.IPv4zero.String(), strconv.Itoa(cfg.Port))
		if err := http.ListenAndServe(addr, app); err != nil {
			log.Fatalln("Error starting server:", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down server")
	if err := dbConn.Close(ctx); err != nil {
		log.Println("Error closing database connection:", err)
	}
	if err := rdb.Close(); err != nil {
		log.Println("Error closing Redis connection:", err)
	}
	log.Println("Shutdown complete")
	os.Exit(0)
}

func loadConfig(ctx context.Context) *utils.Config {
	cfg, err := utils.LoadConfig(ctx)
	if err != nil {
		log.Fatalln("Could not load config:", err)
	}
	log.Println("Loaded config")
	return cfg
}

func connectToPostgres(ctx context.Context, cfg *utils.Config) *pgx.Conn {
	conn, err := pgx.Connect(ctx, cfg.DatabaseUrl)
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

	return conn
}

func connectToRedis(ctx context.Context, cfg *utils.Config) *redis.Client {
	opts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		log.Fatalln("Error parsing Redis URL:", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalln("Error connecting to Redis:", err)
	}
	log.Println("Connected to Redis")

	return rdb
}

func getPublicFs() fs.FS {
	pub, err := fs.Sub(public, "public")
	if err != nil {
		log.Fatalln("Could not find public directory:", err)
	}
	return pub
}
