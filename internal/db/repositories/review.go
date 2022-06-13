package repositories

import (
	"context"
	"errors"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r ReviewRepo) Read(pool *pgxpool.Pool, id int) (models.Review, error) {
	query := `
	SELECT name, description, created_at, updated_at, review_url, review_score, media_type, genre_list, creator_list, version
	FROM reviews
	WHERE id = $1`

	var review models.Review
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := pool.QueryRow(ctx, query, id).Scan(&review.Name, &review.Description, &review.CreatedAt, &review.UpdatedAt, &review.ReviewURL, &review.ReviewScore, &review.MediaType, &review.GenreList, &review.CreatorList, &review.Version)
	if err != nil {
		return models.Review{}, err
	}

	return review, nil
}

func (r ReviewRepo) ReadAll(pool *pgxpool.Pool) ([]models.Review, error) {
	// query := `
	// SELECT *
	// FROM reviews`

	return nil, nil
}

func (r ReviewRepo) Update(pool *pgxpool.Pool, id int, m models.Review) error {
	// version here is used as a simple locking mechanism to prevent data race
	query := `
	UPDATE reviews 
	SET name=$1, description=$2, updated_at=$3, review_score=$4, media_type=$5, genre_list=$6, creator_list=$7, version = version + 1
	WHERE id = $8 and version = $9`

	args := []interface{}{m.Name, m.Description, m.UpdatedAt, m.ReviewScore, m.MediaType, m.GenreList, m.CreatorList, id, m.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	commandTag, err := pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (r ReviewRepo) Delete(pool *pgxpool.Pool, id int) error {
	query := `
	DELETE FROM reviews
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	commandTag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("record not found")
	}

	return nil
}
