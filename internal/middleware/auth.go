package middleware

import (
	"context"
	"net/http"

	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/harry713j/minurly/internal/server"
	"github.com/harry713j/minurly/internal/utils"
)

type AuthMiddleware struct {
	server *server.Server
}

type contextKey string

const userIdKey contextKey = "userId"

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

func (a *AuthMiddleware) VerifySession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session from the store
		session, err := a.server.Config.Auth.SessionStore.Get(r, "session")
		if err != nil {
			a.server.LoggerService.Logger.Warn().Str("err", err.Error()).Msg("no session cookie found from request")
			utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
			return
		}

		userId, ok := session.Values["userId"].(string)

		if !ok || userId == "" {
			a.server.LoggerService.Logger.Warn().Msg("userId not found from request")
			utils.RespondError(w, apperrors.NewUnauthorizedErr("Unauthorized"))
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userId)

		a.server.LoggerService.Logger.Info().Str("userId", userId).Msg("set the context from auth middleware")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserId(r *http.Request) (string, bool) {
	userId, ok := r.Context().Value(userIdKey).(string)

	return userId, ok
}
