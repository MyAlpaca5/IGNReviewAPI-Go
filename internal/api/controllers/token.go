package controllers

import (
	"fmt"
	"net/http"
	"time"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TokenController struct {
	userRepo repositories.User
	Repo     repositories.Token
	Pool     *pgxpool.Pool
}

// CreateUserHandler handles "POST /api/tokens/authentication" endpoint. It will insert a new authentication token entry to the database.
func (ctrl TokenController) CreateAuthenticationTokenHandler(c *gin.Context) {
	var credential struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credential); err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    r_errors.GetBindingErrorStr(err),
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// check if not null fields are given by user, if not, return error
	if credential.Username == nil || credential.Password == nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Input Error - the following fields must be given: username, password",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// check if user exists
	user, err := ctrl.userRepo.ReadByUsername(*credential.Username)
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Request Error - no record for this user",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// check if password matches
	valid, err := user.ComparePassword(*credential.Password)
	if err != nil || !valid {
		response := r_errors.ResponseError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Request Error - authentication failed",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// create a new token for user
	token, err := models.NewToken(user.ID, time.Now().Add(24*time.Hour).UTC(), models.RoleSimple)
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Request Error - cannot create authentication token, please try later",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	// insert new record into database
	err = ctrl.Repo.Create(ctrl.Pool, token)
	if err != nil {
		response := r_errors.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Request Error - cannot create authentication token, please try later",
		}
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("authentication token: %s, expiry: %d-%02d-%02dT%02d:%02d:%02d", token.TokenString, token.Expiry.Year(), token.Expiry.Month(), token.Expiry.Day(), token.Expiry.Hour(), token.Expiry.Minute(), token.Expiry.Second())})
}
