package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/database"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/logging"
	"github.com/harry713j/minurly/internal/middleware"
	"github.com/harry713j/minurly/internal/repository"
	"github.com/harry713j/minurly/internal/routes"
	"github.com/harry713j/minurly/internal/server"
	"github.com/harry713j/minurly/internal/service"
	"github.com/harry713j/minurly/internal/validation"
)

const DefaultContextTimeout = 30

func main() {
	validation.InitValidator()
	cfg := config.LoadConfig()
	logger := logging.NewLogger(cfg)

	client, db := database.New(cfg, logger)
	authRepo := repository.NewAuthRepo(db, logger.Logger)
	userRepo := repository.NewUserRepo(db, logger.Logger)
	urlRepo := repository.NewUrlRepo(db, logger.Logger)

	repository := repository.NewRepository(authRepo, userRepo, urlRepo)

	srv, err := server.New(cfg, logger, client)

	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to initialize server")
	}

	authService := service.NewAuthService(repository, logger.Logger)
	userService := service.NewUserService(cfg, repository, logger.Logger)
	urlService := service.NewUrlService(repository, logger.Logger)

	services := service.NewService(authService, userService, urlService)
	// handler
	handlers := handler.NewHandler(cfg, services, logger.Logger)
	middleware := middleware.NewMiddleware(srv)

	r := routes.NewRouter(cfg, handlers, middleware)
	srv.SetupHttpServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Fatal().Err(err).Msg("failed to start the server")
		}
	}()
	// wait for interrupt signal
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(DefaultContextTimeout*time.Second))

	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	stop()
	cancel()

	logger.Logger.Info().Msg("server stopped gracefully")
}
