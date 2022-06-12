package repositories

import (
	"fmt"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReviewRepo struct{}

func (r ReviewRepo) Create(pool *pgxpool.Pool, m models.Review) error {
	fmt.Printf("Called ReviewRepo Create! TotalConns is %v", pool.Stat().TotalConns())
	return nil
}

func (r ReviewRepo) Read(pool *pgxpool.Pool, id int) (models.Review, error) {
	fmt.Printf("Called ReviewRepo Read! TotalConns is %v", pool.Stat().TotalConns())
	return models.Review{}, nil
}

func (r ReviewRepo) ReadAll(pool *pgxpool.Pool) ([]models.Review, error) {
	fmt.Printf("Called ReviewRepo ReadAll! TotalConns is %v", pool.Stat().TotalConns())
	return []models.Review{}, nil
}

func (r ReviewRepo) Update(pool *pgxpool.Pool, id int, m models.Review) error {
	fmt.Printf("Called ReviewRepo Update! TotalConns is %v", pool.Stat().TotalConns())
	return nil
}

func (r ReviewRepo) Delete(pool *pgxpool.Pool, id int) error {
	fmt.Printf("Called ReviewRepo Delete! TotalConns is %v", pool.Stat().TotalConns())
	return nil
}
