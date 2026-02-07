package main

import (
	"log"
	"os"

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

	cfg := config{
		addr: addr,
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
