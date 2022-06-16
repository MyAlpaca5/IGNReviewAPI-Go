package controllers

import (
	"fmt"
	"net/http"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserController struct {
	Repo repositories.User
	Pool *pgxpool.Pool
}

// CreateUserHandler handles "POST /api/users" endpoint. It will insert a new user entry to the database.
func (ctrl UserController) CreateUserHandler(c *gin.Context) {
	var userIn models.UserIn
	if err := c.ShouldBindJSON(&userIn); err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    RequestErr(err),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// check if not null fields are given by user, if not, return error
	if userIn.Username == nil || userIn.Password == nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Input Error - the following fields must be given: username, password",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// validate a few fields
	if err := userIn.Validate(); err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Input Error - request data does not satisfy requirements: %s", err.Error()),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	user, err := userIn.ToUser()
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Parse Error - %s", err.Error()),
		}
		c.JSON(response.StatusCode, response)
	}

	id, err := ctrl.Repo.Create(ctrl.Pool, user)
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("DB Error - %s", err.Error()),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("user created successfully, assigned id is %d", id)})
}
