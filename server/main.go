package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/controller"
	"github.com/harry713j/minurly/db"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	router := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{os.Getenv("ALLOWED_ORIGIN")}),               // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allow specific methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Allow specific headers
		handlers.AllowCredentials(), // Allow credentials (cookies, authorization headers)
	)

	subroute := router.PathPrefix("/api/v1").Subrouter()

	subroute.HandleFunc("/short", controller.HandleCreateUrl).Methods(http.MethodPost)
	subroute.HandleFunc("/{short-code}", controller.HandleGetUrl).Methods(http.MethodGet)
	subroute.HandleFunc("/{short-code}", controller.HandleDeleteUrl).Methods(http.MethodDelete)

	handler := handlers.LoggingHandler(os.Stdout, cors(router))
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("🚀 Server has started on Port %v\n", port)
	log.Fatal(server.ListenAndServe())
	defer db.Client.Disconnect(context.TODO())
}
