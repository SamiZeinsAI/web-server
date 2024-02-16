package database

import "sync"

type DB struct {
	path   string
	mux    *sync.RWMutex
	NextId int
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}
type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}
