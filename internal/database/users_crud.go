package database

import (
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		Email:    email,
		Id:       id,
		Password: hash,
	}
	return user, nil
}
