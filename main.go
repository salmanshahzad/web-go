package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"salmanshahzad.com/web-go/api"
	"salmanshahzad.com/web-go/database"
	"salmanshahzad.com/web-go/models"
)

func main() {
	loadEnvVars()
	database.Db = connectToPostgres()
	database.Rdb = connectToRedis()

	app := bootstrapApp()
	setupGracefulShutdown()

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	if err := app.Listen(addr); err != nil {
		log.Fatalln("Error starting server:", err)
	}
}

func loadEnvVars() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Could not load environment variables from .env:", err)
	}
	log.Println("Loaded environment variables from .env")
}

func connectToPostgres() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}
	log.Println("Connected to database")

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalln("Error performing database migrations:", err)
	}
	log.Println("Completed database migrations")

	return db
}

func connectToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	log.Println("Connected to Redis")

	return rdb
}

func bootstrapApp() *fiber.App {
	apiRouter := fiber.New()
	apiRouter.Mount("/health", api.NewHealthRouter())
	apiRouter.Mount("/session", api.NewSessionRouter())
	apiRouter.Mount("/user", api.NewUserRouter())

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Mount("/api", apiRouter)
	app.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("public"),
	}))

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	var e *fiber.Error
	if errors.As(err, &e) && e.Code == 500 {
		log.Println("Internal Server Error", e.Error())
		return err
	}
	return c.Status(e.Code).SendString(e.Message)
}

func setupGracefulShutdown() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting down server")
		database.Rdb.Close()
		log.Println("Disconnected from Redis")
		os.Exit(0)
	}()
}
