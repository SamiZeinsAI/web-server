package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
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
	db.mux.Lock()
	defer db.mux.Unlock()
	chirpsSlice, err := db.GetChirps()
	fmt.Println(chirpsSlice)
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
	db.mux.Lock()
	defer db.mux.Unlock()
	type returnVals struct {
		Body string `json:"body"`
		Id   int    `json:"id"`
	}
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 500, "Error when decoding request")
		return
	}
	chirp, err := db.CreateChirp(params.Body)
	if err != nil {
		msg := fmt.Sprintf("Error creating chirp: %s\n", err)
		RespondWithError(w, 500, msg)
	}
	dbStruct, err := db.loadDB()
	if err != nil {
		msg := fmt.Sprintf("Error loading database: %s\n", err)
		RespondWithError(w, 500, msg)
	}
	dbStruct.Chirps[db.NextId] = chirp
	db.writeDB(dbStruct)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(returnVals{
		Id:   db.NextId,
		Body: params.Body,
	})
	if err != nil {
		msg := fmt.Sprintf("Error encoding response: %s\n", err)
		RespondWithError(w, 500, msg)
		return
	}
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	if len(body) > 140 {
		return Chirp{}, errors.New("Chirp is too long")
	}
	cleanedChirp := cleanInput(body)
	db.NextId++
	return Chirp{
		Body: cleanedChirp,
		Id:   db.NextId,
	}, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	for _, v := range dbStruct.Chirps {
		chirps = append(chirps, v)
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	return chirps, nil
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path:   path,
		mux:    &sync.RWMutex{},
		NextId: 0,
	}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	err = db.writeDB(DBStructure{
		Chirps: map[int]Chirp{},
	})
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	err := os.WriteFile(db.path, []byte(""), 0666)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dat, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	dbStruc := DBStructure{}
	err = json.Unmarshal(dat, &dbStruc)
	if err != nil {
		return DBStructure{}, err
	}
	return dbStruc, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	os.WriteFile(db.path, dat, 0666)
	return nil
}
