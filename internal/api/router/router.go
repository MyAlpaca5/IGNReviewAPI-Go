package router

import (
	"reflect"
	"strings"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/controllers"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/middlewares"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(pool *pgxpool.Pool) *gin.Engine {
	// Default With the Logger and Recovery middleware already attached
	router := gin.Default()

	// --- Set Custom Validator Methods for JSON data ---
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
	}

	// --- Set Global Middlewares ---
	router.Use(middlewares.RequestRateLimiter(2, 4))
	router.Use(middlewares.BodySizeLimiter(65_536))

	// --- Set Error Handlers ---
	// router.NoRoute(notFoundErrorHandler)
	router.HandleMethodNotAllowed = true
	// router.NoMethod(replyWithUnsupportedHTTPMethodError)

	// --- Create Controllers ---
	var healthcheckController = controllers.HealthcheckController{}
	var reviewController = controllers.ReviewController{Repo: repositories.Review{}, Pool: pool}
	var userController = controllers.UserController{Repo: repositories.User{}, Pool: pool}

	// --- Set Routes and Handlers ---
	router.GET("/api/healthcheck", healthcheckController.HealthcheckHandler)
	router.GET("/api/reviews/:id", reviewController.ShowReviewHandler)
	router.DELETE("/api/reviews/:id", reviewController.DeleteReviewHandler)
	router.PATCH("/api/reviews/:id", reviewController.UpdateReviewHandler)
	router.GET("/api/reviews", reviewController.ListReviewsHandler)
	router.POST("/api/reviews", reviewController.CreateReviewHandler)
	router.POST("/api/users", userController.CreateUserHandler)

	return router
}
