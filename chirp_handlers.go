package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SamiZeinsAI/web-server/internal/auth"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	_, userID, err := auth.AuthenticateUser(r.Header, cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s\n", err))
		return
	}
	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s\n", err))
		return
	}
	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
	chirp, ok := dbStructure.Chirps[chirpID]
	if !ok || chirp.AuthorID != userID {
		RespondWithError(w, 403, "User is not author of this shirp")
		return
	}
	delete(dbStructure.Chirps, chirpID)
	w.WriteHeader(200)

}

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {

	chirpsSlice, err := cfg.DB.GetChirps()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}

	dat, err := json.Marshal(chirpsSlice)
	if err != nil {
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func (cfg *apiConfig) PostChirpHandler(w http.ResponseWriter, r *http.Request) {

	type returnVals struct {
		Body     string `json:"body"`
		Id       int    `json:"id"`
		AuthorID int    `json:"author_id"`
	}
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, 500, "Error when decoding request")
		return
	}
	_, userID, err := auth.AuthenticateUser(r.Header, cfg.jwtSecret)
	if err != nil {
		msg := fmt.Sprintf("Error authenticating user: %s\n", err)
		RespondWithError(w, 400, msg)
		return
	}

	chirp, err := cfg.DB.CreateChirp(params.Body, userID)
	if err != nil {
		msg := fmt.Sprintf("Error creating chirp: %s\n", err)
		RespondWithError(w, 500, msg)
		return
	}

	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		msg := fmt.Sprintf("Error loading database: %s\n", err)
		RespondWithError(w, 500, msg)
		return
	}
	dbStructure.Chirps[len(dbStructure.Chirps)+1] = chirp
	cfg.DB.WriteDB(dbStructure)
	RespondWithJSON(w, 201, returnVals{
		Id:       len(dbStructure.Chirps),
		Body:     params.Body,
		AuthorID: userID,
	})
}

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {

	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("%s\n", err))
		return
	}
	val, exists := dbStructure.Chirps[id]
	if !exists {
		RespondWithError(w, 404, fmt.Sprintln("The database doesn't contain that id"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(val)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("%s\n", err))
		return
	}
}
