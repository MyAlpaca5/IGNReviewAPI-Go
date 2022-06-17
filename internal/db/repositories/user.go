package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	pool *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) User {
	return User{
		pool: pool,
	}
}

func (u User) Create(m models.User) (int, error) {
	query := `
	INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3)
	RETURNING id`

	var id int
	args := []interface{}{m.Username, m.Password, m.Email}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if e, ok := err.(*pgconn.PgError); ok && e.Code == "23505" {
			return 0, errors.New("unique constraint violation, possibly duplicate 'username'")
		}
		return 0, err
	}

	return id, nil
}

func (u User) ReadByUsername(username string) (models.User, error) {
	query := `
	SELECT id, created_at, updated_at, username, password, email
	FROM users
	WHERE username = $1`

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := u.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// func (u User) UpdateByID(pool *pgxpool.Pool, m models.User, id int) error {
// 	// updated_at here is used as a simple locking mechanism to prevent data race
// 	query := `
// 	UPDATE users
// 	SET username=$1, password=$2, email=$3, updated_at=now()
// 	WHERE id = $4 and updated_at = $5`

// 	args := []interface{}{m.Username, m.PasswordHash, m.Email, m.ID, m.UpdatedAt}
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	commandTag, err := pool.Exec(ctx, query, args...)
// 	if err != nil {
// 		return err
// 	}

// 	if commandTag.RowsAffected() == 0 {
// 		return errors.New("record not found")
// 	}

// 	return nil
// }
