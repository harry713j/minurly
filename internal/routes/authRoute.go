package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/middleware"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/google/login", handler.HandleLogin).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/callback", handler.HandleLoginCallback).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/logout", middleware.VerifyLogin(handler.HandleLogout)).Methods(http.MethodGet)
	r.HandleFunc("/auth/status", middleware.VerifyLogin(handler.HandleGetAuthStatus)).Methods(http.MethodGet)
}
