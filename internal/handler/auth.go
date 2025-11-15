package handler

import (
	"net/http"

	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/middleware"
	"github.com/harry713j/minurly/internal/service"
	"github.com/harry713j/minurly/internal/utils"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	cfg      *config.Config
	services *service.Service
	log      zerolog.Logger
}

func NewAuthHandler(cfg *config.Config, services *service.Service, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		cfg:      cfg,
		services: services,
		log:      logger,
	}
}

func (a *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateRandomStrings(16)

	if err != nil {
		a.log.Err(err).Msg("failed to generate the random string for oauth")
		utils.RespondError(w, apperrors.NewInternalServerErr())
		return
	}

	session, err := a.cfg.Auth.SessionStore.Get(r, "oauthstate")
	if err != nil {
		a.log.Debug().Err(err).Msg("oauthstate cookie not found, creating new session")
	}

	session.Values["state"] = state
	if err = session.Save(r, w); err != nil {
		a.log.Error().Err(err).Msg("failed to save the session")
		utils.RespondError(w, apperrors.NewInternalServerErr())
		return
	}

	// url will be created by oauth client
	// look for this state
	url := a.cfg.Auth.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (a *AuthHandler) HandleLoginCallback(w http.ResponseWriter, r *http.Request) {
	session, err := a.cfg.Auth.SessionStore.Get(r, "oauthstate")
	if err != nil {
		a.log.Warn().Str("err", err.Error()).Msg("cookie named oauthstate not found from client request")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
		return
	}

	storedState, ok := session.Values["state"].(string)
	if !ok {
		a.log.Warn().Msg("state not found from the cookie")
		utils.RespondError(w, apperrors.NewBadRequestErr("invalid session state"))
		return
	}

	// Read state from Google callback
	returnedState := r.URL.Query().Get("state")
	if returnedState == "" || returnedState != storedState {
		a.log.Warn().Str("oauth_state", returnedState).Str("session_state", storedState).Msg("state not matching")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("invalid oauth state"))
		return
	}

	code := r.URL.Query().Get("code") // oauth client send a ?code=random in query
	if code == "" {
		a.log.Warn().Msg("empty code from oauth")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("invalid oauth code"))
		return
	}

	user, err := a.services.User.Create(code)
	if err != nil {
		a.log.Warn().Str("err", err.Error()).Msg("user creation failed ")
		utils.RespondError(w, err)
		return
	}

	// create the session
	sessionDoc, err := a.services.Auth.CreateSession(user.ID)
	if err != nil {
		a.log.Warn().Str("err", err.Error()).Msg("session creation failed")
		utils.RespondError(w, err)
		return
	}

	appSession, err := a.cfg.Auth.SessionStore.Get(r, "session")
	if err != nil {
		a.log.Debug().Err(err).Msg("session cookie not found, creating new session")
	}

	appSession.Values["sessionId"] = sessionDoc.SessionId
	appSession.Values["userId"] = user.ID.Hex()

	if err := appSession.Save(r, w); err != nil {
		a.log.Error().Err(err).Msg("failed to save the session")
		utils.RespondError(w, apperrors.NewInternalServerErr())
		return
	}

	http.Redirect(w, r, a.cfg.Server.CORSAllowedOrigins[0], http.StatusSeeOther)
}

func (a *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserId(r)

	if !ok {
		a.log.Warn().Msg("userId value not found from the request context")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
		return
	}

	err := a.services.Auth.DeleteSession(userId)
	if err != nil {
		a.log.Err(err).Msg("failed to delete the session")
		utils.RespondError(w, err)
		return
	}

	session, err := a.cfg.Auth.SessionStore.Get(r, "session")
	if err != nil {
		a.log.Warn().Msg("session cookie not found")
		utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
		return
	}

	session.Options.MaxAge = -1 // Marks it as expired
	if err = session.Save(r, w); err != nil {
		a.log.Err(err).Msg("failed to remove the session cookie")
		utils.RespondError(w, apperrors.NewInternalServerErr())
		return
	}

	// Clear cookie explicitly
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Log out successful"}, a.log)

	a.log.Info().Str("userId", userId).Msg("successfully logout the user")
}
