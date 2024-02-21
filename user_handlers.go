package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SamiZeinsAI/web-server/internal/database"
	"github.com/SamiZeinsAI/web-server/internal/database/auth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(key *jwt.Token) (interface{}, error) { return []byte(cfg.jwtSecret), nil },
	)
	if err != nil {
		RespondWithError(w, 401, fmt.Sprintf("%s\n", err))
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

func (cfg *apiConfig) PostUserLoginHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type returnVals struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
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
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = 86400
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Second * time.Duration(params.ExpiresInSeconds)),
		},
		Subject: fmt.Sprintf("%d", matchingUser.Id),
	})
	signedToken, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	respBody := returnVals{
		Id:    matchingUser.Id,
		Email: matchingUser.Email,
		Token: signedToken,
	}
	err = json.NewEncoder(w).Encode(&respBody)
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
