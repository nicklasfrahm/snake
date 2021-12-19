package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mykilio/mykilio.go/pkg/gateway"
)

func main() {
	// Configure logger.
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})

	// Migrate database.
	db, err := MigrateDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Create new fiber app.
	app := fiber.New(fiber.Config{
		ErrorHandler:          gateway.MiddlewareError(),
		DisableStartupMessage: true,
		Prefork:               false,
	})

	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-Token",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
		MaxAge:           600,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	app.Use(gateway.MiddlewareRedirectSlashes())
	app.Use(gateway.MiddlewareContentType(fiber.MIMEApplicationJSONCharsetUTF8))

	// Mount API router.
	app.Use(API(db))

	// Configure fallback route.
	app.Use(gateway.MiddlewareNotFound())

	// Check if the port is valid.
	port := os.Getenv("PORT")
	if port == "" {
		log.Warn().Msg("Missing environment variable: PORT")
		port = "8080"
		log.Warn().Msgf("Using default port: %s/tcp", port)
	}

	// Start server and listen for incoming connections.
	log.Info().Msgf("Server online: %s/tcp", port)
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}
