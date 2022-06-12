package repositories

import (
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type model interface {
	models.Review
}

// THINK: should I just use any instead of define the model interface for type constrain?
type Repo[M model] interface {
	Create(pool *pgxpool.Pool, m M) (int, error)
	Read(pool *pgxpool.Pool, id int) (M, error)
	ReadAll(pool *pgxpool.Pool) ([]M, error)
	Update(pool *pgxpool.Pool, id int, m M) error
	Delete(pool *pgxpool.Pool, id int) error
}
