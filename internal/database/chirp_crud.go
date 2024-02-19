package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})

}
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func cleanInput(input string) string {
	words := strings.Split(input, " ")
	for i, word := range words {
		word = strings.ToLower(word)
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (db *DB) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {

	chirpsSlice, err := db.GetChirps()
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

func (db *DB) PostChirpHandler(w http.ResponseWriter, r *http.Request) {

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
	chirp, err := db.CreateChirp(params.Body)
	if err != nil {
		msg := fmt.Sprintf("Error creating chirp: %s\n", err)
		RespondWithError(w, 500, msg)
	}
	dbStructure, err := db.loadDB()
	if err != nil {
		msg := fmt.Sprintf("Error loading database: %s\n", err)
		RespondWithError(w, 500, msg)
	}
	dbStructure.Chirps[len(dbStructure.Chirps)+1] = chirp
	db.writeDB(dbStructure)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(returnVals{
		Id:   len(dbStructure.Chirps),
		Body: params.Body,
	})
	if err != nil {
		msg := fmt.Sprintf("Error encoding response: %s\n", err)
		RespondWithError(w, 500, msg)
		return
	}
}

func (db *DB) GetChirpHandler(w http.ResponseWriter, r *http.Request) {

	dbStructure, err := db.loadDB()
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

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	if len(body) > 140 {
		return Chirp{}, errors.New("Chirp is too long")
	}
	return Chirp{
		Body: cleanInput(body),
		Id:   len(dbStructure.Chirps) + 1,
	}, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	for _, v := range dbStructure.Chirps {
		chirps = append(chirps, v)
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	return chirps, nil
}
