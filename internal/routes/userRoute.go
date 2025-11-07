package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/middleware"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/users", middleware.VerifyLogin(handler.HandleGetUser)).Methods(http.MethodGet)
}
