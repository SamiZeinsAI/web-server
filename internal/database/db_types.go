package database

import "sync"

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}
type User struct {
	Id    int    `json:"id"`
	Email string `json:"string"`
}
type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}
