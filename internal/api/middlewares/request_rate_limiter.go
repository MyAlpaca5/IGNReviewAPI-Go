package middlewares

import (
	"net/http"
	"sync"

	r_errors "github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RequestRateLimiter will limit the number of requests a client can sent over a short period of time. The restriction is IP based.
// TODO: set up a eviction policy to delete old IP from the map when the time expired. Idea: LRU cache
func RequestRateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	var (
		ipMap = make(map[string]*rate.Limiter)
		mu    sync.Mutex
	)

	return func(c *gin.Context) {
		// extract client IP address
		ip := c.ClientIP()

		mu.Lock()
		// create a new entry in map, if this is first time this IP shows up
		if _, found := ipMap[ip]; !found {
			ipMap[ip] = rate.NewLimiter(r, b)
		}

		if !ipMap[ip].Allow() {
			mu.Unlock() // need to release the lock here, otherwise, no further request from this IP will be accepted
			response := r_errors.ResponseError{
				StatusCode: http.StatusTooManyRequests,
				Message:    "too many request, please try later",
			}
			c.AbortWithStatusJSON(response.StatusCode, response)
			return
		}
		mu.Unlock()

		c.Next()
	}
}
