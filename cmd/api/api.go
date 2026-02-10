package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/high-la/gopher-social/docs" // This is required to generate swagger docs
	"github.com/high-la/gopher-social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
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
	logger *zap.SugaredLogger
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

		// http://localhost:40100/swagger/doc.json
		docsURL := fmt.Sprintf("%s/v1/swagger/doc.json", app.config.apiURL)

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL), // The url pointing to API definition
		))

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

		// /users
		r.Route("/users", func(r chi.Router) {

			//
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware) // Middleware applied only here

				r.Get("/", app.getUserHandler) // Maps to GET /users/{id}

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	// Docs
	// Parse the external URL to strip the protocol (http://)
	parsedURL, err := url.Parse(app.config.apiURL)
	if err != nil {
		app.logger.Fatalf("invalid API URL: %v", err)
	}

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = parsedURL.Host                // Dynamically gets "localhost:40100"
	docs.SwaggerInfo.Schemes = []string{parsedURL.Scheme} // Dynamically gets "http"
	docs.SwaggerInfo.BasePath = "/v1"

	// .
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
