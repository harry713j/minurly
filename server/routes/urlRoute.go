package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/controller"
	"github.com/harry713j/minurly/middleware"
)

func UrlRoutes(r *mux.Router) {
	r.HandleFunc("/urls", middleware.VerifyLogin(controller.HandleCreateUrl)).Methods(http.MethodPost)
	r.HandleFunc("/urls/{short-code}", controller.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/urls/{short-code}", middleware.VerifyLogin(controller.HandleDeleteUrl)).Methods(http.MethodDelete)
}
