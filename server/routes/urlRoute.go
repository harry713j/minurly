package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/controller"
)

func UrlRoutes(r *mux.Router) {
	r.HandleFunc("/urls", controller.HandleCreateUrl).Methods(http.MethodPost)
	r.HandleFunc("/urls/{short-code}", controller.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/urls/{short-code}", controller.HandleDeleteUrl).Methods(http.MethodDelete)
}
