package handler

import (
	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/service"
	"github.com/rs/zerolog"
)

type Handler struct {
	Auth *AuthHandler
	User *UserHandler
	Url  *UrlHandler
}

func NewHandler(cfg *config.Config, services *service.Service, logger zerolog.Logger) *Handler {
	return &Handler{
		Auth: NewAuthHandler(cfg, services, logger),
		User: NewUserHandler(logger, services),
		Url:  NewUrlHandler(logger, services),
	}
}
