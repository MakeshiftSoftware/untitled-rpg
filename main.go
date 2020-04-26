package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"untitled_rpg/logger"
	"untitled_rpg/migrate"
	"untitled_rpg/server"
	"untitled_rpg/service"
	"untitled_rpg/store"
	"untitled_rpg/token"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

func main() {
	// Seed random
	rand.Seed(time.Now().UnixNano())

	config := loadConfig()

	logger := logger.NewZerologLogger(config.Debug)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info().Msg("Server starting")

	db := dbConnect(logger, config)

	migrate.Migrate(logger, db)

	tokenProvider := token.NewProvider(config.Key)
	accountStore := store.NewAccountStore(db)
	accountService := service.NewAccountService(accountStore)
	authService := service.NewAuthService(accountStore, tokenProvider)

	server := server.NewServer(logger, config.Port, accountService, authService)

	go server.Start()

	logger.Info().Msg("Server started")

	<-c

	logger.Info().Msg("Shutdown started")

	// Graceful shutdown here; stop services

	logger.Info().Msg("Shutdown complete")

	os.Exit(0)
}

// dbConnect connects to postgres and pings the database to test the connection.
func dbConnect(logger logger.Logger, config config) *sqlx.DB {
	db, err := sqlx.Open("pgx", config.Database)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open database")
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to ping database")
	}

	logger.Info().Msg("Connected to database")

	// TODO: these parameters should be set as part of the application's config
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(12)

	var dbVersion string
	if err := db.QueryRow("SELECT version()").Scan(&dbVersion); err != nil {
		logger.Fatal().Err(err).Msg("Failed to query database version")
	}
	logger.Info().Str("databaseVersion", dbVersion).Send()
	return db
}

// config contains the server configuration.
type config struct {
	Debug    bool   `default:"false"` // Debug indicates whether debugging is enabled.
	Database string `required:"true"` // Database is the database connection url.
	Port     int    `required:"true"` // Port is the port that the server listens on.
	Key      string `required:"true"` // Key is the secret key used when generating auth tokens.
}

// loadConfig loads the server configuration from environment.
func loadConfig() config {
	godotenv.Load()

	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse environment variables")
	}

	return config
}
