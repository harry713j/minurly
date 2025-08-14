package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/controller"
	"github.com/harry713j/minurly/middleware"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/google/login", controller.HandleLogin).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/callback", controller.HandleLoginCallback).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/logout", middleware.VerifyLogin(controller.HandleLogout)).Methods(http.MethodGet)
	r.HandleFunc("/auth/status", middleware.VerifyLogin(controller.HandleGetAuthStatus)).Methods(http.MethodGet)
}
