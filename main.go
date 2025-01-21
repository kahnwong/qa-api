package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kahnwong/qa-api/controller"
	_ "github.com/kahnwong/qa-api/docs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	mode   = os.Getenv("MODE")
	logger zerolog.Logger

	apiKey        = os.Getenv("QA_API_KEY")
	protectedURLs = []*regexp.Regexp{
		regexp.MustCompile("^/submit$"),
	}
)

func validateAPIKey(c *fiber.Ctx, key string) (bool, error) {
	hashedAPIKey := sha256.Sum256([]byte(apiKey))
	hashedKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func authFilter(c *fiber.Ctx) bool {
	originalURL := strings.ToLower(c.OriginalURL())

	for _, pattern := range protectedURLs {
		if pattern.MatchString(originalURL) {
			return false
		}
	}
	return true
}

// @title QA API
// @version 1.0
// @contact.name Karn Wong
// @contact.email karn@karnwong.me
// @license.name MIT
// @host localhost:3000
// @BasePath /
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

	app.Get("/swagger/*", swagger.HandlerDefault)

	// 5 requests per 1 minute max
	app.Use("/submit", limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        5,
	}))

	// auth
	app.Use(keyauth.New(keyauth.Config{
		Next:      authFilter,
		KeyLookup: "cookie:access_token",
		Validator: validateAPIKey,
	}))

	app.Get("/", controller.RootController)

	// --- main --- //
	app.Post("/submit", controller.SubmitController)

	// error handling
	if err := app.Listen(listenAddress); err != nil {
		logger.Fatal().Msg("Fiber app error")
	}
}
