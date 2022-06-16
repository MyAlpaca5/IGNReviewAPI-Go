package repositories

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Review struct{}

func (r Review) Create(pool *pgxpool.Pool, m models.Review) (int, error) {
	query := `
	INSERT INTO reviews (name, description, review_url, review_score, media_type, genre_list, creator_list)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	var id int
	args := []interface{}{m.Name, m.Description, m.ReviewURL, m.ReviewScore, m.MediaType, m.GenreList, m.CreatorList}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r Review) ReadByID(pool *pgxpool.Pool, id int) (models.Review, error) {
	query := `
	SELECT name, description, created_at, updated_at, review_url, review_score, media_type, genre_list, creator_list
	FROM reviews
	WHERE id = $1`

	var review models.Review
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := pool.QueryRow(ctx, query, id).Scan(&review.Name, &review.Description, &review.CreatedAt, &review.UpdatedAt, &review.ReviewURL, &review.ReviewScore, &review.MediaType, &review.GenreList, &review.CreatorList)
	if err != nil {
		return models.Review{}, err
	}

	return review, nil
}

func (r Review) ReadAll(pool *pgxpool.Pool, queryParamaters map[string][]string) ([]models.Review, error) {
	whereClause, err := generateWhereClause(queryParamaters)
	if err != nil {
		return nil, err
	}

	orderByClause, err := generateOrderByClause(queryParamaters)
	if err != nil {
		return nil, err
	}

	limitClause, err := generateLimitClause(queryParamaters)
	if err != nil {
		return nil, err
	}

	offsetClause, err := generateOffsetClause(queryParamaters)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
	SELECT * 
	FROM reviews 
	%s %s %s %s`, whereClause, orderByClause, limitClause, offsetClause)

	var reviews []models.Review
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	// https://pkg.go.dev/github.com/jackc/pgx#hdr-Query_Interface
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var review models.Review
		err = rows.Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt, &review.Name, &review.Description, &review.ReviewURL, &review.ReviewScore, &review.MediaType, &review.GenreList, &review.CreatorList)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return nil, err
	}

	return reviews, nil
}

func (r Review) UpdateByID(pool *pgxpool.Pool, id int, m models.Review) error {
	// updated_at here is used as a simple locking mechanism to prevent data race
	query := `
	UPDATE reviews 
	SET name=$1, description=$2, review_url=$3, review_score=$4, media_type=$5, genre_list=$6, creator_list=$7, updated_at=now() 
	WHERE id = $8 and updated_at = $9`

	args := []interface{}{m.Name, m.Description, m.ReviewURL, m.ReviewScore, m.MediaType, m.GenreList, m.CreatorList, id, m.UpdatedAt}
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

func (r Review) DeleteByID(pool *pgxpool.Pool, id int) error {
	query := `
	DELETE FROM reviews
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	commandTag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("record not found")
	}

	return nil
}

// generateWhereClause generate a where clause based on the query parameters
func generateWhereClause(queryParamaters map[string][]string) (string, error) {
	var whereClause = make([]string, 0, 3)

	if name, found := queryParamaters["name"]; found {
		// using PostgreSQL build-in full text search
		whereClause = append(whereClause, fmt.Sprintf("(to_tsvector('simple', name) @@ plainto_tsquery('simple', '%s'))", name[0]))
	}

	if score, found := queryParamaters["scoreMin"]; found {
		val, err := strconv.ParseFloat(score[0], 32)
		if err != nil {
			return "", errors.New("query parameter 'scoreMin' must be float type and 0 <= scoreMin <= 10")
		} else if val < 0.0 || val > 10.0 {
			return "", errors.New("query parameter 'scoreMin' must be 0 <= scoreMin <= 10")
		}
		whereClause = append(whereClause, fmt.Sprintf("review_score >= %s", score[0]))
	}

	if genres, found := queryParamaters["genres"]; found {
		whereClause = append(whereClause, fmt.Sprintf("genre_list @> '{%s}'", strings.Join(genres, ",")))
	}

	if len(whereClause) != 0 {
		return "WHERE " + strings.Join(whereClause, " AND "), nil
	}

	return "", nil
}

var orderOptions = [...]string{"id", "-id", "name", "-name", "review_score", "-review_score"}

// generateOrderByClause generate a order by clause based on the query parameters
func generateOrderByClause(queryParamaters map[string][]string) (string, error) {
	order, found := queryParamaters["order"]
	if !found {
		return "ORDER BY id", nil
	}

	valid := false
	for _, o := range orderOptions {
		if order[0] == o {
			valid = true
			break
		}
	}

	if !valid {
		return "", errors.New(fmt.Sprintf("query parameter 'order' must be one of followings: %v", orderOptions))
	}

	if order[0][0] == '-' {
		return "ORDER BY " + order[0][1:] + " DESC", nil
	}
	return "ORDER BY " + order[0], nil
}

// generateLimitClause generate a limit clause based on the query parameters
func generateLimitClause(queryParamaters map[string][]string) (string, error) {
	if limit, found := queryParamaters["page_size"]; found {
		_, err := strconv.Atoi(limit[0])
		if err != nil {
			return "", errors.New("query parameter 'page_size' must be positive integer type")
		}
		return "LIMIT " + limit[0], nil
	}

	return "LIMIT 10", nil
}

// generateOffsetClause generate a limit clause based on the query parameters
func generateOffsetClause(queryParamaters map[string][]string) (string, error) {
	page_size := 10
	if limit, found := queryParamaters["page_size"]; found {
		p, err := strconv.Atoi(limit[0])
		if err != nil {
			return "", errors.New("query parameter 'page_size' must be positive integer type")
		}
		page_size = p
	}

	if offset, found := queryParamaters["page"]; found {
		page, err := strconv.Atoi(offset[0])
		if err != nil {
			return "", errors.New("query parameter 'page' must be positive integer type")
		}
		return fmt.Sprintf("OFFSET %d", (page-1)*page_size), nil
	}

	return "", nil
}
