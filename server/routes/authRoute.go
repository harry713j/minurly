package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/controller"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/google/login", controller.HandleLogin).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/callback", controller.HandleLoginCallback).Methods(http.MethodGet)
	r.HandleFunc("/auth/google/logout", controller.HandleLogout).Methods(http.MethodGet)
}
