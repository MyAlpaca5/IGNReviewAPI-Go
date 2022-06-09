package router

import (
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/controllers"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// Default With the Logger and Recovery middleware already attached
	router := gin.Default()

	// Bind Middlewares

	// Bind routes and handlers
	router.GET("/api/healthcheck", controllers.HealthcheckGETHandler)

	return router
}
