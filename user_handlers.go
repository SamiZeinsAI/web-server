package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SamiZeinsAI/web-server/internal/auth"
)

func (cfg *apiConfig) PutUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Couldn't find the JWT")
		return
	}

	token, err := auth.ParseToken(tokenString, cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf("%s\n", err))
		return
	}
	issuer, err := token.Claims.GetIssuer()
	fmt.Printf("%s\n", issuer)
	if err != nil || issuer == "chirpy-refresh" {
		RespondWithError(w, 401, "Error finding valid access token")
		return
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}

	cfg.DB.UpdateUser(idInt, params.Email, params.Password)
	err = json.NewEncoder(w).Encode(returnVals{
		Email: params.Email,
		Id:    idInt,
	})
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
}

func (cfg *apiConfig) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s", err))
		return
	}
	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	id := len(dbStructure.Users) + 1
	dbStructure.Users[id] = user
	err = cfg.DB.WriteDB(dbStructure)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}

	respBody := returnVals{
		Id:    id,
		Email: params.Email,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}
