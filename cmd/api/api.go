package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/high-la/gopher-social/docs" // This is required to generate swagger docs
	"github.com/high-la/gopher-social/internal/auth"
	"github.com/high-la/gopher-social/internal/mailer"
	"github.com/high-la/gopher-social/internal/store"
	"github.com/high-la/gopher-social/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
}

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	expiry time.Duration
	issuer string
}

type basicConfig struct {
	username string
	password string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	expiry    time.Duration
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	dsn                string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        time.Duration
}

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
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
		// r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)
		r.Get("/health", app.healthCheckHandler)

		// http://localhost:40100/swagger/doc.json
		docsURL := fmt.Sprintf("%s/v1/swagger/doc.json", app.config.apiURL)

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL), // The url pointing to API definition
		))

		// /posts
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			// 1. Routes that DO NOT need the context middleware (no ID yet)
			r.Post("/", app.createPostHandler)

			// 2. Routes that DO need the context middleware (must have an ID)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware) // Middleware applied only here

				r.Get("/", app.getPostHandler)                                           // Maps to GET /posts/{id}
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler)) // Maps to PATCH /posts/{id}
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))    // Maps to DELETE /posts/{id}
			})
		})

		// /users
		r.Route("/users", func(r chi.Router) {

			r.Put("/activate/{token}", app.activateUserHandler)

			//
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler) // Maps to GET /users/{id}

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})

		})
		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
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

	// Graceful shutdown
	shutdown := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
