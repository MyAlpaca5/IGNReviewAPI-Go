package pasetotoken

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type tokenPayload struct {
	UserID int
	Expiry time.Time
	Role   int
}

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return PasetoMaker{}, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker PasetoMaker) CreateToken(userID int, expiry time.Time, role int) (models.Token, string, error) {
	var payload = &tokenPayload{
		UserID: userID,
		Expiry: expiry,
		Role:   role,
	}

	tokenStr, err := maker.paseto.Encrypt(maker.symmetricKey, &payload, nil)
	if err != nil {
		return models.Token{}, "", err
	}

	var token = models.Token{
		TokenHash: tokenStrToHash(tokenStr),
		UserID:    payload.UserID,
		Expiry:    payload.Expiry,
		Role:      payload.Role,
	}

	return token, tokenStr, nil
}

func (maker PasetoMaker) VerifyToken(tokenStr string) (models.Token, error) {
	payload := &tokenPayload{}

	err := maker.paseto.Decrypt(tokenStr, maker.symmetricKey, payload, nil)
	if err != nil {
		return models.Token{}, ErrInvalidToken
	}

	if time.Now().UTC().After(payload.Expiry) {
		return models.Token{}, ErrExpiredToken
	}

	var token = models.Token{
		TokenHash: tokenStrToHash(tokenStr),
		UserID:    payload.UserID,
		Expiry:    payload.Expiry,
		Role:      payload.Role,
	}

	return token, nil
}

func tokenStrToHash(tokenStr string) []byte {
	tokenHash := sha256.Sum256([]byte(tokenStr))
	return tokenHash[:]
}
