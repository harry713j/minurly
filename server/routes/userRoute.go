package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/controller"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/users", controller.HandleGetUser).Methods(http.MethodGet)
}
