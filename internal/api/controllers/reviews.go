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
	Repo repositories.Repo[models.Review]
	Pool *pgxpool.Pool
}

// ReviewsGETHandler handles "POST /api/reviews" endpoint. It will insert a new review entry to the database.
func (ctrl ReviewController) CreateReviewHandler(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		// check specific error for better response message
		var ginErr gin.Error
		var validationErr validator.ValidationErrors
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
		}

		switch {
		case errors.As(err, &syntaxError):
			response.Message = r_errors.JSONSyntaxError(syntaxError)
		case errors.As(err, &unmarshalTypeError):
			response.Message = r_errors.JSONUnmarshalTypeError(unmarshalTypeError)
		case errors.As(err, &ginErr):
			response.Message = r_errors.GINError(ginErr)
		case errors.As(err, &validationErr):
			response.Message = r_errors.ValidationError(validationErr)
		default:
			response.Message = "Unknown Error - " + err.Error()
		}

		c.JSON(response.StatusCode, response)
		return
	}

	id, err := ctrl.Repo.Create(ctrl.Pool, review)
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "DB Error - cannot create new review",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("review created successfully, assigned id is %d", id)})
}

// ShowReviewHandler handles "GET /api/reviews/:id" endpoint. It will fetch a review entry from database based on id.
func (ctrl ReviewController) ShowReviewHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Path Error - invalid 'id' parameter: %s", c.Param("id"))})
		return
	}

	review, err := ctrl.Repo.Read(ctrl.Pool, int(id))
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - cannot get review with id = %d, err = %s", id, err.Error()),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, review)
}

// UpdateReviewHandler handles "PATCH /api/reviews/:id" endpoint. It will update the review entry in the database based on the user input.
// Note: created_at and review_url fields are deemed as non-mutable, therefore, even user pass new data for those two field, they will be ignored.
func (ctrl ReviewController) UpdateReviewHandler(c *gin.Context) {
	var review models.ReviewForUpdate

	if err := c.ShouldBindJSON(&review); err != nil {
		// check specific error for better response message
		var ginErr gin.Error
		var validationErr validator.ValidationErrors
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
		}

		switch {
		case errors.As(err, &syntaxError):
			response.Message = r_errors.JSONSyntaxError(syntaxError)
		case errors.As(err, &unmarshalTypeError):
			response.Message = r_errors.JSONUnmarshalTypeError(unmarshalTypeError)
		case errors.As(err, &ginErr):
			response.Message = r_errors.GINError(ginErr)
		case errors.As(err, &validationErr):
			response.Message = r_errors.ValidationError(validationErr)
		default:
			response.Message = "Unknown Error - " + err.Error()
		}

		c.JSON(response.StatusCode, response)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Path Error - invalid 'id' parameter: %s", c.Param("id"))})
		return
	}

	// read original values
	original, err := ctrl.Repo.Read(ctrl.Pool, int(id))
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - cannot found the review entry with id = %d", id),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// update values
	var updated = original
	if review.Name != nil {
		if *review.Name == "" {
			c.JSON(http.StatusBadRequest, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: "Validation Error - 'name' cannot be empty string"})
			return
		}

		updated.Name = *review.Name
	}

	if review.Description != nil {
		updated.Description = *review.Description
	}

	if review.ReviewScore != nil {
		updated.ReviewScore = *review.ReviewScore
	}

	if review.MediaType != nil {
		updated.MediaType = *review.MediaType
	}

	if review.GenreList != nil {
		updated.GenreList = review.GenreList
	}

	if review.CreatorList != nil {
		updated.CreatorList = review.CreatorList
	}

	if review.UpdatedAt.Before(original.UpdatedAt) || review.UpdatedAt.Equal(original.UpdatedAt) {
		c.JSON(http.StatusBadRequest, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: "Validation Error - the new value of 'update_at' field is not newer than that of the original 'update_at' field."})
		return
	}
	updated.UpdatedAt = review.UpdatedAt

	// update database
	err = ctrl.Repo.Update(ctrl.Pool, int(id), updated)
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusConflict,
			Message:    fmt.Sprintf("DB Error - cannot update review with id = %d, potential data race happened, please try again.", id),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("review with id = %d updated successfully", id)})
}

// DeleteReviewHandler handles "DELETE /api/reviews/:id" endpoint. It will delete a review entry from database based on id.
func (ctrl ReviewController) DeleteReviewHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Path Error - invalid 'id' parameter: %s", c.Param("id"))})
		return
	}

	err = ctrl.Repo.Delete(ctrl.Pool, int(id))
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - cannot delete review with id = %d, err = %s", id, err.Error()),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("review with id = %d deleted successfully", id)})
}

// ListReviewsHandler handles "GET /api/reviews" endpoint.
func (ctrl ReviewController) ListReviewsHandler(c *gin.Context) {
	reviews, err := ctrl.Repo.ReadAll(ctrl.Pool, c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("DB Error - %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
