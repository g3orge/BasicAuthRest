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

	storage, err := storage.New(cfg.Db.DbPath, cfg.Db.DbName, cfg.Db.CollName)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	r := mux.NewRouter()
	// user struct in user.go
	r.HandleFunc("/login/{guid}", handlers.CreateTokens(storage)).Methods("GET")             //generate access and refresh token, hash refreshT and save to db
	r.HandleFunc("/refresh/{guid}/{refresh}", handlers.RefreshToken(storage)).Methods("GET") //regenerete refresh and access tokens and save hashedRt to db
	r.HandleFunc("/", handlers.CreateUser(storage)).Methods("POST")                          //send user to database

	http.ListenAndServe(":8080", r)
}
