package main

import (
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	mode   = os.Getenv("MODE")
	logger zerolog.Logger
)

func main() {
	// entrypoint
	listenAddress := ""
	isPrettyLog := false
	switch mode {
	case "production":
		listenAddress = ":3000"
	case "development":
		listenAddress = "localhost:3000"
		isPrettyLog = true
	default:
		log.Fatal().Msg("Listen address is not set")
	}

	// app
	app := fiber.New()
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	if isPrettyLog {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	// 10 requests per 1 minute max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        10,
	}))

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to qa-api api")
	})

	// --- main --- //
	app.Get("/submit", func(c *fiber.Ctx) error {
		return c.SendString("Foo")
	})

	// error handling
	if err := app.Listen(listenAddress); err != nil {
		logger.Fatal().Err(err).Msg("Fiber app error")
	}
}
