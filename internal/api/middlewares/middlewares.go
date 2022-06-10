package middlewares

import (
	"fmt"
	"net/http"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/gin-gonic/gin"
)

// BodySizeLimiter will limit how large the request body can be in bytes
func BodySizeLimiter(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, r_errors.ResponseError{StatusCode: http.StatusInternalServerError, Message: fmt.Sprintf("Error - request body exceed %d bytes", maxSize)})
			return
		}

		c.Next()
	}
}
