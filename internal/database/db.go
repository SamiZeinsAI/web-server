package database

import (
	"encoding/json"
	"os"
	"sync"
)

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.EnsureDB()
	if err != nil {
		return nil, err
	}
	err = db.WriteDB(DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RevokedTokens: map[string]RevokedToken{},
	})
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) EnsureDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RevokedTokens: map[string]RevokedToken{},
	}
	dat, err := os.ReadFile(db.path)
	if err != nil {
		db.WriteDB(dbStructure)
		return nil
	}

	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		db.WriteDB(DBStructure{
			Chirps:        map[int]Chirp{},
			Users:         map[int]User{},
			RevokedTokens: map[string]RevokedToken{},
		})
		return nil
	}
	return nil
}

func (db *DB) LoadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dat, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	dbStructure := DBStructure{}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}
	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) WriteDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	os.WriteFile(db.path, dat, 0666)
	return nil
}
