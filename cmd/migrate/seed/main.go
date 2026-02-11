package main

import (
	"log"
	"os"
	"time"

	"github.com/high-la/gopher-social/internal/db"
	"github.com/high-la/gopher-social/internal/store"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env", ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("GOPHER_SOCIAL_DSN")

	conn, err := db.New(dsn, 3, 3, time.Minute*15)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
