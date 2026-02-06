package main

import (
	"log"
	"os"

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

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
