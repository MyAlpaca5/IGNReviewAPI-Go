package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ReviewsGETHandler handles "POST /api/reviews" endpoint.
func CreateReviewHandler(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
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

	c.JSON(http.StatusAccepted, review)
}

// ReviewsGETHandler handles "GET /api/reviews/:id" endpoint. TODO: for now, just return plain text.
func ShowReviewHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(id)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, r_errors.ResponseError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Path Error - invalid id parameter: %s", c.Param("id"))})
		return
	}

	// TEST
	r := models.Review{
		Name:        strconv.Itoa(int(id)),
		Description: "short test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now().Add(time.Hour),
		ReviewURL:   "url to here",
		ReviewScore: 4.5,
		MediaType:   "abc",
		GenreList:   nil,
		CreatorList: []string{},
	}

	c.JSON(http.StatusOK, r)
	// TEST END
}
