package main

import (
	"time"

	"github.com/high-la/gopher-social/internal/auth"
	"github.com/high-la/gopher-social/internal/db"
	"github.com/high-la/gopher-social/internal/env"
	"github.com/high-la/gopher-social/internal/mailer"
	"github.com/high-la/gopher-social/internal/ratelimiter"
	"github.com/high-la/gopher-social/internal/store"
	"github.com/high-la/gopher-social/internal/store/cache"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
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

	// Get env vars
	appEnv := env.GetString("GOPHER_SOCIAL_APP_ENV", "development")
	//
	addr := env.GetString("GOPHER_SOCIAL_ADDR", "")
	apiURL := env.GetString("GOPHER_SOCIAL_EXTERNAL_URL", "")
	frontendURL := env.GetString("GOPHER_SOCIAL_FRONTEND_URL", "")
	// Database
	dsn := env.GetString("GOPHER_SOCIAL_DSN", "")
	maxOpenConnections := env.GetInt("GOPHER_SOCIAL_DB_MAX_OPEN_CONNECTIONS", 30)
	maxIdleConnections := env.GetInt("GOPHER_SOCIAL_DB_MAX_IDLE_CONNECTIONS", 30)
	maxIdleTime := env.GetTime("GOPHER_SOCIAL_DB_MAX_IDLE_TIME", 15*time.Minute)
	// Redis
	redisAddr := env.GetString("GOPHER_SOCIAL_REDIS_ADDRESS", "")
	redisPassword := env.GetString("GOPHER_SOCIAL_REDIS_PASSWORD", "")
	redisDB := env.GetInt("GOPHER_SOCIAL_REDIS_DB", 0)
	redisEnabled := env.GetBool("GOPHER_SOCIAL_REDIS_ENABLED", false)
	// Email
	fromEmail := env.GetString("GOPHER_SOCIAL_FROM_EMAIL", "")
	sendGridApiKey := env.GetString("GOPHER_SOCIAL_SENDGRID_API_KEY", "")
	// Basic Auth
	basicAuthUsername := env.GetString("GOPHER_SOCIAL_BASIC_AUTH_USERNAME", "")
	basicAuthPassword := env.GetString("GOPHER_SOCIAL_BASIC_AUTH_PASSWORD", "")
	// Auth token
	authTokenSecret := env.GetString("GOPHER_SOCIAL_AUTH_TOKEN_SECRET", "")
	// Rate limiter
	reqPerTimeFrame := env.GetInt("GOPHER_SOCIAL_RATELIMITER_REQUEST_COUNT", 20)
	isRateLimiterEnabled := env.GetBool("GOPHER_SOCIAL_RATELIMITER_ENABLED", true)

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
		redisCfg: redisConfig{
			addr:     redisAddr,
			password: redisPassword,
			db:       redisDB,
			enabled:  redisEnabled,
		},
		env: appEnv,
		mail: mailConfig{
			expiry:    time.Hour * 24 * 3,
			fromEmail: fromEmail,
			sendGrid: sendGridConfig{
				apiKey: sendGridApiKey,
			},
		},
		auth: authConfig{
			basic: basicConfig{
				username: basicAuthUsername,
				password: basicAuthPassword,
			},
			token: tokenConfig{
				secret: authTokenSecret,
				expiry: time.Hour * 24 * 3,
				issuer: "gophersocial",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: reqPerTimeFrame,
			TimeFrame:            time.Second * 5,
			Enabled:              isRateLimiterEnabled,
		},
	}

	// Database
	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConnections, cfg.db.maxIdleConnections, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal("unable to connect to the database \n", err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	// Cache (Redis)
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.password, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
	}
	// store
	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	// Mailer
	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	// Auth
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.issuer, cfg.auth.token.issuer)

	// Rate limiter
	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   ratelimiter,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
