package database

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (db *DB) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type returnVals struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s", err))
		return
	}
	params := parameters{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s", err))
		return
	}
	user, err := db.CreateUser(params.Email)
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
	dat, err = json.Marshal(respBody)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		Email: email,
		Id:    id,
	}
	return user, nil

}
