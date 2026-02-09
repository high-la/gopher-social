package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/high-la/gopher-social/internal/store"
)

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	dsn                string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        time.Duration
}

type application struct {
	config config
	store  store.Storage
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Group routes
	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", app.healthCheckHandler)

		// /posts
		r.Route("/posts", func(r chi.Router) {
			// 1. Routes that DO NOT need the context middleware (no ID yet)
			r.Post("/", app.createPostHandler)

			// 2. Routes that DO need the context middleware (must have an ID)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware) // Middleware applied only here

				r.Get("/", app.getPostHandler)       // Maps to GET /posts/{id}
				r.Patch("/", app.updatePostHandler)  // Maps to PATCH /posts/{id}
				r.Delete("/", app.deletePostHandler) // Maps to DELETE /posts/{id}
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
