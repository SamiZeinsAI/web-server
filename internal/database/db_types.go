package database

import (
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Chirps        map[int]Chirp `json:"chirps"`
	Users         map[int]User  `json:"users"`
	RevokedTokens map[string]RevokedToken
}

type RevokedToken struct {
	TimeRevoked time.Time
	TokenString string
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}
type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}
