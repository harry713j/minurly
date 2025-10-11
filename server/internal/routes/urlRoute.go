package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/middleware"
)

func UrlRoutes(r *mux.Router) {
	r.HandleFunc("/urls", middleware.VerifyLogin(handler.HandleCreateUrl)).Methods(http.MethodPost)
	r.HandleFunc("/urls/{short-code}", handler.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/urls/{short-code}", middleware.VerifyLogin(handler.HandleDeleteUrl)).Methods(http.MethodDelete)
}
