package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/middleware"
)

func registerAuthRoutes(r *mux.Router, h *handler.Handler, mw *middleware.Middleware) {
	r.HandleFunc("/auth/google/login", h.Auth.HandleLogin).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/callback", h.Auth.HandleLoginCallback).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/logout", mw.Auth.VerifySession(h.Auth.HandleLogout)).Methods(http.MethodGet)
}
