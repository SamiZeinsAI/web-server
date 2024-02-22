package database

import (
	"errors"
	"sort"
	"strings"
)

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

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	if len(body) > 140 {
		return Chirp{}, errors.New("Chirp is too long")
	}
	return Chirp{
		Body:     cleanInput(body),
		Id:       len(dbStructure.Chirps) + 1,
		AuthorID: authorId,
	}, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}
	dbStructure, err := db.LoadDB()
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
