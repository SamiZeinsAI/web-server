package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SamiZeinsAI/web-server/internal/auth"
	"github.com/SamiZeinsAI/web-server/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) PostLoginHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s\n", err))
		return
	}

	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
	matchingUser := database.User{}
	foundMatch := false
	for _, user := range dbStructure.Users {
		if user.Email == params.Email {
			err := bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
			if err != nil {
				RespondWithError(w, 401, fmt.Sprintf("%s\n", err))
				return
			}
			matchingUser = user
			foundMatch = true
			break
		}
	}
	if !foundMatch {
		RespondWithError(w, 401, fmt.Sprintf("%s\n", err))
		return
	}

	token, err := auth.MakeToken(matchingUser.Id, "chirpy-access", cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
	refreshToken, err := auth.MakeToken(matchingUser.Id, "chirpy-refresh", cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	respBody := returnVals{
		Id:           matchingUser.Id,
		Email:        matchingUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
	err = json.NewEncoder(w).Encode(&respBody)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
}
