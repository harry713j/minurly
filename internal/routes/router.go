package routes

import (
	"net/http"

	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/handler"
	"github.com/harry713j/minurly/internal/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Config, h *handler.Handler, mw *middleware.Middleware) http.Handler {
	router := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedOrigins(cfg.Server.CORSAllowedOrigins),                       // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allow specific methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Allow specific headers
		handlers.AllowCredentials(),                                                  // Allow credentials (cookies, authorization headers)
	)

	// add all the routes here

	subroute := router.PathPrefix("/api/v1").Subrouter()
	registerAuthRoutes(subroute, h, mw)
	registerUserRoutes(subroute, h, mw)
	registerUrlRoutes(subroute, h, mw)

	handler := cors(router)

	return handler
}
