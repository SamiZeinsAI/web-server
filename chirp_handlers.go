package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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
		Body string `json:"body"`
		Id   int    `json:"id"`
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
	chirp, err := cfg.DB.CreateChirp(params.Body)
	if err != nil {
		msg := fmt.Sprintf("Error creating chirp: %s\n", err)
		RespondWithError(w, 500, msg)
	}
	dbStructure, err := cfg.DB.LoadDB()
	if err != nil {
		msg := fmt.Sprintf("Error loading database: %s\n", err)
		RespondWithError(w, 500, msg)
	}
	dbStructure.Chirps[len(dbStructure.Chirps)+1] = chirp
	cfg.DB.WriteDB(dbStructure)
	RespondWithJSON(w, 201, returnVals{
		Id:   len(dbStructure.Chirps),
		Body: params.Body,
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
