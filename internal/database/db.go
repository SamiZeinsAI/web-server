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
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	})
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) EnsureDB() error {
	dat, err := os.ReadFile(db.path)
	if err != nil {
		db.WriteDB(DBStructure{
			Chirps: map[int]Chirp{},
			Users:  map[int]User{},
		})
		return nil
	}
	dbStruc := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	err = json.Unmarshal(dat, &dbStruc)
	if err != nil {
		db.WriteDB(DBStructure{
			Chirps: map[int]Chirp{},
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
	dbStruc := DBStructure{}
	err = json.Unmarshal(dat, &dbStruc)
	if err != nil {
		return DBStructure{}, err
	}
	return dbStruc, nil
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
