package router

import (
	"net/http"

	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Config) http.Handler {
	router := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedOrigins(cfg.Server.CORSAllowedOrigins),                       // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allow specific methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Allow specific headers
		handlers.AllowCredentials(),                                                  // Allow credentials (cookies, authorization headers)
	)

	// add all the routes here
	subroute := router.PathPrefix("/api/v1").Subrouter()
	routes.AuthRoutes(subroute)
	routes.UserRoutes(subroute)
	routes.UrlRoutes(subroute)

	handler := cors(router)

	return handler
}
