package utils

import (
	"encoding/json"
	"log"
	"net/http"
)


func RespondWithJSON[T any](w http.ResponseWriter, code int, payload T){
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshaling the json for payload %v %v\n", payload, err)
		w.WriteHeader(500)
		w.Write([]byte("Server error"))
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, msg string){
	if code > 499 {
		log.Printf("Responding with 500 error: %s\n", msg)
		w.WriteHeader(500)
	}

	type errorResponse struct {
		Error string `json:"error"`
	} 

	RespondWithJSON(w, code, errorResponse {
		Error: msg,
	})
}