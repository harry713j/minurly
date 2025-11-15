package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/harry713j/minurly/internal/middleware"
	"github.com/harry713j/minurly/internal/models"
	"github.com/harry713j/minurly/internal/service"
	"github.com/harry713j/minurly/internal/utils"
	"github.com/harry713j/minurly/internal/validation"
	"github.com/rs/zerolog"
)

type UrlHandler struct {
	services *service.Service
	log      zerolog.Logger
}

type shortUrl struct {
	ShortCode string `json:"shortCode"`
}

func NewUrlHandler(logger zerolog.Logger, services *service.Service) *UrlHandler {
	return &UrlHandler{
		log:      logger,
		services: services,
	}
}

func (u *UrlHandler) HandleCreateUrl(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserId(r)

	if !ok {
		u.log.Warn().Str("userId", userId).Msg("userId value not found from the request context")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
		return
	}

	var req models.CreateUrlPayload
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		u.log.Err(err).Str("func_name", "HandleCreateUrl").Str("userId", userId).Msg("failed to decode the request body")
		utils.RespondError(w, apperrors.NewInternalServerErr())
		return
	}
	// validate the incoming struct
	if err = validation.ValidateStruct(req); err != nil {
		u.log.Warn().Str("func_name", "HandleCreateUrl").Str("userId", userId).Msg("invalid structure")
		utils.RespondError(w, apperrors.NewBadRequestErr("invalid data provided, required an url"))
		return
	}

	url, err := u.services.Url.Create(userId, req.OriginalUrl)
	if err != nil {
		u.log.Err(err).Str("func_name", "HandleCreateUrl").Str("userId", userId).Msg("failed to create short url")
		utils.RespondError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, shortUrl{
		ShortCode: url.ShortCode,
	}, u.log)

	u.log.Info().Str("userId", userId).Msg("successfully send created short url response")
}

func (u *UrlHandler) HandleGetUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode, ok := vars["short-code"]

	if !ok {
		u.log.Warn().Str("short-code", shortCode).Msg("no short-code parameter in request url")
		utils.RespondError(w, apperrors.NewBadRequestErr("invalid params"))
		return
	}

	if shortCode == "" || len(shortCode) != 8 {
		u.log.Warn().Str("short-code", shortCode).Msg("short-code is not valid")
		utils.RespondError(w, apperrors.NewBadRequestErr("invalid params"))
		return
	}

	url, err := u.services.Url.Fetch(shortCode)
	if err != nil {
		u.log.Err(err).Str("short-code", shortCode).Msg("failed to get url ")
		utils.RespondError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message":     "Successfully fetch original url",
		"originalUrl": url.OriginalUrl,
	}, u.log)

	u.log.Info().Str("short-code", shortCode).Msg("successfully send original url response")
}

func (u *UrlHandler) HandleDeleteUrl(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserId(r)

	if !ok {
		u.log.Warn().Str("userId", userId).Msg("userId value not found from the request context")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
		return
	}

	vars := mux.Vars(r)
	shortCode, ok := vars["short-code"]

	if !ok {
		u.log.Warn().Str("short-code", shortCode).Msg("no short-code parameter in request url")
		utils.RespondError(w, apperrors.NewBadRequestErr("invalid params"))
		return
	}

	if err := u.services.Url.Remove(userId, shortCode); err != nil {
		u.log.Err(err).Str("userId", userId).Str("short-code", shortCode).Msg("failed to delete the url")
		utils.RespondError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Successfully deleted url"}, u.log)
	u.log.Info().Str("userId", userId).Str("short-code", shortCode).Msg("successfully send delete url response")
}
