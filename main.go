package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
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
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cfg, err := utils.LoadConfig(ctx)
	if err != nil {
		slog.Error("could not load config", "err", err)
		os.Exit(1)
	}
	slog.Info("loaded config")

	dbConn, err := getDatabaseConnection(ctx, cfg)
	if err != nil {
		slog.Error("could not connect to database", "err", err)
		os.Exit(1)
	}
	slog.Info("connected to database")
	db := database.New(dbConn)

	rdb, err := getRedisConnection(ctx, cfg)
	if err != nil {
		slog.Error("could not connect to Redis", "err", err)
		os.Exit(1)
	}
	slog.Info("connected to Redis")
	sessStore := goredisstore.New(rdb)

	pub, err := fs.Sub(public, "public")
	if err != nil {
		slog.Error("could not get public directory", "err", err)
		os.Exit(1)
	}

	app := app.NewApplication(cfg, db, rdb, pub, sessStore)
	go func() {
		slog.Info(fmt.Sprintf("starting server on port %d", cfg.Port))
		addr := net.JoinHostPort(net.IPv4zero.String(), strconv.Itoa(cfg.Port))
		if err := http.ListenAndServe(addr, app); err != nil {
			slog.Error("error starting server", "err", err)
			os.Exit(1)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	slog.Info("shutting down server")
	if err := dbConn.Close(ctx); err != nil {
		slog.Warn("error closing database connection", "err", err)
	}
	slog.Info("disconnected from database")
	if err := rdb.Close(); err != nil {
		slog.Warn("error closing Redis connection", "err", err)
	}
	slog.Info("disconnected from Redis")
	slog.Info("shutdown complete")
}

func getDatabaseConnection(ctx context.Context, cfg *utils.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	goose.SetBaseFS(migrations)

	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()
	if err := goose.Up(db, "internal/database/migrations"); err != nil && !errors.Is(err, goose.ErrNoMigrationFiles) {
		return nil, err
	}

	return conn, nil
}

func getRedisConnection(ctx context.Context, cfg *utils.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
