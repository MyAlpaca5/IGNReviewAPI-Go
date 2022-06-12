package repositories

import (
	"context"
	"fmt"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReviewRepo struct{}

func (r ReviewRepo) Create(pool *pgxpool.Pool, m models.Review) (int, error) {
	query := `
	INSERT INTO reviews (name, description, created_at, updated_at, review_url, review_score, media_type, genre_list, creator_list)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`

	var id int
	args := []interface{}{m.Name, m.Description, m.CreatedAt, m.UpdatedAt, m.ReviewURL, m.ReviewScore, m.MediaType, m.GenreList, m.CreatorList}
	err := pool.QueryRow(context.Background(), query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r ReviewRepo) Read(pool *pgxpool.Pool, id int) (models.Review, error) {
	query := `
	SELECT name, description, created_at, updated_at, review_url, review_score, media_type, genre_list, creator_list
	FROM reviews
	WHERE id = $1`

	var review models.Review
	err := pool.QueryRow(context.Background(), query, id).Scan(&review.Name, &review.Description, &review.CreatedAt, &review.UpdatedAt, &review.ReviewURL, &review.ReviewScore, &review.MediaType, &review.GenreList, &review.CreatorList)
	if err != nil {
		return models.Review{}, err
	}

	return review, nil
}

func (r ReviewRepo) ReadAll(pool *pgxpool.Pool) ([]models.Review, error) {
	// query := `
	// SELECT name, description, created_at, updated_at, review_url, review_score, media_type, genre_list, creator_list
	// FROM reviews`

	return nil, nil
}

func (r ReviewRepo) Update(pool *pgxpool.Pool, id int, m models.Review) error {
	fmt.Printf("Called ReviewRepo Update! TotalConns is %v", pool.Stat().TotalConns())
	return nil
}

func (r ReviewRepo) Delete(pool *pgxpool.Pool, id int) error {
	fmt.Printf("Called ReviewRepo Delete! TotalConns is %v", pool.Stat().TotalConns())
	return nil
}
