package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	us "github.com/g3orge/BasicAuthRest/internal/user"
	"github.com/gorilla/mux"
)

type Store interface {
	Create(ctx context.Context, user us.User) (string, error)
	GenerateTokens(ctx context.Context, guid string) (string, string, error)
	RefreshToken(ctx context.Context, rt string, guid string) (string, string, error)
}

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func CreateUser(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}

		var u us.User

		if err = json.Unmarshal(body, &u); err != nil {
			log.Println(err)
		}

		id, err1 := s.Create(context.Background(), u)
		if err1 != nil {
			log.Println(err1)
		}

		log.Println(id)
	}
}

func CreateTokens(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		guid := params["guid"]

		at, rt, err := s.GenerateTokens(context.Background(), guid)
		if err != nil {
			log.Println(err)
		}

		resp := AuthTokenResponse{
			AccessToken:  at,
			RefreshToken: rt,
		}

		pl, err1 := json.Marshal(resp)
		if err1 != nil {
			log.Println(err1)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pl)
	}
}

func RefreshToken(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		guid := params["guid"]
		hRT := params["refresh"]

		at, rt, err := s.RefreshToken(context.Background(), hRT, guid)
		if err != nil {
			log.Println(err)
		}

		resp := AuthTokenResponse{
			AccessToken:  at,
			RefreshToken: rt,
		}

		pl, err1 := json.Marshal(resp)
		if err1 != nil {
			log.Println(err1)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pl)
	}
}
