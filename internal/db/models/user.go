package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User struct includes all data related to one user record
type User struct {
	dbBase
	Username string
	Password []byte
	Email    string
	Role     int
}

func (u User) ComparePassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, errors.New("mismatched password")
		default:
			return false, err
		}
	}
	return true, nil
}
