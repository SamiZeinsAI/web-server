package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) UpdateUser(id int, newEmail, newPassword string) error {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	_, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("ID doesn't exist in database")
	}
	user, err := db.CreateUser(newEmail, newPassword)
	if err != nil {
		return err
	}
	user.Id = id
	dbStructure.Users[id] = user
	db.WriteDB(dbStructure)
	return nil
}

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
