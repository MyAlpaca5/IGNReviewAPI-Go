package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/repositories"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReviewController struct {
	Repo repositories.Review
	Pool *pgxpool.Pool
}

// ReviewsGETHandler handles "POST /api/reviews" endpoint. It will insert a new review entry to the database.
func (ctrl ReviewController) CreateReviewHandler(c *gin.Context) {
	var reviewIn models.ReviewIn
	if err := c.ShouldBindJSON(&reviewIn); err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    RequestErr(err),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// check if not null fields are given by user, if not, return error
	if reviewIn.Name == nil || reviewIn.ReviewURL == nil || reviewIn.ReviewScore == nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Input Error - the following fields must be given: name, review_url, review_score",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// validate a few fields
	if err := reviewIn.Validate(); err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Input Error - request data does not satisfy requirements: %s", err.Error()),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	id, err := ctrl.Repo.Create(ctrl.Pool, reviewIn.ToReview())
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "DB Error - cannot create new review",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("review created successfully, assigned id is %d", id)})
}

// ShowReviewHandler handles "GET /api/reviews/:id" endpoint. It will fetch a review entry from database based on id.
func (ctrl ReviewController) ShowReviewHandler(c *gin.Context) {
	// validate path parameter 'id'
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Path Error - invalid 'id' parameter: %s", c.Param("id")),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	review, err := ctrl.Repo.ReadByID(ctrl.Pool, int(id))
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - cannot get review with id(%d), %s", id, err.Error()),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, review.ToReviewOut())
}

// UpdateReviewHandler handles "PATCH /api/reviews/:id" endpoint. It will update the review entry in the database based on the user input.
func (ctrl ReviewController) UpdateReviewHandler(c *gin.Context) {
	var reviewIn models.ReviewIn
	if err := c.ShouldBindJSON(&reviewIn); err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    RequestErr(err),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// validate a few fields
	if err := reviewIn.Validate(); err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Input Error - request data does not satisfy requirements: %s", err.Error()),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// validate path parameter 'id'
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Path Error - invalid 'id' parameter: %s", c.Param("id")),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// read original values
	original, err := ctrl.Repo.ReadByID(ctrl.Pool, int(id))
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - cannot found the review entry with id(%d)", id),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// create updated record
	var updated = original
	if reviewIn.Name != nil {
		updated.Name = *reviewIn.Name
	}
	if reviewIn.Description != nil {
		updated.Description = *reviewIn.Description
	}
	if reviewIn.ReviewScore != nil {
		updated.ReviewScore = *reviewIn.ReviewScore
	}
	if reviewIn.ReviewURL != nil {
		updated.ReviewURL = *reviewIn.ReviewURL
	}
	if reviewIn.MediaType != nil {
		updated.MediaType = *reviewIn.MediaType
	}
	if reviewIn.GenreList != nil {
		updated.GenreList = reviewIn.GenreList
	}
	if reviewIn.CreatorList != nil {
		updated.CreatorList = reviewIn.CreatorList
	}

	err = ctrl.Repo.UpdateByID(ctrl.Pool, int(id), updated)
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusConflict,
			Message:    fmt.Sprintf("DB Error - cannot update review with id = %d, potential data race happened, please try again.", id),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("review with id(%d) updated successfully", id)})
}

// DeleteReviewHandler handles "DELETE /api/reviews/:id" endpoint. It will delete a review entry from database based on id.
func (ctrl ReviewController) DeleteReviewHandler(c *gin.Context) {
	// validate path parameter 'id'
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Path Error - invalid 'id' parameter: %s", c.Param("id")),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	err = ctrl.Repo.DeleteByID(ctrl.Pool, int(id))
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - cannot delete review with id(%d), %s", id, err.Error()),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("review with id(%d) deleted successfully", id)})
}

// ListReviewsHandler handles "GET /api/reviews" endpoint.
func (ctrl ReviewController) ListReviewsHandler(c *gin.Context) {
	reviews, err := ctrl.Repo.ReadAll(ctrl.Pool, c.Request.URL.Query())
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - %s", err.Error()),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	reviewOuts := make([]models.ReviewOut, 0, len(reviews))
	for _, v := range reviews {
		reviewOuts = append(reviewOuts, v.ToReviewOut())
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviewOuts})
}

func RequestErr(err error) string {
	var ginErr gin.Error
	var validationErr validator.ValidationErrors
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return r_errors.JSONSyntaxError(syntaxError)
	case errors.As(err, &unmarshalTypeError):
		return r_errors.JSONUnmarshalTypeError(unmarshalTypeError)
	case errors.As(err, &ginErr):
		return r_errors.GINError(ginErr)
	case errors.As(err, &validationErr):
		return r_errors.ValidationError(validationErr)
	default:
		return "Unknown Error - " + err.Error()
	}
}
