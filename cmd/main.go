package main

import (
	"log"
	"net/http"

	"github.com/g3orge/BasicAuthRest/internal/config"
	"github.com/g3orge/BasicAuthRest/internal/handlers"
	"github.com/g3orge/BasicAuthRest/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.MustLoad()

	_ = cfg
	storage, err := storage.New()
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	_ = storage

	r := mux.NewRouter()

	r.HandleFunc("/login/{guid}", handlers.CreateTokens(storage)).Methods("GET")
	r.HandleFunc("/refresh/{guid}/{refresh}", handlers.RefreshToken(storage)).Methods("GET")
	r.HandleFunc("/", handlers.CreateUser(storage)).Methods("POST")

	http.ListenAndServe(":8080", r)
}
