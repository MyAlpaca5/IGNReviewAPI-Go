package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type contextKey int

const tokenKey contextKey = 1

func contextWithToken(c *gin.Context, token *models.Token) *http.Request {
	ctx := context.WithValue(c.Request.Context(), tokenKey, token)
	return c.Request.WithContext(ctx)
}

func TokenFromContext(c *gin.Context) *models.Token {
	token, ok := c.Request.Context().Value(tokenKey).(*models.Token)
	if !ok {
		return nil
	}
	return token
}

func Authenticate(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Authorization")

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response := r_errors.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Error - Authorization header is not provided",
			}
			c.AbortWithStatusJSON(response.StatusCode, response)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 && headerParts[0] != "Bearer" {
			response := r_errors.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Error - Authorization header contains malformed data",
			}
			c.Header("WWW-Authenticate", "Bearer")
			c.AbortWithStatusJSON(response.StatusCode, response)
			return
		}

		tokenStr := headerParts[1]
		token, err := repositories.Token{}.ReadByToken(pool, models.TokenStrToHash(tokenStr))
		if err != nil || token.Expiry.Before(time.Now().UTC()) {
			response := r_errors.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Error - Bearer token is invalid or expired",
			}
			c.Header("WWW-Authenticate", "Bearer")
			c.AbortWithStatusJSON(response.StatusCode, response)
			return
		}

		c.Request = contextWithToken(c, &token)
		c.Next()
	}
}
