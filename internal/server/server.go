package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/logging"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	Config        *config.Config
	LoggerService *logging.LoggerService
	DB            *mongo.Client
	httpServer    *http.Server
}

func New(cfg *config.Config, logger *logging.LoggerService, dbClient *mongo.Client) (*Server, error) {
	return &Server{
		Config:        cfg,
		LoggerService: logger,
		DB:            dbClient,
	}, nil
}

func (s *Server) SetupHttpServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeOut) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeOut) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeOut) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("http server is not initialized")
	}

	s.LoggerService.Logger.Info().Str("port", s.Config.Server.Port).Str("env", s.Config.Primary.Env).Msg("Server started")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown the server %w", err)
	}

	// close the database
	if err := s.DB.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to close database connection %w", err) // we are not using the pool connection
		// so this step doesn't make much difference
	}

	return nil
}
