package main

import (
	"fmt"
	"os"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/handlers/http"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	postgres "github.com/arasan1289/hexagonal-demo/internal/adapters/storage/db"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/storage/db/repository"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/service"
)

func main() {
	// Initialize config
	config, err := config.New()
	if err != nil {
		fmt.Println("Error initializing config:", err)
		os.Exit(1)
	}
	// Initialize logger
	log, err := logger.Set(config)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		os.Exit(1)
	}
	// Initialize Gorm custom logger
	customGormLogger := logger.NewGormLogger()
	// Initialize DB
	conn, err := postgres.New(config.DB, customGormLogger)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing Postgres DB")
		os.Exit(1)
	}
	defer conn.Close()
	log.Info().Msg("Successfully connected to DB")

	// Migrate DB
	conn.Migrate(&domain.User{})

	log.Info().Msg("Successfully migrated user table")

	// Initialize Handlers
	userRepo := repository.NewUserRepository(conn)
	userSvc := service.NewUserService(userRepo, log)
	UserHandler := http.NewUserHandler(userSvc, config.App, log)

	otpSvc := service.NewOtpService(log, config.App)
	OtpHandler := http.NewOtpHandler(otpSvc, userSvc, log, config.App)

	// Initialize router
	router, err := http.NewRouter(config, log, *UserHandler, *OtpHandler)
	if err != nil {
		log.Error().Err(err).Msg("Error Initializing router")
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	log.Info().Msgf("HTTP server running on %s", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		log.Error().Err(err).Msg("Error starting the HTTP server")
		os.Exit(1)
	}

}
