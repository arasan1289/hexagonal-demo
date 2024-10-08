//	@title			Hexagonal API
//	@version		1.0
//	@description	This is a swagger docs for Hexagonal API.

//	@host		localhost:3000
//	@BasePath	/api/v1

//	@SecurityDefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization
//	@schemes					http https

package main

import (
	"fmt"
	"os"

	_ "github.com/arasan1289/hexagonal-demo/docs"
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
	conn, err := postgres.New(config, customGormLogger)
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

	authSvc := service.NewAuthService(log, config.App)

	otpSvc := service.NewOtpService(log, config.App)
	OtpHandler := http.NewOtpHandler(otpSvc, userSvc, authSvc, log, config.App)

	authhandler := http.NewAuthHandler(authSvc, userSvc, log)

	// Initialize router
	router, err := http.NewRouter(config, log, *UserHandler, *OtpHandler, authSvc, *authhandler)
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
