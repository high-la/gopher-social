package main

import (
	"os"
	"strconv"
	"time"

	"github.com/high-la/gopher-social/internal/db"
	"github.com/high-la/gopher-social/internal/mailer"
	"github.com/high-la/gopher-social/internal/store"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

const version = "0.0.1"

// ..............................swagger directives
//	@title			GopherSocial
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Description for the tool
func main() {

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// env files
	// u can load multiple files
	err := godotenv.Load(".env", ".env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	addr := os.Getenv("GOPHER_SOCIAL_ADDR")
	apiURL := os.Getenv("GOPHER_SOCIAL_EXTERNAL_URL")
	frontendURL := os.Getenv("GOPHER_SOCIAL_FRONTEND_URL")
	dsn := os.Getenv("GOPHER_SOCIAL_DSN")
	maxOpenConnections := getEnvInt("GOPHER_SOCIAL_DB_MAX_OPEN_CONNECTIONS", 30)
	maxIdleConnections := getEnvInt("GOPHER_SOCIAL_DB_MAX_IDLE_CONNECTIONS", 30)
	maxIdleTime := getEnvTime("GOPHER_SOCIAL_DB_MAX_IDLE_TIME", 15*time.Minute)

	cfg := config{
		addr:        addr,
		apiURL:      apiURL,
		frontendURL: frontendURL,
		db: dbConfig{
			dsn:                dsn,
			maxOpenConnections: maxOpenConnections,
			maxIdleConnections: maxIdleConnections,
			maxIdleTime:        maxIdleTime,
		},
		env: os.Getenv("GOPHER_SOCIAL_APP_ENV"),
		mail: mailConfig{
			expiry:    time.Hour * 24 * 3,
			fromEmail: os.Getenv("GOPHER_SOCIAL_FROM_EMAIL"),
			sendGrid: sendGridConfig{
				apiKey: os.Getenv("GOPHER_SOCIAL_SENDGRID_API_KEY"),
			},
		},
	}

	// Database
	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConnections, cfg.db.maxIdleConnections, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal("unable to connect to the database \n", err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	// store
	store := store.NewStorage(db)

	// Mailer
	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}

// getEnvInt reads an env var and converts to int, or returns a default
func getEnvInt(key string, fallback int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return val
}

// getEnvTime reads an env var and parses it as a duration (e.g. "15m", "1h")
func getEnvTime(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}
