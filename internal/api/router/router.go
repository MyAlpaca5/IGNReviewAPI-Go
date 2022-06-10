package router

import (
	"reflect"
	"strings"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/controllers"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewRouter() *gin.Engine {
	// Default With the Logger and Recovery middleware already attached
	router := gin.Default()

	// --- Bind Middlewares ---
	router.Use(middlewares.BodySizeLimiter(65_536))

	// --- Define Custom Validation Mehod ---
	// https://pkg.go.dev/github.com/go-playground/validator/v10#Validate.RegisterTagNameFunc
	// use the names which have been specified for JSON representations of structs, rather than normal Go field names
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	// --- Define Error Handlers ---
	// router.NoRoute(notFoundErrorHandler)
	router.HandleMethodNotAllowed = true
	// router.NoMethod(replyWithUnsupportedHTTPMethodError)

	// --- Bind Routes and Handlers ---
	router.GET("/api/healthcheck", controllers.HealthcheckHandler)
	router.GET("/api/reviews/:id", controllers.ShowReviewHandler)
	router.POST("/api/reviews", controllers.CreateReviewHandler)

	return router
}
