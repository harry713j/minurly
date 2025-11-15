package handler

import (
	"net/http"

	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/harry713j/minurly/internal/middleware"
	"github.com/harry713j/minurly/internal/models"
	"github.com/harry713j/minurly/internal/service"
	"github.com/harry713j/minurly/internal/utils"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	services *service.Service
	log      zerolog.Logger
}

func NewUserHandler(logger zerolog.Logger, services *service.Service) *UserHandler {
	return &UserHandler{
		log:      logger,
		services: services,
	}
}

func (u *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserId(r)

	if !ok {
		u.log.Warn().Str("userId", userId).Msg("userId value not found from the request context")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
		return
	}

	resp, err := u.services.User.Fetch(userId)
	if err != nil {
		u.log.Err(err).Str("userId", userId).Msg("failed to get the user data and their short urls")
		utils.RespondError(w, err)
		return
	}

	type responseUser struct {
		Message string               `json:"message"`
		Data    *models.UserResponse `json:"data"`
	}

	utils.RespondJSON(w, http.StatusOK, responseUser{
		Message: "Successfully fetch user data",
		Data:    resp,
	}, u.log)

	u.log.Info().Str("userId", userId).Msg("successfully send user data response")
}
