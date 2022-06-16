package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Token struct{}

func (t Token) Create(pool *pgxpool.Pool, m models.Token) error {
	query := `
	INSERT INTO tokens (token, userID, expiry, role)
	VALUES ($1, $2, $3, $4)	
	`
	args := []interface{}{m.TokenHash, m.UserID, m.Expiry, m.Role}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	commandTag, err := pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("record not create")
	}

	return nil
}

func (t Token) ReadByToken(pool *pgxpool.Pool, tokenHash []byte) (models.Token, error) {
	query := `
	SELECT userid, expiry, role
	FROM tokens
	WHERE token = $1`

	var token = models.Token{TokenHash: tokenHash}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := pool.QueryRow(ctx, query, tokenHash).Scan(&token.UserID, &token.Expiry, &token.Role)
	if err != nil {
		return models.Token{}, err
	}

	return token, nil
}
