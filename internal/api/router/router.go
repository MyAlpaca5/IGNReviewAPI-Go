package router

import (
	"reflect"
	"regexp"
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

	// --- Set Custom Validator Methods ---
	// https://pkg.go.dev/github.com/go-playground/validator/v10#Validate
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// use the names which have been specified for JSON representations of structs, rather than normal Go field names
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// check the review_url is a valid url to a IGN review article
		v.RegisterValidation("ignurl", func(fl validator.FieldLevel) bool {
			url := fl.Field().String()
			matched, err := regexp.MatchString("^https://www.ign.com/articles/[12][0-9]{3}/(0[1-9]|1[0-2])/(0[1-9]|[1-2][0-9]|3[01])/.+$", url)
			if !matched || err != nil {
				return false
			}
			return true
		})
	}

	// --- Set Middlewares ---
	router.Use(middlewares.BodySizeLimiter(65_536))

	// --- Set Error Handlers ---
	// router.NoRoute(notFoundErrorHandler)
	router.HandleMethodNotAllowed = true
	// router.NoMethod(replyWithUnsupportedHTTPMethodError)

	// --- Set Routes and Handlers ---
	router.GET("/api/healthcheck", controllers.HealthcheckHandler)
	router.GET("/api/reviews/:id", controllers.ShowReviewHandler)
	router.POST("/api/reviews", controllers.CreateReviewHandler)

	return router
}
