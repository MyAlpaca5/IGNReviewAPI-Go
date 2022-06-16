package models

import (
	"errors"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// UserIn struct includes data that is received from a request
// Used pointer here to differentiate from zero value and no input given
type UserIn struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

// User struct includes all data related to one user record
type User struct {
	dbBase
	UserOut
	PasswordHash []byte
}

// UserOut struct includes data that is sent back via a response
type UserOut struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *UserIn) Validate() error {
	if u.Username != nil {
		if *u.Username == "" {
			return errors.New("field 'username' cannot be empty string")
		}

		if strings.ContainsAny(*u.Username, " ") {
			return errors.New("field 'username' cannot contain white space")
		}
	}

	if u.Password != nil {
		pwd := []byte(*u.Password)
		if len(pwd) < 8 {
			return errors.New("field 'password' is too short")
		}

		if len(pwd) > 72 {
			return errors.New("field 'password' is too long")
		}
		// TODO: more password strength check
	}

	if u.Email != nil {
		_, err := mail.ParseAddress(*u.Email)
		if err != nil {
			return errors.New("field 'email' is not a valid email address")
		}
	}

	return nil
}

func (u *UserIn) ToUser() (User, error) {
	var user User
	if u.Username != nil {
		user.Username = *u.Username
	}
	if u.Password != nil {
		err := user.SetPassword(*u.Password)
		if err != nil {
			return User{}, err
		}
	}
	if u.Email != nil {
		user.Email = *u.Email
	}

	return user, nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	u.PasswordHash = hash
	return nil
}

func (u User) ComparePassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
