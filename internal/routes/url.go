package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/middleware"
)

func registerUrlRoutes(r *mux.Router, h *handler.Handler, mw *middleware.Middleware) {
	r.HandleFunc("/urls", mw.Auth.VerifySession(h.Url.HandleCreateUrl)).Methods(http.MethodPost)
	r.HandleFunc("/urls/{short-code}", h.Url.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/urls/{short-code}", mw.Auth.VerifySession(h.Url.HandleDeleteUrl)).Methods(http.MethodDelete)
}
