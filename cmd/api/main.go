package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/high-la/gopher-social/internal/store"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// u can load multiple files
	err := godotenv.Load(".env", ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("GOPHER_SOCIAL_ADDR")
	dsn := os.Getenv("GOPHER_SOCIAL_DSN")
	maxOpenConnections := getEnvInt("GOPHER_SOCIAL_DB_MAX_OPEN_CONNECTIONS", 30)
	maxIdleConnections := getEnvInt("GOPHER_SOCIAL_DB_MAX_IDLE_CONNECTIONS", 30)
	maxIdleTime := getEnvTime("GOPHER_SOCIAL_DB_MAX_IDLE_TIME", 15*time.Minute)

	cfg := config{
		addr: addr,
		db: dbConfig{
			dsn:                dsn,
			maxOpenConnections: maxOpenConnections,
			maxIdleConnections: maxIdleConnections,
			maxIdleTime:        maxIdleTime,
		},
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
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
