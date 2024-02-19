package database

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) PostUserLoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type returnVals struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s\n", err))
		return
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
	matchingUser := User{}
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
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	respBody := returnVals{
		Id:    matchingUser.Id,
		Email: matchingUser.Email,
	}
	err = json.NewEncoder(w).Encode(&respBody)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
}

func (db *DB) PostUserHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := db.CreateUser(params.Email, params.Password)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	dbStructure, err := db.loadDB()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	id := len(dbStructure.Users) + 1
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
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

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		Email:    email,
		Id:       id,
		Password: hash,
	}
	return user, nil
}
