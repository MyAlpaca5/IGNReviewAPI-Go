package schemas

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"golang.org/x/crypto/bcrypt"
)

// UserIn struct includes data that is received from a request
// Used pointer here to differentiate from zero value and no input given
type UserIn struct {
	Username    *string `json:"username"`
	PasswordStr *string `json:"password"`
	Email       *string `json:"email"`
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

	if u.PasswordStr != nil {
		pwd := []byte(*u.PasswordStr)
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

// ToUser creates a User model based on the data in a UserIn schema
func (u UserIn) ToUser() (models.User, error) {
	var user models.User

	if u.Username != nil {
		user.Username = *u.Username
	}
	if u.PasswordStr != nil {
		hash, err := generatePasswordHash(*u.PasswordStr)
		if err != nil {
			return models.User{}, err
		}
		user.Password = hash
	}
	if u.Email != nil {
		user.Email = *u.Email
	}

	return user, nil
}

func generatePasswordHash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
