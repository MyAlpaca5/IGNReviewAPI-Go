package middlewares

import (
	"net/http"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/gin-gonic/gin"
)

func Authorize(minRole int) gin.HandlerFunc {
	// assume the user has already been authenticated
	return func(c *gin.Context) {
		var token = TokenFromContext(c)
		if minRole < token.Role {
			response := r_errors.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Error - no permission to access this page",
			}
			c.AbortWithStatusJSON(response.StatusCode, response)
			return
		}

		c.Next()
	}
}
