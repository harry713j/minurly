package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/middleware"
)

func registerUserRoutes(r *mux.Router, h *handler.Handler, mw *middleware.Middleware) {
	r.HandleFunc("/users/me", mw.Auth.VerifySession(h.User.HandleGetUser)).Methods(http.MethodGet)
}
