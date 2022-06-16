package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"
)

type Token struct {
	TokenString string    `json:"token"`
	Expiry      time.Time `json:"expiry"`
	TokenHash   []byte
	UserID      int
	Role        int
}

func NewToken(userID int, expiry time.Time, role int) (Token, error) {
	var token Token
	tokenString, tokenHash, err := generateToken()
	if err != nil {
		return Token{}, err
	}

	token.TokenString = tokenString
	token.TokenHash = tokenHash
	token.Expiry = expiry
	token.UserID = userID
	token.Role = role

	return token, nil
}

func generateToken() (string, []byte, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "", nil, errors.New("cannot create token")
	}

	tokenStr := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf)
	return tokenStr, TokenStrToHash(tokenStr), nil
}

func TokenStrToHash(tokenStr string) []byte {
	tokenHash := sha256.Sum256([]byte(tokenStr))
	return tokenHash[:]
}
