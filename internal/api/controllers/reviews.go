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
	repo repositories.ReviewRepo
	Pool *pgxpool.Pool
}

// ReviewsGETHandler handles "POST /api/reviews" endpoint.
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

	id, err := ctrl.repo.Create(ctrl.Pool, review)
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

// ReviewsGETHandler handles "GET /api/reviews/:id" endpoint. TODO: for now, just return plain text.
func (ctrl ReviewController) ShowReviewHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(id)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Path Error - invalid id parameter: %s", c.Param("id"))})
		return
	}

	review, err := ctrl.repo.Read(ctrl.Pool, int(id))
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
