package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/controller"
	"github.com/harry713j/minurly/middleware"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/users", middleware.VerifyLogin(controller.HandleGetUser)).Methods(http.MethodGet)
}
